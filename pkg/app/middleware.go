package app

import (
	"net/http"

	"net/url"

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

func onlyAdmin(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := session.Get(r.Context())
		id, _ := s.Get(keyUserID).(string)
		c := internal.GetPrimaryDB()
		defer c.Close()
		b, err := model.IsUserAdmin(c, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !b {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
		h.ServeHTTP(w, r)
	})
}
