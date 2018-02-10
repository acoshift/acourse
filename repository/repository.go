package repository

import (
	"context"
	"database/sql"
)

type scanFunc func(...interface{}) error

// Queryer is the sql DB or Tx
type Queryer interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}
