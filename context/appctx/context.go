package appctx

import (
	"context"
	"net/http"

	"github.com/acoshift/middleware"
	"github.com/acoshift/session"

	"github.com/acoshift/acourse/entity"
)

type (
	userKey struct{}
)

// session id
const sessName = "sess"

// NewUserContext creates new context with user
func NewUserContext(ctx context.Context, user *entity.User) context.Context {
	return context.WithValue(ctx, userKey{}, user)
}

// GetUser gets user from context
func GetUser(ctx context.Context) *entity.User {
	x, _ := ctx.Value(userKey{}).(*entity.User)
	return x
}

// getSession gets session from context
func getSession(ctx context.Context) *session.Session {
	s, err := session.Get(ctx, sessName)
	if err != nil {
		panic(err)
	}
	return s
}

// Repository is appctx middleware storage
type Repository interface {
	GetUser(ctx context.Context, userID string) (*entity.User, error)
}

// Middleware is appctx middleware
func Middleware(repo Repository) middleware.Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			userID := GetUserID(ctx)
			if userID != "" {
				u, err := repo.GetUser(ctx, userID)
				if err == entity.ErrNotFound {
					u = &entity.User{
						ID:       userID,
						Username: userID,
					}
				}
				if err != nil {
					panic(err)
				}
				r = r.WithContext(NewUserContext(ctx, u))
			}
			h.ServeHTTP(w, r)
		})
	}
}
