package middleware

import (
	"net/http"
	"strings"
)

// NonWWWRedirect redirects www to non-www
func NonWWWRedirect() Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			host := strings.TrimPrefix(r.Host, "www.")
			if len(host) < len(r.Host) {
				http.Redirect(w, r, scheme(r)+"://"+host+r.RequestURI, http.StatusMovedPermanently)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}

// WWWRedirect redirects non-www to www
func WWWRedirect() Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !strings.HasPrefix(r.Host, "www.") {
				http.Redirect(w, r, scheme(r)+"://www."+r.Host+r.RequestURI, http.StatusMovedPermanently)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}
