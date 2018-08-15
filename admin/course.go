package admin

import (
	"strconv"

	"github.com/acoshift/hime"
	"github.com/acoshift/paginate"

	"github.com/acoshift/acourse/repository"
	"github.com/acoshift/acourse/view"
)

func (c *ctrl) courses(ctx *hime.Context) error {
	cnt, err := repository.CountCourses(ctx)
	if err != nil {
		return err
	}

	pg, _ := strconv.ParseInt(ctx.FormValue("page"), 10, 64)
	pn := paginate.New(pg, 30, cnt)

	courses, err := repository.ListCourses(ctx, pn.Limit(), pn.Offset())
	if err != nil {
		return err
	}

	p := view.Page(ctx)
	p["Navbar"] = "admin.courses"
	p["Courses"] = courses
	p["Paginate"] = pn
	return ctx.View("admin.courses", p)
}
