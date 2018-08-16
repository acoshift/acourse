package auth

import (
	"github.com/acoshift/hime"

	"github.com/acoshift/acourse/context/appctx"
	"github.com/acoshift/acourse/service"
	"github.com/acoshift/acourse/view"
)

func (c *ctrl) signIn(ctx *hime.Context) error {
	return ctx.View("auth.signin", view.Page(ctx))
}

func (c *ctrl) postSignIn(ctx *hime.Context) error {
	f := appctx.GetFlash(ctx)

	email := ctx.FormValueTrimSpace("email")
	if email == "" {
		f.Add("Errors", "email required")
		return ctx.RedirectToGet()
	}

	err := c.Service.SendSignInMagicLinkEmail(ctx, email)
	if service.IsUIError(err) {
		f.Set("Email", email)
		f.Add("Errors", err.Error())
		return ctx.RedirectToGet()
	}
	if err != nil {
		return err
	}

	f.Set("CheckEmail", true)

	return ctx.RedirectTo("auth.signin.check-email")
}

func (c *ctrl) checkEmail(ctx *hime.Context) error {
	f := appctx.GetFlash(ctx)
	if f.GetBool("CheckEmail") {
		return ctx.Redirect("/")
	}
	return ctx.View("auth.check-email", view.Page(ctx))
}

func (c *ctrl) signInLink(ctx *hime.Context) error {
	linkID := ctx.FormValue("id")
	if len(linkID) == 0 {
		return ctx.RedirectTo("auth.signin")
	}

	f := appctx.GetFlash(ctx)

	userID, err := c.Service.SignInMagicLink(ctx, linkID)
	if service.IsUIError(err) {
		f.Add("Errors", err.Error())
		return ctx.RedirectTo("auth.signin")
	}
	if err != nil {
		return err
	}

	appctx.SetUserID(ctx, userID)
	return ctx.Redirect("/")
}

func (c *ctrl) signInPassword(ctx *hime.Context) error {
	return ctx.View("auth.signin-password", view.Page(ctx))
}

func (c *ctrl) postSignInPassword(ctx *hime.Context) error {
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

	return ctx.SafeRedirect(ctx.FormValue("r"))
}
