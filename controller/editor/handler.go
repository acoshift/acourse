package editor

import (
	"net/http"

	"github.com/acoshift/methodmux"
	"github.com/moonrhythm/hime"
)

// Config is editor config
type Config struct {
	Repository Repository
}

// New creates new editor handler
func New(cfg Config) http.Handler {
	c := &ctrl{cfg}

	mux := http.NewServeMux()
	mux.Handle("/course/create", c.onlyInstructor(methodmux.GetPost(
		hime.Handler(c.courseCreate),
		hime.Handler(c.postCourseCreate),
	)))
	mux.Handle("/course/edit", c.isCourseOwner(methodmux.GetPost(
		hime.Handler(c.courseEdit),
		hime.Handler(c.postCourseEdit),
	)))
	mux.Handle("/content", c.isCourseOwner(methodmux.GetPost(
		hime.Handler(c.contentList),
		hime.Handler(c.postContentList),
	)))
	mux.Handle("/content/create", c.isCourseOwner(methodmux.GetPost(
		hime.Handler(c.contentCreate),
		hime.Handler(c.postContentCreate),
	)))
	// TODO: add middleware
	mux.Handle("/content/edit", methodmux.GetPost(
		hime.Handler(c.contentEdit),
		hime.Handler(c.postContentEdit),
	))

	return mux
}

type ctrl struct {
	Config
}
