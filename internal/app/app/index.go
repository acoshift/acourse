package app

import (
	"github.com/moonrhythm/hime"

	"github.com/acoshift/acourse/internal/app/view"
	"github.com/acoshift/acourse/internal/pkg/course"
)

func index(ctx *hime.Context) error {
	if ctx.URL.Path != "/" {
		return view.NotFound(ctx)
	}

	courses, err := course.GetPublicCards(ctx)
	if err != nil {
		return err
	}

	p := view.Page(ctx)
	p.Data["Courses"] = courses
	return ctx.View("app.index", p)
}
