package auth

import (
	"net/url"

	"github.com/moonrhythm/hime"

	"github.com/acoshift/acourse/internal/context/appctx"
	"github.com/acoshift/acourse/internal/pkg/dispatcher"
	"github.com/acoshift/acourse/internal/pkg/model/app"
	"github.com/acoshift/acourse/internal/pkg/model/auth"
	"github.com/acoshift/acourse/internal/view"
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

	q := auth.SignInPassword{
		Email:    email,
		Password: pass,
	}
	err := dispatcher.Dispatch(ctx, &q)
	if app.IsUIError(err) {
		f.Add("Errors", err.Error())
		return ctx.RedirectToGet()
	}
	if err != nil {
		return err
	}

	appctx.RegenerateSessionID(ctx)
	appctx.SetUserID(ctx, q.Result)

	rd, _ := url.QueryUnescape(ctx.FormValue("r"))
	return ctx.SafeRedirect(rd)
}
