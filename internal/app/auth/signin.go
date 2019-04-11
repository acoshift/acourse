package auth

import (
	"net/url"

	"github.com/moonrhythm/hime"

	"github.com/acoshift/acourse/internal/app/view"
	app2 "github.com/acoshift/acourse/internal/pkg/app"
	"github.com/acoshift/acourse/internal/pkg/auth"
	"github.com/acoshift/acourse/internal/pkg/context/appctx"
)

func getSignIn(ctx *hime.Context) error {
	return ctx.View("auth.signin", view.Page(ctx))
}

func postSignIn(ctx *hime.Context) error {
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

	userID, err := auth.SignInPassword(ctx, email, pass)
	if app2.IsUIError(err) {
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
