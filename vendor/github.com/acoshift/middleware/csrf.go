package middleware

import (
	"net/http"
	"strings"
)

// CSRFConfig is the csrf config
type CSRFConfig struct {
	Origins          []string
	ForbiddenHandler http.Handler
	IgnoreProto      bool
	Force            bool
}

// CSRF creates new csrf middleware
func CSRF(c CSRFConfig) func(http.Handler) http.Handler {
	if c.ForbiddenHandler == nil {
		c.ForbiddenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Forbidden", http.StatusForbidden)
		})
	}

	normalize := func(s string) (string, bool) {
		return s, true
	}

	origins := make([]string, len(c.Origins))
	copy(origins, c.Origins)
	if c.IgnoreProto {
		for i := range origins {
			origins[i], _ = removeProto(origins[i])
		}

		normalize = removeProto
	}

	checkOrigin := func(r *http.Request) bool {
		origin := r.Header.Get("Origin")

		if c.Force || origin != "" {
			origin, b := normalize(origin)
			if !b {
				return false
			}

			for _, allow := range origins {
				if origin == allow {
					return true
				}
			}
			return false
		}

		return true
	}

	checkReferer := func(r *http.Request) bool {
		referer := r.Referer()

		if c.Force || referer != "" {
			referer, b := normalize(referer)
			if !b {
				return false
			}

			for _, allow := range origins {
				if strings.HasPrefix(referer, allow+"/") {
					return true
				}
			}
			return false
		}

		return true
	}

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				if !checkOrigin(r) || !checkReferer(r) {
					c.ForbiddenHandler.ServeHTTP(w, r)
					return
				}
			}
			h.ServeHTTP(w, r)
		})
	}
}

func removeProto(s string) (string, bool) {
	prefix := strings.Index(s, "://")
	if prefix >= 0 {
		return s[prefix+3:], true
	}
	return s, false
}
