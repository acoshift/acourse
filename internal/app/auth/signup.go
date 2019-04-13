package auth

import (
	"net/url"

	"github.com/moonrhythm/hime"

	"github.com/acoshift/acourse/internal/app/view"
	"github.com/acoshift/acourse/internal/pkg/auth"
	"github.com/acoshift/acourse/internal/pkg/context/appctx"
)

func getSignUp(ctx *hime.Context) error {
	return ctx.View("auth.signup", view.Page(ctx))
}

func postSignUp(ctx *hime.Context) error {
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

	f.Set("Email", email)
	userID, err := auth.SignUp(ctx, email, pass)
	if err == auth.ErrEmailNotAvailable {
		f.Add("Errors", "อีเมลนี้ถูกสมัครแล้ว")
		return ctx.RedirectToGet()
	}
	if err == auth.ErrUsernameNotAvailable {
		f.Add("Errors", "username นี้ถูกใช้งานแล้ว")
		return ctx.RedirectToGet()
	}
	if err != nil {
		f.Add("Errors", err.Error())
		return ctx.RedirectToGet()
	}

	f.Del("Email")

	appctx.RegenerateSessionID(ctx)
	appctx.SetUserID(ctx, userID)

	rd, _ := url.QueryUnescape(ctx.FormValue("r"))
	return ctx.SafeRedirect(rd)
}
