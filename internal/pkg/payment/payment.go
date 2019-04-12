package payment

import (
	"context"

	"github.com/acoshift/pgsql/pgctx"
)

// Status values
const (
	Pending = iota
	Accepted
	Rejected
	Refunded
)

// SetStatus sets payment status
func SetStatus(ctx context.Context, id string, status int) error {
	// language=SQL
	_, err := pgctx.Exec(ctx, `
		update payments
		set
			status = $2,
			updated_at = now(),
			at = now()
		where id = $1
	`, id, status)
	return err
}

// HasPending checks is has course pending payment
func HasPending(ctx context.Context, userID, courseID string) (exists bool, err error) {
	// language=SQL
	err = pgctx.QueryRow(ctx, `
		select exists (
			select 1
			from payments
			where user_id = $1 and course_id = $2 and status = $3
		)
	`, userID, courseID, Pending).Scan(&exists)
	return
}
