package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"cloud.google.com/go/errorreporting"
	"github.com/acoshift/header"
	"github.com/acoshift/middleware"
)

// Alias functions
var (
	Chain = middleware.Chain
	CSRF = middleware.CSRF
)

// Alias types
type (
	CSRFConfig = middleware.CSRFConfig
)

// ErrorLogger logs error and send error page back to response
func ErrorLogger(errClient *errorreporting.Client) middleware.Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					log.Println(err)
					debug.PrintStack()

					if errClient != nil {
						nerr, _ := err.(error)
						errClient.Report(errorreporting.Entry{
							Error: nerr,
							Req:   r,
							Stack: debug.Stack(),
						})
					}
				}
			}()
			h.ServeHTTP(w, r)
		})
	}
}

// SetHeaders sets default headers
func SetHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(header.XContentTypeOptions, "nosniff")
		w.Header().Set(header.XXSSProtection, "1; mode=block")
		w.Header().Set(header.XFrameOptions, "deny")
		// w.Header().Set(header.ContentSecurityPolicy, "img-src https: data:; font-src https: data:; media-src https:;")
		w.Header().Set(header.CacheControl, "no-store")
		w.Header().Set(header.ReferrerPolicy, "same-origin")
		h.ServeHTTP(w, r)
	})
}
