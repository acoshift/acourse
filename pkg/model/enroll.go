package model

import "database/sql"

// Enroll an user to a course
func Enroll(tx *sql.Tx, userID string, courseID int64) error {
	_, err := tx.Exec(`
		insert into enrolls
			(user_id, course_id)
		values
			($1, $2)
	`, userID, courseID)
	if err != nil {
		return err
	}
	return nil
}

// IsEnrolled returns true if user enrolled a given course
func IsEnrolled(userID string, courseID int64) (bool, error) {
	var p int
	err := db.QueryRow(`select 1 from enrolls where user_id = $1 and course_id = $2`, userID, courseID).Scan(&p)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
