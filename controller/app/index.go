package app

import (
	"github.com/moonrhythm/hime"

	"github.com/acoshift/acourse/controller/share"
	"github.com/acoshift/acourse/view"
)

func (c *ctrl) index(ctx *hime.Context) error {
	if ctx.URL.Path != "/" {
		return share.NotFound(ctx)
	}

	courses, err := c.Repository.ListPublicCourses(ctx)
	if err != nil {
		return err
	}

	p := view.Page(ctx)
	p.Data["Courses"] = courses
	return ctx.View("app.index", p)
}
