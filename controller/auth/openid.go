package auth

import (
	"net/http"

	"github.com/moonrhythm/dispatcher"
	"github.com/moonrhythm/hime"

	"github.com/acoshift/acourse/context/appctx"
	"github.com/acoshift/acourse/model/auth"
	"github.com/acoshift/acourse/service"
)

var allowProvider = map[string]bool{
	"google.com": true,
	"github.com": true,
}

func (c *ctrl) openID(ctx *hime.Context) error {
	p := ctx.FormValue("p")

	q := auth.GenerateOpenIDURI{Provider: p}
	err := dispatcher.Dispatch(ctx, &q)
	if service.IsUIError(err) {
		// TODO: redirect to sign in page
		return ctx.Status(http.StatusBadRequest).String(err.Error())
	}
	if err != nil {
		return err
	}

	appctx.SetOpenIDState(ctx, q.Result.State)
	return ctx.Redirect(q.Result.RedirectURI)
}

func (c *ctrl) openIDCallback(ctx *hime.Context) error {
	sessID := appctx.GetOpenIDState(ctx)
	appctx.DelOpenIDState(ctx)

	q := auth.SignInOpenIDCallback{URI: ctx.RequestURI, State: sessID}
	err := dispatcher.Dispatch(ctx, &q)
	if service.IsUIError(err) {
		// TODO: redirect to sign in page
		return ctx.Status(http.StatusBadRequest).String(err.Error())
	}
	if err != nil {
		return err
	}

	appctx.RegenerateSessionID(ctx)
	appctx.SetUserID(ctx, q.Result)
	return ctx.RedirectTo("app.index")
}
