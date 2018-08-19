package editor

import (
	"time"

	"github.com/acoshift/hime"

	"github.com/acoshift/acourse/context/appctx"
	"github.com/acoshift/acourse/service"
	"github.com/acoshift/acourse/view"
)

func (c *ctrl) courseCreate(ctx *hime.Context) error {
	return ctx.View("editor.course-create", view.Page(ctx))
}

func (c *ctrl) postCourseCreate(ctx *hime.Context) error {
	f := appctx.GetFlash(ctx)

	var (
		title     = ctx.PostFormValueTrimSpace("title")
		shortDesc = ctx.PostFormValueTrimSpace("shortDesc")
		desc      = ctx.PostFormValue("desc")
		start     time.Time
		// assignment, _ = strconv.ParseBool(ctx.FormValue("assignment"))
	)
	if len(title) == 0 {
		f.Add("Errors", "title required")
		return ctx.RedirectToGet()
	}

	if v := ctx.FormValue("start"); v != "" {
		start, _ = time.Parse("2006-01-02", v)
	}

	image, _ := ctx.FormFileHeaderNotEmpty("image")

	courseID, err := c.Service.CreateCourse(ctx, &service.CreateCourse{
		Title:     title,
		ShortDesc: shortDesc,
		LongDesc:  desc,
		Image:     image,
		Start:     start,
	})
	if service.IsUIError(err) {
		f.Add("Errors", err.Error())
		return ctx.RedirectToGet()
	}
	if err != nil {
		return err
	}

	link, _ := c.Repository.GetCourseURL(ctx, courseID)
	if link == "" {
		return ctx.RedirectTo("app.course", courseID)
	}
	return ctx.RedirectTo("app.course", link)
}

func (c *ctrl) courseEdit(ctx *hime.Context) error {
	id := ctx.FormValue("id")
	course, err := c.Repository.GetCourse(ctx, id)
	if err != nil {
		return err
	}

	p := view.Page(ctx)
	p["Course"] = course
	return ctx.View("editor.course-edit", p)
}

func (c *ctrl) postCourseEdit(ctx *hime.Context) error {
	id := ctx.FormValue("id")

	f := appctx.GetFlash(ctx)

	var (
		title     = ctx.FormValue("title")
		shortDesc = ctx.FormValue("shortDesc")
		desc      = ctx.FormValue("desc")
		start     time.Time
		// assignment, _ = strconv.ParseBool(ctx.FormValue("assignment"))
	)
	if len(title) == 0 {
		f.Add("Errors", "title required")
		return ctx.RedirectToGet()
	}

	if v := ctx.FormValue("start"); len(v) > 0 {
		start, _ = time.Parse("2006-01-02", v)
	}

	image, _ := ctx.FormFileHeaderNotEmpty("image")

	err := c.Service.UpdateCourse(ctx, &service.UpdateCourse{
		ID:        id,
		Title:     title,
		ShortDesc: shortDesc,
		LongDesc:  desc,
		Image:     image,
		Start:     start,
	})
	if service.IsUIError(err) {
		f.Add("Errors", err.Error())
		return ctx.RedirectToGet()
	}
	if err != nil {
		return err
	}

	link, _ := c.Repository.GetCourseURL(ctx, id)
	if link == "" {
		return ctx.RedirectTo("app.course", id)
	}
	return ctx.RedirectTo("app.course", link)
}
