package app

import (
	"github.com/acoshift/methodmux"
	"github.com/acoshift/prefixhandler"
	"github.com/moonrhythm/hime"
	"github.com/moonrhythm/httpmux"
)

// Mount mounts app handlers
func Mount(m *httpmux.Mux) {
	m.Handle("/", methodmux.Get(
		hime.Handler(index),
	))
	m.Handle("/signout", methodmux.Post(
		hime.Handler(signOut),
	))

	profile := m.Group("/profile", mustSignedIn)
	profile.Handle("/", methodmux.Get(
		hime.Handler(getProfile),
	))
	profile.Handle("/edit", methodmux.GetPost(
		hime.Handler(getProfileEdit),
		hime.Handler(postProfileEdit),
	))

	// course
	m.Handle("/course/", prefixhandler.New("/course", courseIDKey{}, newCourseHandler()))
}
