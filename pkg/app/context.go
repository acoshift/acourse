package app

import (
	"context"
	"database/sql"

	"github.com/acoshift/session"
	"github.com/garyburd/redigo/redis"
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
func NewUserContext(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userKey{}, user)
}

// GetUser gets user from context
func GetUser(ctx context.Context) *User {
	x, _ := ctx.Value(userKey{}).(*User)
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

// NewRedisPoolContext creates new context with redis pool
func NewRedisPoolContext(ctx context.Context, pool *redis.Pool, prefix string) context.Context {
	ctx = context.WithValue(ctx, redisPoolKey{}, pool)
	ctx = context.WithValue(ctx, redisPrefixKey{}, prefix)
	return ctx
}

// GetRedisPool gets redis pool from context
func GetRedisPool(ctx context.Context) (*redis.Pool, string) {
	return ctx.Value(redisPoolKey{}).(*redis.Pool), ctx.Value(redisPrefixKey{}).(string)
}

// NewCachePoolContext creates new context with cache pool
func NewCachePoolContext(ctx context.Context, pool *redis.Pool, prefix string) context.Context {
	ctx = context.WithValue(ctx, cachePoolKey{}, pool)
	ctx = context.WithValue(ctx, cachePrefixKey{}, prefix)
	return ctx
}

// GetCachePool gets cache pool from context
func GetCachePool(ctx context.Context) (*redis.Pool, string) {
	return ctx.Value(cachePoolKey{}).(*redis.Pool), ctx.Value(cachePrefixKey{}).(string)
}
