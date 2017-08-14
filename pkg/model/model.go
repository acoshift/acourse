package model

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
)

type scanFunc func(...interface{}) error

// Errors
var (
	ErrNotFound = errors.New("not found")
)

// DB is the sql.DB context interface
type DB interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

func strID(id int64) string {
	return strconv.FormatInt(id, 10)
}
