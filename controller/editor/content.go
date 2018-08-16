package editor

import (
	"net/http"

	"github.com/acoshift/hime"

	"github.com/acoshift/acourse/context/appctx"
	"github.com/acoshift/acourse/controller/share"
	"github.com/acoshift/acourse/entity"
	"github.com/acoshift/acourse/repository"
	"github.com/acoshift/acourse/view"
)

func (c *ctrl) contentList(ctx *hime.Context) error {
	id := ctx.FormValue("id")

	course, err := repository.GetCourse(ctx, id)
	if err == entity.ErrNotFound {
		return share.NotFound(ctx)
	}
	if err != nil {
		return err
	}
	course.Contents, err = repository.GetCourseContents(ctx, id)
	if err != nil {
		return err
	}

	p := view.Page(ctx)
	p["Course"] = course
	return ctx.View("editor.content", p)
}

func (c *ctrl) postContentList(ctx *hime.Context) error {
	id := ctx.FormValue("id")

	if ctx.FormValue("action") == "delete" {
		contentID := ctx.FormValue("contentId")

		err := repository.DeleteCourseContent(ctx, id, contentID)
		if err != nil {
			return err
		}
	}
	return ctx.RedirectToGet()
}

func (c *ctrl) contentCreate(ctx *hime.Context) error {
	id := ctx.FormValue("id")

	course, err := repository.GetCourse(ctx, id)
	if err != nil {
		return err
	}

	p := view.Page(ctx)
	p["Course"] = course
	return ctx.View("editor.content-create", p)
}

func (c *ctrl) postContentCreate(ctx *hime.Context) error {
	id := ctx.FormValue("id")

	var (
		title   = ctx.FormValue("title")
		desc    = ctx.FormValue("desc")
		videoID = ctx.FormValue("videoId")
	)

	_, err := repository.RegisterCourseContent(ctx, &entity.RegisterCourseContent{
		CourseID:  id,
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

	content, err := repository.GetCourseContent(ctx, id)
	if err == entity.ErrNotFound {
		return share.NotFound(ctx)
	}
	if err != nil {
		return err
	}

	course, err := repository.GetCourse(ctx, content.CourseID)
	if err != nil {
		return err
	}

	user := appctx.GetUser(ctx)
	// user is not course owner
	if user.ID != course.UserID {
		return ctx.Status(http.StatusForbidden).StatusText()
	}

	p := view.Page(ctx)
	p["Course"] = course
	p["Content"] = content
	return ctx.View("editor.content-edit", p)
}

func (c *ctrl) postContentEdit(ctx *hime.Context) error {
	// course content id
	id := ctx.FormValue("id")

	content, err := repository.GetCourseContent(ctx, id)
	if err == entity.ErrNotFound {
		return share.NotFound(ctx)
	}
	if err != nil {
		return err
	}

	course, err := repository.GetCourse(ctx, content.CourseID)
	if err != nil {
		return err
	}

	user := appctx.GetUser(ctx)
	// user is not course owner
	if user.ID != course.UserID {
		return ctx.Status(http.StatusForbidden).StatusText()
	}

	var (
		title   = ctx.FormValue("title")
		desc    = ctx.FormValue("desc")
		videoID = ctx.FormValue("videoId")
	)

	err = repository.UpdateCourseContent(ctx, course.ID, id, title, desc, videoID)
	if err != nil {
		return err
	}

	return ctx.RedirectTo("editor.content", ctx.Param("id", course.ID))
}
