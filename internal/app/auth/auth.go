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

func (c *ctrl) signUp(ctx *hime.Context) error {
	return ctx.View("auth.signup", view.Page(ctx))
}

func (c *ctrl) postSignUp(ctx *hime.Context) error {
	f := appctx.GetFlash(ctx)

	email := ctx.PostFormValueTrimSpace("email")
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

	q := auth.SignUp{
		Email:    email,
		Password: pass,
	}
	err := dispatcher.Dispatch(ctx, &q)
	if app.IsUIError(err) {
		f.Set("Email", email)
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

func (c *ctrl) resetPassword(ctx *hime.Context) error {
	return ctx.View("auth.reset-password", view.Page(ctx))
}

func (c *ctrl) postResetPassword(ctx *hime.Context) error {
	f := appctx.GetFlash(ctx)
	f.Set("OK", "1")

	email := ctx.PostFormValueTrimSpace("email")
	if email == "" {
		f.Add("Errors", "email required")
		return ctx.RedirectToGet()
	}

	err := dispatcher.Dispatch(ctx, &auth.SendPasswordResetEmail{Email: email})
	if app.IsUIError(err) {
		f.Add("Errors", err.Error())
		return ctx.RedirectToGet()
	}
	if err != nil {
		return err
	}

	return ctx.RedirectToGet()
}
