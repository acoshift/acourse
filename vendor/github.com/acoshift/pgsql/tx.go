package pgsql

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

// TxOptions is the transaction options
type TxOptions struct {
	sql.TxOptions
	MaxAttempts int
}

const (
	defaultMaxAttempts = 10
)

// RunInTx runs fn inside retryable transaction
func RunInTx(db *sql.DB, opts *TxOptions, fn func(*sql.Tx) error) error {
	return RunInTxContext(context.Background(), db, opts, fn)
}

// RunInTxContext runs fn inside retryable transaction with context
func RunInTxContext(ctx context.Context, db *sql.DB, opts *TxOptions, fn func(*sql.Tx) error) error {
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
		if err == nil {
			return nil
		}
		pqErr, ok := err.(*pq.Error)
		if retryable := ok && (pqErr.Code == "40001"); !retryable {
			return err
		}
	}

	return nil
}
