package editor

import (
	"net/http"

	"github.com/acoshift/methodmux"
	"github.com/moonrhythm/hime"
	"github.com/moonrhythm/httpmux"

	"github.com/acoshift/acourse/internal/app/view"
	"github.com/acoshift/acourse/internal/pkg/bus"
	"github.com/acoshift/acourse/internal/pkg/context/appctx"
	"github.com/acoshift/acourse/internal/pkg/model"
	"github.com/acoshift/acourse/internal/pkg/model/course"
)

// Mount mounts editor handlers
func Mount(m *httpmux.Mux) {
	instructorMux := m.Group("/editor", onlyInstructor)
	instructorMux.Handle("/course/create", methodmux.GetPost(
		hime.Handler(getCourseCreate),
		hime.Handler(postCourseCreate),
	))

	courseOwnerMux := m.Group("/editor", onlyCourseOwner)
	courseOwnerMux.Handle("/course/edit", methodmux.GetPost(
		hime.Handler(getCourseEdit),
		hime.Handler(postCourseEdit),
	))
	courseOwnerMux.Handle("/content", methodmux.GetPost(
		hime.Handler(getContentList),
		hime.Handler(postContentList),
	))
	courseOwnerMux.Handle("/content/create", methodmux.GetPost(
		hime.Handler(getContentCreate),
		hime.Handler(postContentCreate),
	))
	// TODO: add middleware for course content owner
	m.Handle("/editor/content/edit", methodmux.GetPost(
		hime.Handler(getContentEdit),
		hime.Handler(postContentEdit),
	))
}

func onlyInstructor(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := appctx.GetUser(r.Context())
		if u == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if !u.Role.Instructor {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func onlyCourseOwner(h http.Handler) http.Handler {
	return hime.Handler(func(ctx *hime.Context) error {
		u := appctx.GetUser(ctx)
		if u == nil {
			return ctx.RedirectTo("auth.signin")
		}

		id := ctx.FormValue("id")

		ownerID := course.GetUserID{ID: id}
		err := bus.Dispatch(ctx, &ownerID)
		if err == model.ErrNotFound {
			return view.NotFound(ctx)
		}
		if err != nil {
			return err
		}

		if ownerID.Result != u.ID {
			return ctx.Redirect("/")
		}
		return ctx.Handle(h)
	})
}
