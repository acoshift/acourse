package app

import (
	"github.com/acoshift/hime"

	"github.com/acoshift/acourse/repository"
)

func index(ctx hime.Context) hime.Result {
	if ctx.Request().URL.Path != "/" {
		return notFound(ctx)
	}

	courses, err := repository.ListPublicCourses(db, cachePool, cachePrefix)
	must(err)

	page := newPage(ctx)
	page["Courses"] = courses
	return ctx.View("index", page)
}
