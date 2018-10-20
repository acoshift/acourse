package auth

import (
	"net/http"

	"github.com/acoshift/methodmux"
	"github.com/moonrhythm/hime"
)

// New creates new auth handler
func New() http.Handler {
	c := &ctrl{}

	mux := http.NewServeMux()

	mux.Handle("/signin", methodmux.GetPost(
		hime.Handler(c.signIn),
		hime.Handler(c.postSignIn),
	))
	mux.Handle("/reset/password", methodmux.GetPost(
		hime.Handler(c.resetPassword),
		hime.Handler(c.postResetPassword),
	))
	mux.Handle("/openid", methodmux.Get(
		hime.Handler(c.openID),
	))
	mux.Handle("/openid/callback", methodmux.Get(
		hime.Handler(c.openIDCallback),
	))
	mux.Handle("/signup", methodmux.GetPost(
		hime.Handler(c.signUp),
		hime.Handler(c.postSignUp),
	))

	return mux
}

type ctrl struct{}
