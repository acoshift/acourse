package repository

import (
	"database/sql"
)

// Enroll an user to a course
func Enroll(q Queryer, userID string, courseID string) error {
	_, err := q.Exec(`
		INSERT INTO enrolls (user_id, course_id)
		VALUES ($1, $2);
	`, userID, courseID)
	if err != nil {
		return err
	}
	return nil
}

// IsEnrolled returns true if user enrolled a given course
func IsEnrolled(q Queryer, userID string, courseID string) (bool, error) {
	var p int
	err := q.QueryRow(`
		SELECT 1
		  FROM enrolls
		 WHERE user_id = $1
		   AND course_id = $2;
	`, userID, courseID).Scan(&p)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
