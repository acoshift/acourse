package internal

import (
	"net/http"
)

// Turbolinks sets Turbolinks-Location to request uri
func Turbolinks(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Turbolinks-Location", r.RequestURI)
		h.ServeHTTP(w, r)
	})
}
