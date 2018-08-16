package repository

import (
	"context"

	"github.com/acoshift/acourse/context/sqlctx"
)

// IsEnrolled returns true if user enrolled a given course
func IsEnrolled(ctx context.Context, userID string, courseID string) (enrolled bool, err error) {
	q := sqlctx.GetQueryer(ctx)

	err = q.QueryRow(`
		select exists (
			select 1
			from enrolls
			where user_id = $1 and course_id = $2
		)
	`, userID, courseID).Scan(&enrolled)
	return
}
