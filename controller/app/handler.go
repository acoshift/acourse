package app

import (
	"net/http"

	"github.com/acoshift/hime"
	"github.com/acoshift/methodmux"
	"github.com/acoshift/prefixhandler"

	"github.com/acoshift/acourse/service"
)

// Config is the app config
type Config struct {
	BaseURL string
	Service service.Service
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
	{
		m := http.NewServeMux()
		m.Handle("/", methodmux.Get(
			hime.Handler(c.courseView),
		))
		m.Handle("/content", mustSignedIn(methodmux.Get(
			hime.Handler(c.courseContent),
		)))
		m.Handle("/enroll", mustSignedIn(methodmux.GetPost(
			hime.Handler(c.courseEnroll),
			hime.Handler(c.postCourseEnroll),
		)))
		m.Handle("/assignment", mustSignedIn(methodmux.Get(
			hime.Handler(c.courseAssignment),
		)))

		mux.Handle("/course/", prefixhandler.New("/course", courseURLKey{}, m))
	}

	return mux
}

type ctrl struct {
	Config
}
