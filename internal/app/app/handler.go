package app

import (
	"net/http"

	"github.com/acoshift/methodmux"
	"github.com/acoshift/prefixhandler"
	"github.com/moonrhythm/hime"
)

// Mount mounts app handlers
func Mount(m *http.ServeMux, baseURL string) {
	c := &ctrl{baseURL}

	m.Handle("/", methodmux.Get(
		hime.Handler(c.index),
	))
	m.Handle("/signout", methodmux.Post(
		hime.Handler(c.signOut),
	))

	// profile
	m.Handle("/profile", mustSignedIn(methodmux.Get(
		hime.Handler(c.profile),
	)))
	m.Handle("/profile/edit", mustSignedIn(methodmux.GetPost(
		hime.Handler(c.profileEdit),
		hime.Handler(c.postProfileEdit),
	)))

	// course
	m.Handle("/course/", prefixhandler.New("/course", courseIDKey{}, newCourseHandler(c)))
}

type ctrl struct {
	BaseURL string
}
