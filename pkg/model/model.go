package model

import (
	"context"
	"database/sql"
)

type scanFunc func(...interface{}) error

// DB is the sql.DB context interface
type DB interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}
