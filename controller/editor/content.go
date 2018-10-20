package editor

import (
	"net/http"

	"github.com/moonrhythm/hime"

	"github.com/acoshift/acourse/context/appctx"
	"github.com/acoshift/acourse/controller/share"
	"github.com/acoshift/acourse/entity"
	"github.com/acoshift/acourse/service"
	"github.com/acoshift/acourse/view"
)

func (c *ctrl) contentList(ctx *hime.Context) error {
	id := ctx.FormValue("id")

	course, err := c.Repository.GetCourse(ctx, id)
	if err == entity.ErrNotFound {
		return share.NotFound(ctx)
	}
	if err != nil {
		return err
	}
	course.Contents, err = c.Service.ListCourseContents(ctx, id)
	if err != nil {
		return err
	}

	p := view.Page(ctx)
	p.Data["Course"] = course
	return ctx.View("editor.content", p)
}

func (c *ctrl) postContentList(ctx *hime.Context) error {
	if ctx.FormValue("action") == "delete" {
		contentID := ctx.FormValue("contentId")

		err := c.Service.DeleteCourseContent(ctx, contentID)
		if service.IsUIError(err) {
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

	course, err := c.Repository.GetCourse(ctx, id)
	if err != nil {
		return err
	}

	p := view.Page(ctx)
	p.Data["Course"] = course
	return ctx.View("editor.content-create", p)
}

func (c *ctrl) postContentCreate(ctx *hime.Context) error {
	id := ctx.FormValue("id")

	var (
		title   = ctx.FormValue("title")
		desc    = ctx.FormValue("desc")
		videoID = ctx.FormValue("videoId")
	)

	_, err := c.Service.CreateCourseContent(ctx, &entity.RegisterCourseContent{
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

	content, err := c.Service.GetCourseContent(ctx, id)
	if err == entity.ErrNotFound {
		return share.NotFound(ctx)
	}
	if err != nil {
		return err
	}

	course, err := c.Repository.GetCourse(ctx, content.CourseID)
	if err != nil {
		return err
	}

	user := appctx.GetUser(ctx)
	// user is not course owner
	if user.ID != course.UserID {
		return ctx.Status(http.StatusForbidden).StatusText()
	}

	p := view.Page(ctx)
	p.Data["Course"] = course
	p.Data["Content"] = content
	return ctx.View("editor.content-edit", p)
}

func (c *ctrl) postContentEdit(ctx *hime.Context) error {
	// course content id
	id := ctx.FormValue("id")

	content, err := c.Service.GetCourseContent(ctx, id)
	if err == entity.ErrNotFound {
		return share.NotFound(ctx)
	}
	if err != nil {
		return err
	}

	course, err := c.Repository.GetCourse(ctx, content.CourseID)
	if err != nil {
		return err
	}

	user := appctx.GetUser(ctx)
	// user is not course owner
	// TODO: move to service
	if user.ID != course.UserID {
		return ctx.Status(http.StatusForbidden).StatusText()
	}

	var (
		title   = ctx.FormValue("title")
		desc    = ctx.FormValue("desc")
		videoID = ctx.FormValue("videoId")
	)

	err = c.Service.UpdateCourseContent(ctx, id, title, desc, videoID)
	if err != nil {
		return err
	}

	return ctx.RedirectTo("editor.content", ctx.Param("id", course.ID))
}
