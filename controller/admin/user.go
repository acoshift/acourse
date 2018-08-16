package admin

import (
	"strconv"

	"github.com/acoshift/hime"
	"github.com/acoshift/paginate"

	"github.com/acoshift/acourse/repository"
	"github.com/acoshift/acourse/view"
)

func (c *ctrl) users(ctx *hime.Context) error {
	cnt, err := repository.CountUsers(ctx)
	if err != nil {
		return err
	}

	pg, _ := strconv.ParseInt(ctx.FormValue("page"), 10, 64)
	pn := paginate.New(pg, 30, cnt)

	users, err := repository.ListUsers(ctx, pn.Limit(), pn.Offset())
	if err != nil {
		return err
	}

	p := view.Page(ctx)
	p["Navbar"] = "admin.users"
	p["Users"] = users
	p["Paginate"] = pn
	return ctx.View("admin.users", p)
}
