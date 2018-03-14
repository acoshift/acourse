package prefixhandler

import (
	"context"
	"net/http"
	"strings"
)

// New creates new prefix handler
func New(prefix string, ctxKey interface{}, h http.Handler) http.Handler {
	if prefix == "" {
		prefix = "/"
	}

	return http.StripPrefix(prefix, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ps := strings.SplitN(strings.TrimPrefix(r.URL.Path, "/"), "/", 2)

		// WithContext already clone r and r.URL
		r = r.WithContext(context.WithValue(r.Context(), ctxKey, ps[0]))
		r.URL.Path = "/"
		if len(ps) > 1 {
			r.URL.Path += ps[1]
		}

		h.ServeHTTP(w, r)
	}))
}

// Get gets prefix from context key
func Get(ctx context.Context, ctxKey interface{}) string {
	p, _ := ctx.Value(ctxKey).(string)
	return p
}
