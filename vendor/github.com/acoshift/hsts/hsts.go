// Package hsts provides net/http middleware for hsts
// see https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Strict-Transport-Security
package hsts

import (
	"net/http"
	"strconv"
	"time"

	"github.com/acoshift/header"
	"github.com/acoshift/middleware"
)

// Config is the hsts config
type Config struct {
	Skipper           middleware.Skipper
	MaxAge            time.Duration
	IncludeSubDomains bool
	Preload           bool
}

// New creates new CORS middleware
func New(config Config) func(http.Handler) http.Handler {
	if config.Skipper == nil {
		config.Skipper = middleware.DefaultSkipper
	}

	hs := "max-age=" + strconv.FormatInt(int64(config.MaxAge/time.Second), 10)
	if config.IncludeSubDomains {
		hs += "; includeSubDomains"
	}
	if config.Preload {
		hs += "; preload"
	}

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if config.Skipper(r) {
				h.ServeHTTP(w, r)
				return
			}
			w.Header().Set(header.StrictTransportSecurity, hs)
			h.ServeHTTP(w, r)
		})
	}
}
