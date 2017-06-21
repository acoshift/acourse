package servertiming

import (
	"net/http"
	"time"

	"github.com/acoshift/middleware"
)

// Middleware is the server timing middleware
func Middleware() middleware.Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nw := &responseWriter{
				ResponseWriter: w,
				t:              time.Now(),
			}
			h.ServeHTTP(nw, r)
		})
	}
}
