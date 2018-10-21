package admin

import (
	"strconv"

	"github.com/acoshift/paginate"
	"github.com/moonrhythm/hime"

	"github.com/acoshift/acourse/view"
)

func (c *ctrl) users(ctx *hime.Context) error {
	cnt, err := countUsers(ctx)
	if err != nil {
		return err
	}

	pg, _ := strconv.ParseInt(ctx.FormValue("page"), 10, 64)
	pn := paginate.New(pg, 30, cnt)

	users, err := listUsers(ctx, pn.Limit(), pn.Offset())
	if err != nil {
		return err
	}

	p := view.Page(ctx)
	p.Data["Navbar"] = "admin.users"
	p.Data["Users"] = users
	p.Data["Paginate"] = pn
	return ctx.View("admin.users", p)
}
