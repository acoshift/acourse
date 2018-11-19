package app

import (
	"net/http"

	"github.com/acoshift/methodmux"
	"github.com/acoshift/prefixhandler"
	"github.com/moonrhythm/hime"
)

// Mount mounts app handlers
func Mount(m *http.ServeMux) {
	m.Handle("/", methodmux.Get(
		hime.Handler(index),
	))
	m.Handle("/signout", methodmux.Post(
		hime.Handler(signOut),
	))

	// profile
	m.Handle("/profile", mustSignedIn(methodmux.Get(
		hime.Handler(profile),
	)))
	m.Handle("/profile/edit", mustSignedIn(methodmux.GetPost(
		hime.Handler(profileEdit),
		hime.Handler(postProfileEdit),
	)))

	// course
	m.Handle("/course/", prefixhandler.New("/course", courseIDKey{}, newCourseHandler()))
}
