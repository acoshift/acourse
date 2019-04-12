package admin

import (
	"strconv"

	"github.com/acoshift/paginate"
	"github.com/moonrhythm/hime"

	"github.com/acoshift/acourse/internal/app/view"
	"github.com/acoshift/acourse/internal/pkg/admin"
)

func getCourses(ctx *hime.Context) error {
	cnt, err := admin.CountCourses(ctx)
	if err != nil {
		return err
	}

	pg, _ := strconv.ParseInt(ctx.FormValue("page"), 10, 64)
	pn := paginate.New(pg, 30, cnt)

	list, err := admin.GetCourses(ctx, pn.Limit(), pn.Offset())
	if err != nil {
		return err
	}

	p := view.Page(ctx)
	p.Data["Navbar"] = "admin.courses"
	p.Data["Courses"] = list
	p.Data["Paginate"] = pn
	return ctx.View("admin.courses", p)
}
