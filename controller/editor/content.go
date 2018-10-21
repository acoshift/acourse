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

	x, err := c.Repository.GetCourse(ctx, id)
	if err == entity.ErrNotFound {
		return share.NotFound(ctx)
	}
	if err != nil {
		return err
	}

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
