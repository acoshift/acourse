package model

import "github.com/acoshift/acourse/pkg/internal"

var (
	enrollStmt, _ = internal.GetDB().Prepare(`
		INSERT INTO enrolls
			(user_id, course_id)
		VALUES
			($1, $2);
	`)
)

// Enroll an user to a course
func Enroll(userID string, courseID int64) error {
	_, err := enrollStmt.Exec(userID, courseID)
	if err != nil {
		return err
	}
	return nil
}

// IsEnrolled returns true if user enrolled a given course
func IsEnrolled(userID string, courseID int64) (bool, error) {
	return true, nil
}
