package course

import (
	"context"

	"github.com/acoshift/pgsql/pgctx"
)

// InsertEnroll inserts enroll
func InsertEnroll(ctx context.Context, courseID, userID string) error {
	// language=SQL
	_, err := pgctx.Exec(ctx, `
		insert into enrolls
			(user_id, course_id)
		values
			($1, $2)
	`, userID, courseID)
	return err
}

// IsEnroll checks is user enrolled a course
func IsEnroll(ctx context.Context, userID, courseID string) (bool, error) {
	var b bool

	// language=SQL
	err := pgctx.QueryRow(ctx, `
		select exists (
			select 1
			from enrolls
			where user_id = $1 and course_id = $2
		)
	`, userID, courseID).Scan(&b)
	return b, err
}
