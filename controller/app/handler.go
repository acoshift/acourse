package app

import (
	"net/http"

	"github.com/acoshift/methodmux"
	"github.com/acoshift/prefixhandler"
	"github.com/moonrhythm/hime"
)

// Config is the app config
type Config struct {
	BaseURL    string
	Repository Repository
}

// New creates new app
func New(cfg Config) http.Handler {
	c := &ctrl{cfg}

	mux := http.NewServeMux()

	mux.Handle("/", methodmux.Get(
		hime.Handler(c.index),
	))
	mux.Handle("/signout", methodmux.Post(
		hime.Handler(c.signOut),
	))

	// profile
	mux.Handle("/profile", mustSignedIn(methodmux.Get(
		hime.Handler(c.profile),
	)))
	mux.Handle("/profile/edit", mustSignedIn(methodmux.GetPost(
		hime.Handler(c.profileEdit),
		hime.Handler(c.postProfileEdit),
	)))

	// course
	mux.Handle("/course/", prefixhandler.New("/course", courseIDKey{}, newCourseHandler(c)))

	return mux
}

type ctrl struct {
	Config
}
