package editor

import (
	"time"

	"github.com/moonrhythm/dispatcher"
	"github.com/moonrhythm/hime"

	"github.com/acoshift/acourse/internal/context/appctx"
	"github.com/acoshift/acourse/internal/model/app"
	"github.com/acoshift/acourse/internal/model/course"
	"github.com/acoshift/acourse/internal/view"
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

	q := course.Create{
		UserID:    appctx.GetUserID(ctx),
		Title:     title,
		ShortDesc: shortDesc,
		LongDesc:  desc,
		Image:     image,
		Start:     start,
	}
	err := dispatcher.Dispatch(ctx, &q)
	if app.IsUIError(err) {
		f.Add("Errors", err.Error())
		return ctx.RedirectToGet()
	}
	if err != nil {
		return err
	}

	link := course.GetURL{ID: q.Result}
	dispatcher.Dispatch(ctx, &link)
	if link.Result == "" {
		return ctx.RedirectTo("app.course", q.Result)
	}
	return ctx.RedirectTo("app.course", link.Result)
}

func (c *ctrl) courseEdit(ctx *hime.Context) error {
	id := ctx.FormValue("id")
	getCourse := course.Get{ID: id}
	err := dispatcher.Dispatch(ctx, &getCourse)
	if err != nil {
		return err
	}
	x := getCourse.Result

	p := view.Page(ctx)
	p.Data["Course"] = x
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

	q := course.Update{
		ID:        id,
		Title:     title,
		ShortDesc: shortDesc,
		LongDesc:  desc,
		Image:     image,
		Start:     start,
	}
	err := dispatcher.Dispatch(ctx, &q)
	if app.IsUIError(err) {
		f.Add("Errors", err.Error())
		return ctx.RedirectToGet()
	}
	if err != nil {
		return err
	}

	link := course.GetURL{ID: id}
	dispatcher.Dispatch(ctx, &link)
	if link.Result == "" {
		return ctx.RedirectTo("app.course", id)
	}
	return ctx.RedirectTo("app.course", link.Result)
}
