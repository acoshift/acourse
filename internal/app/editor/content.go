package editor

import (
	"net/http"

	"github.com/moonrhythm/hime"

	"github.com/acoshift/acourse/internal/app/view"
	"github.com/acoshift/acourse/internal/pkg/app"
	"github.com/acoshift/acourse/internal/pkg/context/appctx"
	"github.com/acoshift/acourse/internal/pkg/course"
)

func getContentList(ctx *hime.Context) error {
	id := ctx.FormValue("id")

	c, err := course.Get(ctx, id)
	if err == course.ErrNotFound {
		return view.NotFound(ctx)
	}
	if err != nil {
		return err
	}

	contents, err := course.GetContents(ctx, id)
	if err != nil {
		return err
	}
	c.Contents = contents

	p := view.Page(ctx)
	p.Data["Course"] = c
	return ctx.View("editor.content", p)
}

func postContentList(ctx *hime.Context) error {
	if ctx.FormValue("action") == "delete" {
		contentID := ctx.FormValue("contentId")

		err := course.DeleteContent(ctx, contentID)
		if app.IsUIError(err) {
			// TODO: use flash
			return ctx.Status(http.StatusBadRequest).Error(err.Error())
		}
		if err != nil {
			return err
		}
	}
	return ctx.RedirectToGet()
}

func getContentCreate(ctx *hime.Context) error {
	id := ctx.FormValue("id")

	c, err := course.Get(ctx, id)
	if err == course.ErrNotFound {
		return view.NotFound(ctx)
	}
	if err != nil {
		return err
	}

	p := view.Page(ctx)
	p.Data["Course"] = c
	return ctx.View("editor.content-create", p)
}

func postContentCreate(ctx *hime.Context) error {
	id := ctx.FormValue("id")

	var (
		title   = ctx.FormValue("title")
		desc    = ctx.FormValue("desc")
		videoID = ctx.FormValue("videoId")
	)

	_, err := course.CreateContent(ctx, &course.CreateContentArgs{
		ID:        id,
		Title:     title,
		LongDesc:  desc,
		VideoID:   videoID,
		VideoType: course.Youtube,
	})
	if err != nil {
		return err
	}

	return ctx.RedirectTo("editor.content", ctx.Param("id", ctx.FormValue("id")))
}

func getContentEdit(ctx *hime.Context) error {
	// course content id
	id := ctx.FormValue("id")

	content, err := course.GetContent(ctx, id)
	if err == course.ErrNotFound {
		return view.NotFound(ctx)
	}
	if err != nil {
		return err
	}

	c, err := course.Get(ctx, content.CourseID)
	if err != nil {
		return err
	}

	user := appctx.GetUser(ctx)
	// user is not course owner
	if user.ID != c.Owner.ID {
		return ctx.Status(http.StatusForbidden).StatusText()
	}

	p := view.Page(ctx)
	p.Data["Course"] = c
	p.Data["Content"] = content
	return ctx.View("editor.content-edit", p)
}

func postContentEdit(ctx *hime.Context) error {
	// course content id
	id := ctx.FormValue("id")

	content, err := course.GetContent(ctx, id)
	if err == course.ErrNotFound {
		return view.NotFound(ctx)
	}
	if err != nil {
		return err
	}

	c, err := course.Get(ctx, content.CourseID)
	if err != nil {
		return err
	}

	user := appctx.GetUser(ctx)
	// user is not course owner
	// TODO: move to service
	if user.ID != c.Owner.ID {
		return ctx.Status(http.StatusForbidden).StatusText()
	}

	var (
		title   = ctx.FormValue("title")
		desc    = ctx.FormValue("desc")
		videoID = ctx.FormValue("videoId")
	)

	err = course.UpdateContent(ctx, &course.UpdateContentArgs{
		ContentID: id,
		Title:     title,
		Desc:      desc,
		VideoID:   videoID,
	})
	if err != nil {
		return err
	}

	return ctx.RedirectTo("editor.content", ctx.Param("id", content.CourseID))
}
