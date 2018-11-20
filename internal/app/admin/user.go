package admin

import (
	"strconv"

	"github.com/acoshift/paginate"
	"github.com/moonrhythm/hime"

	"github.com/acoshift/acourse/internal/app/view"
	"github.com/acoshift/acourse/internal/pkg/bus"
	"github.com/acoshift/acourse/internal/pkg/model/admin"
)

func getUsers(ctx *hime.Context) error {
	cnt := admin.CountUsers{}
	err := bus.Dispatch(ctx, &cnt)
	if err != nil {
		return err
	}

	pg, _ := strconv.ParseInt(ctx.FormValue("page"), 10, 64)
	pn := paginate.New(pg, 30, cnt.Result)

	list := admin.ListUsers{Limit: pn.Limit(), Offset: pn.Offset()}
	err = bus.Dispatch(ctx, &list)
	if err != nil {
		return err
	}

	p := view.Page(ctx)
	p.Data["Navbar"] = "admin.users"
	p.Data["Users"] = list.Result
	p.Data["Paginate"] = pn
	return ctx.View("admin.users", p)
}
