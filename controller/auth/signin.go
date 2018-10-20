package auth

import (
	"net/url"

	"github.com/moonrhythm/hime"

	"github.com/acoshift/acourse/context/appctx"
	"github.com/acoshift/acourse/service"
	"github.com/acoshift/acourse/view"
)

func (c *ctrl) signIn(ctx *hime.Context) error {
	return ctx.View("auth.signin", view.Page(ctx))
}

func (c *ctrl) postSignIn(ctx *hime.Context) error {
	f := appctx.GetFlash(ctx)

	email := ctx.PostFormValue("email")
	if email == "" {
		f.Add("Errors", "email required")
	}
	pass := ctx.PostFormValue("password")
	if pass == "" {
		f.Add("Errors", "password required")
	}
	if f.Has("Errors") {
		f.Set("Email", email)
		return ctx.RedirectToGet()
	}

	userID, err := c.Service.SignInPassword(ctx, email, pass)
	if service.IsUIError(err) {
		f.Add("Errors", err.Error())
		return ctx.RedirectToGet()
	}
	if err != nil {
		return err
	}

	appctx.RegenerateSessionID(ctx)
	appctx.SetUserID(ctx, userID)

	rd, _ := url.QueryUnescape(ctx.FormValue("r"))
	return ctx.SafeRedirect(rd)
}
