package editor

import (
	"time"

	"github.com/moonrhythm/hime"

	"github.com/acoshift/acourse/internal/app/view"
	"github.com/acoshift/acourse/internal/pkg/app"
	"github.com/acoshift/acourse/internal/pkg/context/appctx"
	"github.com/acoshift/acourse/internal/pkg/course"
)

func getCourseCreate(ctx *hime.Context) error {
	return ctx.View("editor.course-create", view.Page(ctx))
}

func postCourseCreate(ctx *hime.Context) error {
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

	courseID, err := course.Create(ctx, &course.CreateArgs{
		UserID:    appctx.GetUserID(ctx),
		Title:     title,
		ShortDesc: shortDesc,
		LongDesc:  desc,
		Image:     image,
		Start:     start,
	})
	if app.IsUIError(err) {
		f.Add("Errors", err.Error())
		return ctx.RedirectToGet()
	}
	if err != nil {
		return err
	}

	link, _ := course.GetURL(ctx, courseID)
	if link == "" {
		return ctx.RedirectTo("app.course", courseID)
	}
	return ctx.RedirectTo("app.course", link)
}

func getCourseEdit(ctx *hime.Context) error {
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
	return ctx.View("editor.course-edit", p)
}

func postCourseEdit(ctx *hime.Context) error {
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

	err := course.Update(ctx, &course.UpdateArgs{
		ID:        id,
		Title:     title,
		ShortDesc: shortDesc,
		LongDesc:  desc,
		Image:     image,
		Start:     start,
	})
	if app.IsUIError(err) {
		f.Add("Errors", err.Error())
		return ctx.RedirectToGet()
	}
	if err != nil {
		return err
	}

	link, _ := course.GetURL(ctx, id)
	if link == "" {
		return ctx.RedirectTo("app.course", id)
	}
	return ctx.RedirectTo("app.course", link)
}
