package auth

import (
	"github.com/moonrhythm/hime"

	"github.com/acoshift/acourse/internal/app/view"
	"github.com/acoshift/acourse/internal/pkg/bus"
	"github.com/acoshift/acourse/internal/pkg/context/appctx"
	"github.com/acoshift/acourse/internal/pkg/model/app"
	"github.com/acoshift/acourse/internal/pkg/model/auth"
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

	err := bus.Dispatch(ctx, &auth.SendPasswordResetEmail{Email: email})
	if app.IsUIError(err) {
		f.Add("Errors", err.Error())
		return ctx.RedirectToGet()
	}
	if err != nil {
		return err
	}

	return ctx.RedirectToGet()
}
