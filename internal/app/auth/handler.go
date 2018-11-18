package auth

import (
	"net/http"

	"github.com/acoshift/methodmux"
	"github.com/moonrhythm/hime"

	"github.com/acoshift/acourse/internal/pkg/context/appctx"
)

// Mount mounts auth handlers
func Mount(m *http.ServeMux) {
	mux := http.NewServeMux()
	mux.Handle("/auth/signin", methodmux.GetPost(
		hime.Handler(getSignIn),
		hime.Handler(postSignIn),
	))
	mux.Handle("/auth/reset/password", methodmux.GetPost(
		hime.Handler(getResetPassword),
		hime.Handler(postResetPassword),
	))
	mux.Handle("/auth/openid", methodmux.Get(
		hime.Handler(getOpenID),
	))
	mux.Handle("/auth/openid/callback", methodmux.Get(
		hime.Handler(getOpenIDCallback),
	))
	mux.Handle("/auth/signup", methodmux.GetPost(
		hime.Handler(getSignUp),
		hime.Handler(postSignUp),
	))

	m.Handle("/auth/", notSignedIn(mux))
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
