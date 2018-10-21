package sqlctx

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/acoshift/middleware"
	"github.com/acoshift/pgsql"
)

// NewContext creates new db context
func NewContext(ctx context.Context, db *sql.DB) context.Context {
	ctx = context.WithValue(ctx, ctxKeyDB{}, db)
	ctx = context.WithValue(ctx, ctxKeyQueryer{}, db)
	return ctx
}

// Middleware creates new db middleware
func Middleware(db *sql.DB) middleware.Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(NewContext(r.Context(), db))
			h.ServeHTTP(w, r)
		})
	}
}

type queryer interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type (
	ctxKeyDB      struct{}
	ctxKeyQueryer struct{}
)

// Exec runs exec context
func Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return q(ctx).ExecContext(ctx, query, args...)
}

// Query runs query context
func Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return q(ctx).QueryContext(ctx, query, args...)
}

// QueryRow runs query row context
func QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return q(ctx).QueryRowContext(ctx, query, args...)
}

func q(ctx context.Context) queryer {
	return ctx.Value(ctxKeyQueryer{}).(queryer)
}

// RunInTx runs f inside tx
func RunInTx(ctx context.Context, f func(context.Context) error) error {
	if _, ok := ctx.Value(ctxKeyQueryer{}).(*sql.Tx); ok {
		return f(ctx)
	}

	db := ctx.Value(ctxKeyDB{}).(*sql.DB)

	return pgsql.RunInTx(db, nil, func(tx *sql.Tx) error {
		return f(context.WithValue(ctx, ctxKeyQueryer{}, tx))
	})
}
