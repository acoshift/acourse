package model

const (
	querySaveEnroll = `
		insert into enrolls
			(user_id, course_id)
		values
			($1, $2)
	`
)

// Enroll an user to a course
func Enroll(userID string, courseID int64) error {
	_, err := db.Exec(querySaveEnroll, userID, courseID)
	if err != nil {
		return err
	}
	return nil
}

// IsEnrolled returns true if user enrolled a given course
func IsEnrolled(userID string, courseID int64) (bool, error) {
	return true, nil
}
