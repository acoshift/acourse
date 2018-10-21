package admin

import (
	"strconv"

	"github.com/acoshift/paginate"
	"github.com/moonrhythm/dispatcher"
	"github.com/moonrhythm/hime"

	"github.com/acoshift/acourse/model/admin"
	"github.com/acoshift/acourse/view"
)

func (c *ctrl) users(ctx *hime.Context) error {
	cnt := admin.CountUsers{}
	err := dispatcher.Dispatch(ctx, &cnt)
	if err != nil {
		return err
	}

	pg, _ := strconv.ParseInt(ctx.FormValue("page"), 10, 64)
	pn := paginate.New(pg, 30, cnt.Result)

	list := admin.ListUsers{Limit: pn.Limit(), Offset: pn.Offset()}
	err = dispatcher.Dispatch(ctx, &list)
	if err != nil {
		return err
	}

	p := view.Page(ctx)
	p.Data["Navbar"] = "admin.users"
	p.Data["Users"] = list.Result
	p.Data["Paginate"] = pn
	return ctx.View("admin.users", p)
}
