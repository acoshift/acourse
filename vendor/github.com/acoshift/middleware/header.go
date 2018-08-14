package middleware

import "net/http"

// AddHeader creates new middleware that adds a header to response
func AddHeader(key, value string) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			addHeaderIfNotExists(w.Header(), key, value)
			h.ServeHTTP(w, r)
		})
	}
}
