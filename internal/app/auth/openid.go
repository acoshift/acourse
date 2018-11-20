package auth

import (
	"net/http"

	"github.com/moonrhythm/hime"

	"github.com/acoshift/acourse/internal/pkg/bus"
	"github.com/acoshift/acourse/internal/pkg/context/appctx"
	"github.com/acoshift/acourse/internal/pkg/model/app"
	"github.com/acoshift/acourse/internal/pkg/model/auth"
)

var allowProvider = map[string]bool{
	"google.com": true,
	"github.com": true,
}

func getOpenID(ctx *hime.Context) error {
	p := ctx.FormValue("p")

	q := auth.GenerateOpenIDURI{Provider: p}
	err := bus.Dispatch(ctx, &q)
	if app.IsUIError(err) {
		// TODO: redirect to sign in page
		return ctx.Status(http.StatusBadRequest).String(err.Error())
	}
	if err != nil {
		return err
	}

	appctx.SetOpenIDState(ctx, q.Result.State)
	return ctx.Redirect(q.Result.RedirectURI)
}

func getOpenIDCallback(ctx *hime.Context) error {
	sessID := appctx.GetOpenIDState(ctx)
	appctx.DelOpenIDState(ctx)

	q := auth.SignInOpenIDCallback{URI: ctx.RequestURI, State: sessID}
	err := bus.Dispatch(ctx, &q)
	if app.IsUIError(err) {
		// TODO: redirect to sign in page
		return ctx.Status(http.StatusBadRequest).String(err.Error())
	}
	if err != nil {
		return err
	}

	appctx.RegenerateSessionID(ctx)
	appctx.SetUserID(ctx, q.Result)
	return ctx.RedirectTo("app.index")
}
