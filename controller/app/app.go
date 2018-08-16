package app

import (
	"math/rand"
	"net/http"

	"github.com/acoshift/go-firebase-admin"
	"github.com/acoshift/header"
	"github.com/acoshift/hime"
	"github.com/acoshift/methodmux"
	"github.com/acoshift/prefixhandler"

	"github.com/acoshift/acourse/file"
	"github.com/acoshift/acourse/image"
	"github.com/acoshift/acourse/notify"
	"github.com/acoshift/acourse/view"
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

	methodmux.FallbackHandler = hime.Handler(notFound)

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

	// editor
	{
		m := http.NewServeMux()
		m.Handle("/create", onlyInstructor(methodmux.GetPost(
			hime.Handler(editorCreate),
			hime.Handler(postEditorCreate),
		)))
		m.Handle("/course", isCourseOwner(methodmux.GetPost(
			hime.Handler(editorCourse),
			hime.Handler(postEditorCourse),
		)))
		m.Handle("/content", isCourseOwner(methodmux.GetPost(
			hime.Handler(editorContent),
			hime.Handler(postEditorContent),
		)))
		m.Handle("/content/create", isCourseOwner(methodmux.GetPost(
			hime.Handler(editorContentCreate),
			hime.Handler(postEditorContentCreate),
		)))
		// TODO: add middleware ?
		m.Handle("/content/edit", methodmux.GetPost(
			hime.Handler(editorContentEdit),
			hime.Handler(postEditorContentEdit),
		))

		mux.Handle("/editor/", http.StripPrefix("/editor", m))
	}

	return mux
}

var notFoundImages = []string{
	"https://storage.googleapis.com/acourse/static/9961f3c1-575f-4b98-af4f-447566ee1cb3.png",
	"https://storage.googleapis.com/acourse/static/b14a40c9-d3a4-465d-9453-ce7fcfbc594c.png",
}

func notFound(ctx *hime.Context) error {
	p := view.Page(ctx)
	p["Image"] = notFoundImages[rand.Intn(len(notFoundImages))]
	ctx.ResponseWriter().Header().Set(header.XContentTypeOptions, "nosniff")
	return ctx.Status(http.StatusNotFound).View("error.not-found", p)
}
