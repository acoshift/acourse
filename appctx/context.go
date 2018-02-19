package appctx

import (
	"context"

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

// GetSession gets session from context
func GetSession(ctx context.Context) *session.Session {
	return session.Get(ctx, sessName)
}
