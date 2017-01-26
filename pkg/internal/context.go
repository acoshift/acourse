package internal

import "context"

type contextKey int

// Context Keys
const (
	contextKeyUserID contextKey = iota
)

// WithUserID creates new context with user id value
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, contextKeyUserID, userID)
}

// GetUserID extract user id from context
func GetUserID(ctx context.Context) string {
	userID, _ := ctx.Value(contextKeyUserID).(string)
	return userID
}
