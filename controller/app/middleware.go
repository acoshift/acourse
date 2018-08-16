package app

import (
	"net/http"
	"net/url"

	"github.com/acoshift/hime"

	"github.com/acoshift/acourse/appsess"
	"github.com/acoshift/acourse/context/appctx"
	"github.com/acoshift/acourse/entity"
	"github.com/acoshift/acourse/repository"
)

func mustSignedIn(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := appctx.GetSession(r.Context())
		id := appsess.GetUserID(s)
		if len(id) == 0 {
			http.Redirect(w, r, "/signin?r="+url.QueryEscape(r.RequestURI), http.StatusFound)
			return
		}
		h.ServeHTTP(w, r)
	})
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

func isCourseOwner(h http.Handler) http.Handler {
	return hime.Handler(func(ctx *hime.Context) error {
		u := appctx.GetUser(ctx)
		if u == nil {
			return ctx.Redirect("signin")
		}

		id := ctx.FormValue("id")

		ownerID, err := repository.GetCourseUserID(ctx, id)
		if err == entity.ErrNotFound {
			return notFound(ctx)
		}
		if err != nil {
			return err
		}

		if ownerID != u.ID {
			return ctx.Redirect("/")
		}
		return ctx.Handle(h)
	})
}
