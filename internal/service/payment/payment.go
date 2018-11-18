package payment

import (
	"context"

	"github.com/acoshift/acourse/internal/pkg/context/sqlctx"
	"github.com/acoshift/acourse/internal/pkg/dispatcher"
	"github.com/acoshift/acourse/internal/pkg/model/payment"
)

// Init inits payment
func Init() {
	dispatcher.Register(setStatus)
	dispatcher.Register(hasPending)
}

func setStatus(ctx context.Context, m *payment.SetStatus) error {
	_, err := sqlctx.Exec(ctx, `
		update payments
		set
			status = $2,
			updated_at = now(),
			at = now()
		where id = $1
	`, m.ID, m.Status)
	return err
}

func hasPending(ctx context.Context, m *payment.HasPending) error {
	return sqlctx.QueryRow(ctx, `
		select exists (
			select 1
			from payments
			where user_id = $1 and course_id = $2 and status = $3
		)
	`, m.UserID, m.CourseID, payment.Pending).Scan(&m.Result)
}
