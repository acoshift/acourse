package middleware

import (
	"net/http"
	"strings"
)

// NonWWWRedirect redirects www to non-www
func NonWWWRedirect() Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			host := strings.TrimPrefix(r.Host, prefixWWW)
			if len(host) < len(r.Host) {
				http.Redirect(w, r, scheme(r)+"://"+host+r.RequestURI, http.StatusMovedPermanently)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}
