package app

import (
	"net/http"

	"github.com/acoshift/go-firebase-admin"
	"github.com/acoshift/hime"
	"github.com/acoshift/methodmux"
	"github.com/acoshift/prefixhandler"

	"github.com/acoshift/acourse/file"
	"github.com/acoshift/acourse/image"
	"github.com/acoshift/acourse/notify"
)

var (
	auth               *firebase.Auth
	adminNotifier      notify.AdminNotifier
	baseURL            string
	fileStorage        file.Storage
	imageResizeEncoder image.JPEGResizeEncoder
)

// New creates new app
func New(config Config) http.Handler {
	auth = config.Auth
	adminNotifier = config.AdminNotifier
	baseURL = config.BaseURL
	fileStorage = config.FileStorage
	imageResizeEncoder = config.ImageResizeEncoder

	mux := http.NewServeMux()

	mux.Handle("/", methodmux.Get(
		hime.Handler(index),
	))
	mux.Handle("/signout", methodmux.Post(
		hime.Handler(signOut),
	))

	// profile
	mux.Handle("/profile", mustSignedIn(methodmux.Get(
		hime.Handler(profile),
	)))
	mux.Handle("/profile/edit", mustSignedIn(methodmux.GetPost(
		hime.Handler(profileEdit),
		hime.Handler(postProfileEdit),
	)))

	// course
	{
		m := http.NewServeMux()
		m.Handle("/", methodmux.Get(
			hime.Handler(courseView),
		))
		m.Handle("/content", mustSignedIn(methodmux.Get(
			hime.Handler(courseContent),
		)))
		m.Handle("/enroll", mustSignedIn(methodmux.GetPost(
			hime.Handler(courseEnroll),
			hime.Handler(postCourseEnroll),
		)))
		m.Handle("/assignment", mustSignedIn(methodmux.Get(
			hime.Handler(courseAssignment),
		)))

		mux.Handle("/course/", prefixhandler.New("/course", courseURLKey{}, m))
	}

	return mux
}
