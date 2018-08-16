package auth

import (
	"net/http"

	"github.com/acoshift/hime"

	"github.com/acoshift/acourse/context/appctx"
	"github.com/acoshift/acourse/service"
)

var allowProvider = map[string]bool{
	"google.com": true,
	"github.com": true,
}

func (c *ctrl) openID(ctx *hime.Context) error {
	p := ctx.FormValue("p")

	redirect, state, err := c.Service.GenerateOpenIDURI(ctx, p)
	if service.IsUIError(err) {
		// TODO: redirect to sign in page
		return ctx.Status(http.StatusBadRequest).String(err.Error())
	}
	if err != nil {
		return err
	}

	appctx.SetOpenIDState(ctx, state)
	return ctx.Redirect(redirect)
}

func (c *ctrl) openIDCallback(ctx *hime.Context) error {
	sessID := appctx.GetOpenIDState(ctx)
	appctx.DelOpenIDState(ctx)

	userID, err := c.Service.SignInOpenIDCallback(ctx, ctx.Request().RequestURI, sessID)
	if service.IsUIError(err) {
		// TODO: redirect to sign in page
		return ctx.Status(http.StatusBadRequest).String(err.Error())
	}
	if err != nil {
		return err
	}

	appctx.RegenerateSessionID(ctx)
	appctx.SetUserID(ctx, userID)
	return ctx.RedirectTo("app.index")
}
