package appctx

import (
	"context"

	"github.com/acoshift/acourse/pkg/model"
)

type userKey struct{}

// WithUser creates new context with user value
func WithUser(ctx context.Context, user *model.User) context.Context {
	return context.WithValue(ctx, userKey{}, user)
}

// GetUser gets user from context
func GetUser(ctx context.Context) *model.User {
	x, _ := ctx.Value(userKey{}).(*model.User)
	return x
}
