package pgsql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

// ErrAbortTx rollbacks transaction and return nil error
var ErrAbortTx = errors.New("pgsql: abort tx")

// BeginTxer type
type BeginTxer interface {
	BeginTx(context.Context, *sql.TxOptions) (*sql.Tx, error)
}

// TxOptions is the transaction options
type TxOptions struct {
	sql.TxOptions
	MaxAttempts int
}

const (
	defaultMaxAttempts = 10
)

// RunInTx runs fn inside retryable transaction.
//
// see RunInTxContext for more info.
func RunInTx(db BeginTxer, opts *TxOptions, fn func(*sql.Tx) error) error {
	return RunInTxContext(context.Background(), db, opts, fn)
}

// RunInTxContext runs fn inside retryable transaction with context.
// It use Serializable isolation level if tx options isolation is setted to sql.LevelDefault.
//
// RunInTxContext DO NOT handle panic.
// But when panic, it will rollback the transaction.
func RunInTxContext(ctx context.Context, db BeginTxer, opts *TxOptions, fn func(*sql.Tx) error) error {
	if opts == nil {
		opts = &TxOptions{}
	}
	// override invalid max attempts
	if opts.MaxAttempts <= 0 {
		opts.MaxAttempts = defaultMaxAttempts
	}
	// override default isolation level to serializable
	if opts.Isolation == sql.LevelDefault {
		opts.Isolation = sql.LevelSerializable
	}

	f := func() error {
		tx, err := db.BeginTx(ctx, &opts.TxOptions)
		if err != nil {
			return err
		}
		// use defer to also rollback when panic
		defer tx.Rollback()

		err = fn(tx)
		if err != nil {
			return err
		}
		return tx.Commit()
	}

	for i := 0; i < opts.MaxAttempts; i++ {
		err := f()
		if err == nil || err == ErrAbortTx {
			return nil
		}
		pqErr, ok := err.(*pq.Error)
		if retryable := ok && (pqErr.Code == "40001"); !retryable {
			return err
		}
	}

	return nil
}
