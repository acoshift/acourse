package appctx

import (
	"context"

	"github.com/acoshift/session"

	"github.com/acoshift/acourse/entity"
)

type (
	userKey        struct{}
	xsrfKey        struct{}
	redisPoolKey   struct{}
	redisPrefixKey struct{}
	cachePoolKey   struct{}
	cachePrefixKey struct{}
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

// NewXSRFTokenContext creates new context with XSRF Token
func NewXSRFTokenContext(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, xsrfKey{}, token)
}

// GetXSRFToken gets xsrf token from context
func GetXSRFToken(ctx context.Context) string {
	x, _ := ctx.Value(xsrfKey{}).(string)
	return x
}

// GetSession gets session from context
func GetSession(ctx context.Context) *session.Session {
	return session.Get(ctx, sessName)
}
