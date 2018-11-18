package app

import (
	"net/http"
	"net/url"

	"github.com/moonrhythm/hime"

	"github.com/acoshift/acourse/internal/pkg/context/appctx"
)

func mustSignedIn(h http.Handler) http.Handler {
	return hime.Handler(func(ctx *hime.Context) error {
		id := appctx.GetUserID(ctx)
		if len(id) == 0 {
			return ctx.RedirectTo("auth.signin", ctx.Param("r", url.QueryEscape(ctx.RequestURI)))
		}
		return ctx.Handle(h)
	})
}
