package editor

import (
	"net/http"

	"github.com/moonrhythm/hime"

	"github.com/acoshift/acourse/internal/context/appctx"
	"github.com/acoshift/acourse/internal/controller/share"
	"github.com/acoshift/acourse/internal/entity"
	"github.com/acoshift/acourse/internal/pkg/dispatcher"
	"github.com/acoshift/acourse/internal/pkg/model/course"
)

func (c *ctrl) onlyInstructor(h http.Handler) http.Handler {
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

func (c *ctrl) isCourseOwner(h http.Handler) http.Handler {
	return hime.Handler(func(ctx *hime.Context) error {
		u := appctx.GetUser(ctx)
		if u == nil {
			return ctx.RedirectTo("auth.signin")
		}

		id := ctx.FormValue("id")

		ownerID := course.GetUserID{ID: id}
		err := dispatcher.Dispatch(ctx, &ownerID)
		if err == entity.ErrNotFound {
			return share.NotFound(ctx)
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
