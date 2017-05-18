package internal

import "context"

type userKey struct{}

// WithUser creates new context with user value
func WithUser(ctx context.Context, user interface{}) context.Context {
	return context.WithValue(ctx, userKey{}, user)
}

// GetUser gets user from context
func GetUser(ctx context.Context) interface{} {
	return ctx.Value(userKey{})
}
