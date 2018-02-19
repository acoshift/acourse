package csrf

import (
	"net/http"
	"strings"

	"github.com/acoshift/header"
	"github.com/acoshift/middleware"
)

// Config is the csrf config
type Config struct {
	Origins          []string
	ForbiddenHandler http.Handler
}

// New creates new csrf middleware
func New(config Config) middleware.Middleware {
	if config.ForbiddenHandler == nil {
		config.ForbiddenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Forbidden", http.StatusForbidden)
		})
	}

	checkOrigin := func(r *http.Request) bool {
		origin := r.Header.Get(header.Origin)
		if origin != "" {
			for _, allow := range config.Origins {
				if origin == allow {
					return true
				}
			}
		}

		return false
	}

	checkReferer := func(r *http.Request) bool {
		referer := r.Referer()
		if referer != "" {
			for _, allow := range config.Origins {
				if strings.HasPrefix(referer, allow+"/") {
					return true
				}
			}
		}
		return false
	}

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				if !checkOrigin(r) && !checkReferer(r) {
					config.ForbiddenHandler.ServeHTTP(w, r)
					return
				}
			}
			h.ServeHTTP(w, r)
		})
	}
}
