package redirecthttps

import (
	"net/http"

	"github.com/acoshift/header"
	"github.com/acoshift/middleware"
)

// Config is redirect https middleware config
type Config struct {
	Skipper middleware.Skipper
	Mode    Mode
}

// Mode is the redirect https mode
type Mode int

const (
	// OnlyConnectionState check only connection state from r.TLS
	OnlyConnectionState Mode = iota

	// OnlyProxy check only X-Forwarded-Proto in request header
	OnlyProxy

	// All check both X-Forwarded-Proto and Request
	All
)

// New creates new redirect https middleware
func New(config Config) middleware.Middleware {
	if config.Skipper == nil {
		config.Skipper = middleware.DefaultSkipper
	}

	var (
		checkProxy, checkRequest func(*http.Request) bool
	)

	if config.Mode == OnlyProxy || config.Mode == All {
		checkProxy = func(r *http.Request) bool {
			return r.Header.Get(header.XForwardedProto) == "http"
		}
	} else {
		checkProxy = func(*http.Request) bool { return false }
	}

	if config.Mode == OnlyConnectionState || config.Mode == All {
		checkRequest = func(r *http.Request) bool {
			return r.TLS != nil
		}
	} else {
		checkRequest = func(*http.Request) bool { return false }
	}

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if config.Skipper(r) {
				h.ServeHTTP(w, r)
				return
			}

			if checkProxy(r) || checkRequest(r) {
				http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
