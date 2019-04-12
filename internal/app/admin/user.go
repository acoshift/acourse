package admin

import (
	"strconv"

	"github.com/acoshift/paginate"
	"github.com/moonrhythm/hime"

	"github.com/acoshift/acourse/internal/app/view"
	"github.com/acoshift/acourse/internal/pkg/admin"
)

func getUsers(ctx *hime.Context) error {
	cnt, err := admin.CountUsers(ctx)
	if err != nil {
		return err
	}

	pg, _ := strconv.ParseInt(ctx.FormValue("page"), 10, 64)
	pn := paginate.New(pg, 30, cnt)

	list, err := admin.GetUsers(ctx, pn.Limit(), pn.Offset())
	if err != nil {
		return err
	}

	p := view.Page(ctx)
	p.Data["Navbar"] = "admin.users"
	p.Data["Users"] = list
	p.Data["Paginate"] = pn
	return ctx.View("admin.users", p)
}
