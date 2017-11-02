package app

import (
	"context"
	"database/sql"
)

type (
	userKey      struct{}
	xsrfKey      struct{}
	courseURLKey struct{}
	dbKey        struct{}
)

// WithUser creates new context with user value
func WithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userKey{}, user)
}

// GetUser gets user from context
func GetUser(ctx context.Context) *User {
	x, _ := ctx.Value(userKey{}).(*User)
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

// WithCourseURL creates new context with course url value
func WithCourseURL(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, courseURLKey{}, v)
}

// GetCourseURL gets course url from context
func GetCourseURL(ctx context.Context) string {
	x, _ := ctx.Value(courseURLKey{}).(string)
	return x
}

// WithDatabase creates new context with database connection
func WithDatabase(ctx context.Context, v *sql.DB) context.Context {
	return context.WithValue(ctx, dbKey{}, v)
}

// GetDatabase gets database connection from context or panic
func GetDatabase(ctx context.Context) *sql.DB {
	return ctx.Value(dbKey{}).(*sql.DB)
}
