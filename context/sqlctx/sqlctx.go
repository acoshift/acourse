package sqlctx

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/acoshift/pgsql"

	"github.com/acoshift/middleware"
)

// Queryer type
type Queryer interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type key int

const (
	keyDB key = iota
	keyQueryer
)

// NewDBContext creates new db context
func NewDBContext(ctx context.Context, db *sql.DB) context.Context {
	ctx = context.WithValue(ctx, keyDB, db)
	ctx = context.WithValue(ctx, keyQueryer, db)
	return ctx
}

// Middleware creates new db middleware
func Middleware(db *sql.DB) middleware.Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = NewDBContext(ctx, db)
			r = r.WithContext(ctx)
			h.ServeHTTP(w, r)
		})
	}
}

// GetQueryer gets queryer from context
func GetQueryer(ctx context.Context) Queryer {
	return ctx.Value(keyQueryer).(Queryer)
}

// RunInTx runs f inside tx
func RunInTx(ctx context.Context, f func(context.Context) error) error {
	db := ctx.Value(keyDB).(*sql.DB)

	return pgsql.RunInTx(db, nil, func(tx *sql.Tx) error {
		ctx := context.WithValue(ctx, keyQueryer, tx)
		return f(ctx)
	})
}
