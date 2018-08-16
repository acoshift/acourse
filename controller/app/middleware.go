package app

import (
	"net/http"
	"net/url"

	"github.com/acoshift/acourse/appsess"
	"github.com/acoshift/acourse/context/appctx"
)

func mustSignedIn(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := appctx.GetSession(r.Context())
		id := appsess.GetUserID(s)
		if len(id) == 0 {
			http.Redirect(w, r, "/signin?r="+url.QueryEscape(r.RequestURI), http.StatusFound)
			return
		}
		h.ServeHTTP(w, r)
	})
}
