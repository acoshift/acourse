package repository

import (
	"database/sql"
)

type scanFunc func(...interface{}) error

// Queryer is the sql DB or Tx
type Queryer interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}
