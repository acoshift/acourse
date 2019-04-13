package auth

import (
	"github.com/moonrhythm/hime"

	"github.com/acoshift/acourse/internal/app/view"
	"github.com/acoshift/acourse/internal/pkg/auth"
	"github.com/acoshift/acourse/internal/pkg/context/appctx"
)

func getResetPassword(ctx *hime.Context) error {
	return ctx.View("auth.reset-password", view.Page(ctx))
}

func postResetPassword(ctx *hime.Context) error {
	f := appctx.GetFlash(ctx)
	f.Set("OK", "1")

	email := ctx.PostFormValueTrimSpace("email")
	if email == "" {
		f.Add("Errors", "email required")
		return ctx.RedirectToGet()
	}

	err := auth.SendPasswordResetEmail(ctx, email)
	if err != nil {
		f.Add("Errors", err.Error())
		return ctx.RedirectToGet()
	}

	return ctx.RedirectToGet()
}
