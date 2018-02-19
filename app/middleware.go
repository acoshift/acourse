package app

import (
	"database/sql"
	"log"
	"net/http"
	"net/url"
	"runtime/debug"

	"github.com/acoshift/header"
	"github.com/acoshift/hime"
	"github.com/acoshift/middleware"

	"github.com/acoshift/acourse/appctx"
	"github.com/acoshift/acourse/entity"
	"github.com/acoshift/acourse/repository"
)

func errorRecovery(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				log.Println(err)
				debug.PrintStack()
			}
		}()
		h.ServeHTTP(w, r)
	})
}

func mustSignedIn(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := appctx.GetSession(r.Context())
		id := getUserID(s)
		if len(id) == 0 {
			http.Redirect(w, r, "/signin?r="+url.QueryEscape(r.RequestURI), http.StatusFound)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func mustNotSignedIn(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := appctx.GetSession(r.Context())
		id := getUserID(s)
		if len(id) > 0 {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func fetchUser() middleware.Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			s := appctx.GetSession(ctx)
			id := getUserID(s)
			if len(id) > 0 {
				u, err := repository.GetUser(db, id)
				if err == entity.ErrNotFound {
					u = &entity.User{
						ID:       id,
						Username: id,
					}
				}
				r = r.WithContext(appctx.NewUserContext(ctx, u))
			}
			h.ServeHTTP(w, r)
		})
	}
}

func onlyAdmin(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := appctx.GetUser(r.Context())
		if u == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if !u.Role.Admin.Bool {
			http.Error(w, "Forbidden", http.StatusForbidden)
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
		if !u.Role.Instructor.Bool {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func isCourseOwner(h http.Handler) http.Handler {
	return hime.H(func(ctx hime.Context) hime.Result {
		u := appctx.GetUser(ctx)
		if u == nil {
			return ctx.Redirect("signin")
		}

		id := ctx.FormValue("id")

		var ownerID string
		err := db.QueryRowContext(ctx, `select user_id from courses where id = $1`, id).Scan(&ownerID)
		if err == sql.ErrNoRows {
			return notFound(ctx)
		}
		must(err)

		if ownerID != u.ID {
			return ctx.Redirect("/")
		}
		return ctx.Handle(h)
	})
}

func setHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(header.XContentTypeOptions, "nosniff")
		w.Header().Set(header.XXSSProtection, "1; mode=block")
		w.Header().Set(header.XFrameOptions, "deny")
		w.Header().Set(header.ContentSecurityPolicy, "img-src https: data:; font-src https: data:; media-src https:;")
		h.ServeHTTP(w, r)
	})
}

func cache(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(header.CacheControl, "public, max-age=31536000")
		h.ServeHTTP(w, r)
	})
}
