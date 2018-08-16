package app

import (
	"net/http"
	"net/url"

	"github.com/acoshift/hime"

	"github.com/acoshift/acourse/appsess"
	"github.com/acoshift/acourse/context/appctx"
)

func mustSignedIn(h http.Handler) http.Handler {
	return hime.Handler(func(ctx *hime.Context) error {
		s := appctx.GetSession(ctx)
		id := appsess.GetUserID(s)
		if len(id) == 0 {
			return ctx.RedirectTo("auth.signin", ctx.Param("r", url.QueryEscape(ctx.Request().RequestURI)))
		}
		return ctx.Handle(h)
	})
}
