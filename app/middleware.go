package app

import (
	"database/sql"
	"log"
	"net/http"
	"net/url"
	"runtime/debug"

	"github.com/acoshift/header"
	"github.com/acoshift/middleware"
	"golang.org/x/net/xsrftoken"

	"github.com/acoshift/acourse/appctx"
	"github.com/acoshift/acourse/entity"
	"github.com/acoshift/acourse/repository"
	"github.com/acoshift/acourse/view"
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

func csrf(baseURL, xsrfSecret string) middleware.Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var id string
			if u := appctx.GetUser(r.Context()); u != nil {
				id = u.ID
			}
			if r.Method == http.MethodPost {
				origin := r.Header.Get(header.Origin)
				if len(origin) > 0 {
					if origin != baseURL {
						http.Error(w, "Not allow cross-site post", http.StatusBadRequest)
						return
					}
				}

				x := r.FormValue("X")
				if !xsrftoken.Valid(x, xsrfSecret, id, r.URL.Path) {
					http.Error(w, "invalid xsrf token, go back, refresh and try again...", http.StatusBadRequest)
					return
				}
				h.ServeHTTP(w, r)
				return
			}
			token := xsrftoken.Generate(xsrfSecret, id, r.URL.Path)
			ctx := appctx.NewXSRFTokenContext(r.Context(), token)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
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
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		u := appctx.GetUser(ctx)
		if u == nil {
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}

		id := r.FormValue("id")

		var ownerID string
		err := db.QueryRowContext(ctx, `select user_id from courses where id = $1`, id).Scan(&ownerID)
		if err == sql.ErrNoRows {
			view.NotFound(w, r)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if ownerID != u.ID {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		h.ServeHTTP(w, r)
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
