package appctx

import (
	"context"
	"database/sql"

	"github.com/acoshift/session"

	"github.com/acoshift/acourse/entity"
)

type (
	userKey        struct{}
	xsrfKey        struct{}
	courseURLKey   struct{}
	dbKey          struct{}
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

// NewCourseURLContext creates new context with course url
func NewCourseURLContext(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, courseURLKey{}, v)
}

// GetCourseURL gets course url from context
func GetCourseURL(ctx context.Context) string {
	x, _ := ctx.Value(courseURLKey{}).(string)
	return x
}

// GetSession gets session from context
func GetSession(ctx context.Context) *session.Session {
	return session.Get(ctx, sessName)
}

// DB type
type DB interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// Tx type
type Tx interface {
	Rollback() error
	Commit() error
}

// NewDatabaseContext creates new context with database connection
func NewDatabaseContext(ctx context.Context, v *sql.DB) context.Context {
	return context.WithValue(ctx, dbKey{}, v)
}

// GetDatabase gets database connection from context or panic
func GetDatabase(ctx context.Context) DB {
	return ctx.Value(dbKey{}).(DB)
}

// NewTransactionContext creates new context with transaction
func NewTransactionContext(ctx context.Context) (context.Context, Tx, error) {
	db := ctx.Value(dbKey{}).(*sql.DB)
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	return context.WithValue(ctx, dbKey{}, tx), tx, nil
}

// GetTransaction gets database but panic if db is not transaction
func GetTransaction(ctx context.Context) DB {
	x := ctx.Value(dbKey{})
	if _, ok := x.(Tx); !ok {
		panic("database is not transaction")
	}
	return x.(DB)
}
