package auth

import (
	"net/http"

	"github.com/acoshift/go-firebase-admin"
	"github.com/acoshift/hime"
	"github.com/acoshift/methodmux"

	"github.com/acoshift/acourse/email"
	"github.com/acoshift/acourse/file"
	"github.com/acoshift/acourse/image"
	"github.com/acoshift/acourse/notify"
)

// Config is auth config
type Config struct {
	Auth               *firebase.Auth
	EmailSender        email.Sender
	AdminNotifier      notify.AdminNotifier
	FileStorage        file.Storage
	ImageResizeEncoder image.JPEGResizeEncoder
	BaseURL            string
}

// New creates new auth handler
func New(cfg Config) http.Handler {
	c := &ctrl{cfg}

	mux := http.NewServeMux()

	mux.Handle("/signin", methodmux.GetPost(
		hime.Handler(c.signIn),
		hime.Handler(c.postSignIn),
	))
	mux.Handle("/signin/password", methodmux.GetPost(
		hime.Handler(c.signInPassword),
		hime.Handler(c.postSignInPassword),
	))
	mux.Handle("/signin/check-email", methodmux.Get(
		hime.Handler(c.checkEmail),
	))
	mux.Handle("/signin/link", methodmux.Get(
		hime.Handler(c.signInLink),
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

type ctrl struct {
	Config
}
