package editor

import (
	"net/http"

	"github.com/moonrhythm/dispatcher"
	"github.com/moonrhythm/hime"

	"github.com/acoshift/acourse/context/appctx"
	"github.com/acoshift/acourse/controller/share"
	"github.com/acoshift/acourse/entity"
	"github.com/acoshift/acourse/model/app"
	"github.com/acoshift/acourse/model/course"
	"github.com/acoshift/acourse/view"
)

func (c *ctrl) contentList(ctx *hime.Context) error {
	id := ctx.FormValue("id")

	getCourse := course.Get{ID: id}
	err := dispatcher.Dispatch(ctx, &getCourse)
	if err == entity.ErrNotFound {
		return share.NotFound(ctx)
	}
	if err != nil {
		return err
	}
	x := getCourse.Result

	{
		q := course.ListContents{ID: id}
		err := dispatcher.Dispatch(ctx, &q)
		if err != nil {
			return err
		}
		x.Contents = q.Result
	}

	p := view.Page(ctx)
	p.Data["Course"] = x
	return ctx.View("editor.content", p)
}

func (c *ctrl) postContentList(ctx *hime.Context) error {
	if ctx.FormValue("action") == "delete" {
		contentID := ctx.FormValue("contentId")

		err := dispatcher.Dispatch(ctx, &course.DeleteContent{ContentID: contentID})
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

func (c *ctrl) contentCreate(ctx *hime.Context) error {
	id := ctx.FormValue("id")

	getCourse := course.Get{ID: id}
	err := dispatcher.Dispatch(ctx, &getCourse)
	if err != nil {
		return err
	}
	x := getCourse.Result

	p := view.Page(ctx)
	p.Data["Course"] = x
	return ctx.View("editor.content-create", p)
}

func (c *ctrl) postContentCreate(ctx *hime.Context) error {
	id := ctx.FormValue("id")

	var (
		title   = ctx.FormValue("title")
		desc    = ctx.FormValue("desc")
		videoID = ctx.FormValue("videoId")
	)

	err := dispatcher.Dispatch(ctx, &course.CreateContent{
		ID:        id,
		Title:     title,
		LongDesc:  desc,
		VideoID:   videoID,
		VideoType: entity.Youtube,
	})
	if err != nil {
		return err
	}

	return ctx.RedirectTo("editor.content", ctx.Param("id", ctx.FormValue("id")))
}

func (c *ctrl) contentEdit(ctx *hime.Context) error {
	// course content id
	id := ctx.FormValue("id")

	getContent := course.GetContent{ContentID: id}
	err := dispatcher.Dispatch(ctx, &getContent)
	if err == entity.ErrNotFound {
		return share.NotFound(ctx)
	}
	if err != nil {
		return err
	}
	content := getContent.Result

	getCourse := course.Get{ID: content.CourseID}
	err = dispatcher.Dispatch(ctx, &getCourse)
	if err != nil {
		return err
	}
	x := getCourse.Result

	user := appctx.GetUser(ctx)
	// user is not course owner
	if user.ID != x.UserID {
		return ctx.Status(http.StatusForbidden).StatusText()
	}

	p := view.Page(ctx)
	p.Data["Course"] = x
	p.Data["Content"] = content
	return ctx.View("editor.content-edit", p)
}

func (c *ctrl) postContentEdit(ctx *hime.Context) error {
	// course content id
	id := ctx.FormValue("id")

	getContent := course.GetContent{ContentID: id}
	err := dispatcher.Dispatch(ctx, &getContent)
	if err == entity.ErrNotFound {
		return share.NotFound(ctx)
	}
	if err != nil {
		return err
	}
	content := getContent.Result
	if err == entity.ErrNotFound {
		return share.NotFound(ctx)
	}
	if err != nil {
		return err
	}

	{
		getCourse := course.Get{ID: content.CourseID}
		err := dispatcher.Dispatch(ctx, &getCourse)
		if err != nil {
			return err
		}
		x := getCourse.Result

		user := appctx.GetUser(ctx)
		// user is not course owner
		// TODO: move to service
		if user.ID != x.UserID {
			return ctx.Status(http.StatusForbidden).StatusText()
		}
	}

	var (
		title   = ctx.FormValue("title")
		desc    = ctx.FormValue("desc")
		videoID = ctx.FormValue("videoId")
	)

	err = dispatcher.Dispatch(ctx, &course.UpdateContent{
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
