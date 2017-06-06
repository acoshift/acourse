package appctx

import (
	"context"

	"github.com/acoshift/acourse/pkg/model"
)

type (
	userKey struct{}
	xsrfKey struct{}
)

// WithUser creates new context with user value
func WithUser(ctx context.Context, user *model.User) context.Context {
	return context.WithValue(ctx, userKey{}, user)
}

// GetUser gets user from context
func GetUser(ctx context.Context) *model.User {
	x, _ := ctx.Value(userKey{}).(*model.User)
	return x
}

// WithXSRFToken creates new context with xsrf value
func WithXSRFToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, xsrfKey{}, token)
}

// GetXSRFToken gets xsrf token from context
func GetXSRFToken(ctx context.Context) string {
	x, _ := ctx.Value(xsrfKey{}).(string)
	return x
}
