package auth

import (
	"net/http"

	"github.com/moonrhythm/hime"

	"github.com/acoshift/acourse/internal/pkg/app"
	"github.com/acoshift/acourse/internal/pkg/auth"
	"github.com/acoshift/acourse/internal/pkg/context/appctx"
)

var allowProvider = map[string]bool{
	"google.com": true,
	"github.com": true,
}

func getOpenID(ctx *hime.Context) error {
	p := ctx.FormValue("p")

	redirectURI, state, err := auth.GenerateOpenIDURI(ctx, p)
	if app.IsUIError(err) {
		// TODO: redirect to sign in page
		return ctx.Status(http.StatusBadRequest).String(err.Error())
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
	if app.IsUIError(err) {
		// TODO: redirect to sign in page
		return ctx.Status(http.StatusBadRequest).String(err.Error())
	}
	if err != nil {
		return err
	}

	appctx.RegenerateSessionID(ctx)
	appctx.SetUserID(ctx, userID)
	return ctx.RedirectTo("app.index")
}
