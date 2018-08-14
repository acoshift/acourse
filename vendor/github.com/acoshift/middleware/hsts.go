package middleware

import (
	"net/http"
	"strconv"
	"time"
)

// HSTSConfig is the HSTS config
type HSTSConfig struct {
	Skipper           Skipper
	MaxAge            time.Duration
	IncludeSubDomains bool
	Preload           bool
}

// Pre-defiend config
var (
	DefaultHSTS = HSTSConfig{
		Skipper:           SkipHTTP,
		MaxAge:            31536000 * time.Second,
		IncludeSubDomains: false,
		Preload:           false,
	}

	PreloadHSTS = HSTSConfig{
		Skipper:           SkipHTTP,
		MaxAge:            63072000 * time.Second,
		IncludeSubDomains: true,
		Preload:           true,
	}
)

// HSTS creates new HSTS middleware
func HSTS(config HSTSConfig) Middleware {
	if config.Skipper == nil {
		config.Skipper = DefaultSkipper
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
			w.Header().Set("Strict-Transport-Security", hs)
			h.ServeHTTP(w, r)
		})
	}
}

// HSTSPreload is the short-hand for HSTS(PreloadHSTS)
func HSTSPreload() Middleware {
	return HSTS(PreloadHSTS)
}
