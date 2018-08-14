package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

// CORSConfig is the cors config
type CORSConfig struct {
	Skipper          Skipper
	AllowAllOrigins  bool
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	ExposeHeaders    []string
	MaxAge           time.Duration
}

// DefaultCORS is the default cors config for public api
var DefaultCORS = CORSConfig{
	AllowAllOrigins: true,
	AllowMethods: []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
	},
	AllowHeaders: []string{
		"Content-Type",
		"Authorization",
	},
	MaxAge: time.Hour,
}

// CORS creates new CORS middleware
func CORS(config CORSConfig) Middleware {
	if config.Skipper == nil {
		config.Skipper = DefaultSkipper
	}

	allowMethods := strings.Join(config.AllowMethods, ",")
	allowHeaders := strings.Join(config.AllowHeaders, ",")
	exposeHeaders := strings.Join(config.ExposeHeaders, ",")

	maxAge := ""
	if config.MaxAge > time.Duration(0) {
		maxAge = strconv.FormatInt(int64(config.MaxAge/time.Second), 10)
	}

	allowOrigins := make(map[string]struct{})
	for _, v := range config.AllowOrigins {
		allowOrigins[v] = struct{}{}
	}

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if config.Skipper(r) {
				h.ServeHTTP(w, r)
				return
			}

			if origin := r.Header.Get("Origin"); origin != "" {
				h := w.Header()

				if config.AllowAllOrigins {
					h.Set("Access-Control-Allow-Origin", "*")
				} else if _, ok := allowOrigins[origin]; ok {
					h.Set("Access-Control-Allow-Origin", origin)
				} else {
					w.WriteHeader(http.StatusForbidden)
					return
				}

				if config.AllowCredentials {
					h.Set("Access-Control-Allow-Credentials", "true")
				}

				if r.Method == http.MethodOptions {
					if allowMethods != "" {
						h.Set("Access-Control-Allow-Methods", allowMethods)
					}
					if allowHeaders != "" {
						h.Set("Access-Control-Allow-Headers", allowHeaders)
					}
					if maxAge != "" {
						h.Set("Access-Control-Max-Age", maxAge)
					}

					if !config.AllowAllOrigins {
						h.Add("Vary", "Origin")
					}

					h.Add("Vary", "Access-Control-Request-Method")
					h.Add("Vary", "Access-Control-Request-Headers")

					w.WriteHeader(http.StatusNoContent)
					return
				}

				if exposeHeaders != "" {
					h.Set("Access-Control-Expose-Headers", exposeHeaders)
				}

				if !config.AllowAllOrigins {
					h.Set("Vary", "Origin")
				}
			}
			h.ServeHTTP(w, r)
		})
	}
}
