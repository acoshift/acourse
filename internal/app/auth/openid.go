package auth

import (
	"net/http"

	"github.com/moonrhythm/hime"

	"github.com/acoshift/acourse/internal/pkg/auth"
	"github.com/acoshift/acourse/internal/pkg/context/appctx"
)

func getOpenID(ctx *hime.Context) error {
	p := ctx.FormValue("p")

	redirectURI, state, err := auth.GenerateOpenIDURI(ctx, p)
	if err == auth.ErrInvalidProvider {
		// TODO: redirect to sign in page
		return ctx.Status(http.StatusBadRequest).String("Provider not allowed")
	}
	if err != nil {
		return err
	}

	appctx.SetOpenIDState(ctx, state)
	return ctx.Redirect(redirectURI)
}

func getOpenIDCallback(ctx *hime.Context) error {
	sessID := appctx.GetOpenIDState(ctx)
	appctx.DelOpenIDState(ctx)

	userID, err := auth.SignInOpenIDCallback(ctx, ctx.RequestURI, sessID)
	if err == auth.ErrInvalidCallbackURI {
		// TODO: redirect to sign in page
		return ctx.Status(http.StatusBadRequest).String("Invalid callback uri")
	}
	if err == auth.ErrEmailNotAvailable {
		// TODO: redirect to sign in page
		return ctx.Status(http.StatusBadRequest).String("อีเมลนี้ถูกสมัครแล้ว")
	}
	if err == auth.ErrUsernameNotAvailable {
		// TODO: redirect to sign in page
		return ctx.Status(http.StatusBadRequest).String("username นี้ถูกใช้งานแล้ว")
	}
	if err != nil {
		return err
	}

	appctx.RegenerateSessionID(ctx)
	appctx.SetUserID(ctx, userID)
	return ctx.RedirectTo("app.index")
}
