package app

import (
	"net/http"
	"net/url"
	"runtime/debug"

	"github.com/acoshift/acourse/pkg/internal"
	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/session"
)

func recovery(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				var p string
				switch t := r.(type) {
				case string:
					p = t
				case error:
					p = t.Error()
				default:
					p = "unknown"
				}
				debug.PrintStack()
				http.Error(w, p, http.StatusInternalServerError)
			}
		}()
		h.ServeHTTP(w, r)
	})
}

func mustSignedIn(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := session.Get(r.Context())
		id, _ := s.Get(keyUserID).(string)
		if len(id) == 0 {
			http.Redirect(w, r, "/signin?r="+url.QueryEscape(r.RequestURI), http.StatusFound)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func mustNotSignedIn(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := session.Get(r.Context())
		id, _ := s.Get(keyUserID).(string)
		if len(id) > 0 {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func fetchUser(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		s := session.Get(ctx)
		id, _ := s.Get(keyUserID).(string)
		if len(id) > 0 {
			u, err := model.GetUser(id)
			if err == model.ErrNotFound {
				u = &model.User{}
				u.ID = id
			}
			r = r.WithContext(internal.WithUser(ctx, u))
		}
		h.ServeHTTP(w, r)
	})
}

func onlyAdmin(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, _ := internal.GetUser(r.Context()).(*model.User)
		if u == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if !u.Role.Admin {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		h.ServeHTTP(w, r)
	})
}
