package auth

import (
	"github.com/acoshift/hime"

	"github.com/acoshift/acourse/appsess"
	"github.com/acoshift/acourse/context/appctx"
	"github.com/acoshift/acourse/service"
	"github.com/acoshift/acourse/view"
)

func (c *ctrl) signUp(ctx *hime.Context) error {
	return ctx.View("auth.signup", view.Page(ctx))
}

func (c *ctrl) postSignUp(ctx *hime.Context) error {
	f := appctx.GetSession(ctx).Flash()

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

	userID, err := c.Service.SignUp(ctx, email, pass)
	if service.IsUIError(err) {
		f.Set("Email", email)
		f.Add("Errors", err.Error())
		return ctx.RedirectToGet()
	}
	if err != nil {
		return err
	}

	s := appctx.GetSession(ctx)
	appsess.SetUserID(s, userID)

	return ctx.SafeRedirect(ctx.FormValue("r"))
}

func (c *ctrl) resetPassword(ctx *hime.Context) error {
	return ctx.View("auth.reset-password", view.Page(ctx))
}

func (c *ctrl) postResetPassword(ctx *hime.Context) error {
	f := appctx.GetSession(ctx).Flash()
	f.Set("OK", "1")

	email := ctx.PostFormValueTrimSpace("email")
	if email == "" {
		f.Add("Errors", "email required")
		return ctx.RedirectToGet()
	}

	err := c.Service.SendPasswordResetEmail(ctx, email)
	if service.IsUIError(err) {
		f.Add("Errors", err.Error())
		return ctx.RedirectToGet()
	}
	if err != nil {
		return err
	}

	return ctx.RedirectToGet()
}
