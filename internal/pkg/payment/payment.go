package payment

import (
	"context"

	"github.com/acoshift/acourse/internal/pkg/context/sqlctx"
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
	_, err := sqlctx.Exec(ctx, `
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
func HasPending(ctx context.Context, userID, courseID string) (bool, error) {
	var b bool
	err := sqlctx.QueryRow(ctx, `
		select exists (
			select 1
			from payments
			where user_id = $1 and course_id = $2 and status = $3
		)
	`, userID, courseID, Pending).Scan(&b)
	return b, err
}
