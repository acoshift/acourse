package auth

import (
	"net/http"

	"github.com/acoshift/methodmux"
	"github.com/moonrhythm/hime"
	"github.com/moonrhythm/httpmux"

	"github.com/acoshift/acourse/internal/pkg/context/appctx"
)

// Mount mounts auth handlers
func Mount(m *httpmux.Mux) {
	mux := m.Group("/auth", notSignedIn)
	mux.Handle("/signin", methodmux.GetPost(
		hime.Handler(getSignIn),
		hime.Handler(postSignIn),
	))
	mux.Handle("/reset/password", methodmux.GetPost(
		hime.Handler(getResetPassword),
		hime.Handler(postResetPassword),
	))
	mux.Handle("/openid", methodmux.Get(
		hime.Handler(getOpenID),
	))
	mux.Handle("/openid/callback", methodmux.Get(
		hime.Handler(getOpenIDCallback),
	))
	mux.Handle("/signup", methodmux.GetPost(
		hime.Handler(getSignUp),
		hime.Handler(postSignUp),
	))
}

func notSignedIn(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := appctx.GetUserID(r.Context())
		if len(id) > 0 {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		h.ServeHTTP(w, r)
	})
}
