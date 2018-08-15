package app

import (
	"github.com/acoshift/hime"

	"github.com/acoshift/acourse/repository"
	"github.com/acoshift/acourse/view"
)

func index(ctx *hime.Context) error {
	if ctx.Request().URL.Path != "/" {
		return notFound(ctx)
	}

	courses, err := repository.ListPublicCourses(ctx)
	if err != nil {
		return err
	}

	p := view.Page(ctx)
	p["Courses"] = courses
	return ctx.View("index", p)
}
