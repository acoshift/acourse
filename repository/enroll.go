package repository

import (
	"context"

	"github.com/acoshift/acourse/context/sqlctx"
)

// RegisterEnroll create enroll data for an user to a course
func RegisterEnroll(ctx context.Context, userID string, courseID string) error {
	q := sqlctx.GetQueryer(ctx)

	_, err := q.Exec(`
		insert into enrolls
			(user_id, course_id)
		values
			($1, $2)
	`, userID, courseID)
	return err
}

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
