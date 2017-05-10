package internal

import "context"

// contextKey is the internal acourse's context key type
// use for store value in context
type contextKey int

// context key values
const (
	_       contextKey = iota
	keyUser            // user object
)

// WithUser creates new context with user value
func WithUser(ctx context.Context, user interface{}) context.Context {
	return context.WithValue(ctx, keyUser, user)
}

// GetUser gets user from context
func GetUser(ctx context.Context) interface{} {
	return ctx.Value(keyUser)
}
