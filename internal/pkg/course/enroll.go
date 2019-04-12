package course

import (
	"context"

	"github.com/acoshift/pgsql/pgctx"
)

// IsEnroll checks is user enrolled a course
func IsEnroll(ctx context.Context, userID, courseID string) (bool, error) {
	var b bool
	err := pgctx.QueryRow(ctx, `
		select exists (
			select 1
			from enrolls
			where user_id = $1 and course_id = $2
		)
	`, userID, courseID).Scan(&b)
	return b, err
}
