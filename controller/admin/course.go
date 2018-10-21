package admin

import (
	"strconv"

	"github.com/acoshift/paginate"
	"github.com/moonrhythm/hime"

	"github.com/acoshift/acourse/view"
)

func (c *ctrl) courses(ctx *hime.Context) error {
	cnt, err := countCourses(ctx)
	if err != nil {
		return err
	}

	pg, _ := strconv.ParseInt(ctx.FormValue("page"), 10, 64)
	pn := paginate.New(pg, 30, cnt)

	courses, err := listCourses(ctx, pn.Limit(), pn.Offset())
	if err != nil {
		return err
	}

	p := view.Page(ctx)
	p.Data["Navbar"] = "admin.courses"
	p.Data["Courses"] = courses
	p.Data["Paginate"] = pn
	return ctx.View("admin.courses", p)
}
