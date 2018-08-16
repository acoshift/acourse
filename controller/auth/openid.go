package auth

import (
	"net/http"

	"github.com/acoshift/hime"

	"github.com/acoshift/acourse/appsess"
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

	s := appctx.GetSession(ctx)
	appsess.SetOpenIDSessionID(s, state)
	return ctx.Redirect(redirect)
}

func (c *ctrl) openIDCallback(ctx *hime.Context) error {
	s := appctx.GetSession(ctx)
	sessID := appsess.GetOpenIDSessionID(s)
	appsess.DelOpenIDSessionID(s)

	userID, err := c.Service.SignInOpenIDCallback(ctx, ctx.Request().RequestURI, sessID)
	if service.IsUIError(err) {
		// TODO: redirect to sign in page
		return ctx.Status(http.StatusBadRequest).String(err.Error())
	}
	if err != nil {
		return err
	}

	s.Regenerate()
	appsess.SetUserID(s, userID)
	return ctx.RedirectTo("app.index")
}
