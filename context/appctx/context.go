package appctx

import (
	"context"
	"net/http"

	"github.com/acoshift/session"

	"github.com/acoshift/acourse/appsess"
	"github.com/acoshift/acourse/entity"
	"github.com/acoshift/acourse/repository"
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

// GetSession gets session from context
func GetSession(ctx context.Context) *session.Session {
	s, err := session.Get(ctx, sessName)
	if err != nil {
		panic(err)
	}
	return s
}

// Middleware is appctx middleware
func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		s := GetSession(ctx)
		id := appsess.GetUserID(s)
		if len(id) > 0 {
			u, err := repository.GetUser(ctx, id)
			if err == entity.ErrNotFound {
				u = &entity.User{
					ID:       id,
					Username: id,
				}
			}
			r = r.WithContext(NewUserContext(ctx, u))
		}
		h.ServeHTTP(w, r)
	})
}
