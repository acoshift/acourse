package app

import (
	"github.com/moonrhythm/hime"

	"github.com/acoshift/acourse/internal/app/view"
)

func index(ctx *hime.Context) error {
	if ctx.URL.Path != "/" {
		return view.NotFound(ctx)
	}

	courses, err := listPublicCourses(ctx)
	if err != nil {
		return err
	}

	p := view.Page(ctx)
	p.Data["Courses"] = courses
	return ctx.View("app.index", p)
}
