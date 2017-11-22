package session

import (
	"context"
	"net/http"

	"github.com/acoshift/middleware"
)

type (
	managerKey struct{}
	requestKey struct{}
	storageKey struct{}
)

// Middleware is the Manager middleware wrapper
//
// New(config).Middleware()
func Middleware(config Config) middleware.Middleware {
	return New(config).Middleware()
}

// Middleware injects session manager into request's context.
//
// All data changed before write response writer's header will be save.
func (m *Manager) Middleware() middleware.Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// inject manager
			ctx = context.WithValue(ctx, managerKey{}, m)

			// inject request
			ctx = context.WithValue(ctx, requestKey{}, r)

			// inject session saver
			storage := make(map[string]*Session)
			ctx = context.WithValue(ctx, storageKey{}, storage)

			nr := r.WithContext(ctx)
			nw := sessionWriter{
				ResponseWriter: w,
				beforeWriteHeader: func() {
					for _, s := range storage {
						m.Save(w, s)
					}
				},
			}
			h.ServeHTTP(&nw, nr)
		})
	}
}

// Get gets session from context
func Get(ctx context.Context, name string) *Session {
	m, _ := ctx.Value(managerKey{}).(*Manager)
	if m == nil {
		// request not pass middleware
		return nil
	}

	// try get session from storage first
	// to preserve session data from difference handler
	storage := ctx.Value(storageKey{}).(map[string]*Session)
	if s, ok := storage[name]; ok {
		return s
	}

	// get session from manager
	s := m.Get(ctx.Value(requestKey{}).(*http.Request), name)

	// save session to storage for later get
	storage[name] = s
	return s
}
