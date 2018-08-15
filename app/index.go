package app

import (
	"github.com/acoshift/hime"

	"github.com/acoshift/acourse/repository"
)

func index(ctx *hime.Context) error {
	if ctx.Request().URL.Path != "/" {
		return notFound(ctx)
	}

	courses, err := repository.ListPublicCourses(ctx)
	if err != nil {
		return err
	}

	p := page(ctx)
	p["Courses"] = courses
	return ctx.View("index", p)
}
