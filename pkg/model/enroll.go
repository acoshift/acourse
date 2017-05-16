package model

import (
	"github.com/acoshift/acourse/pkg/internal"
	"github.com/garyburd/redigo/redis"
)

var (
	enrollStmt, _ = internal.GetDB().Prepare(`
		INSERT INTO enrolls
			(user_id, course_id)
		VALUES
			(?, ?);
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
func IsEnrolled(c redis.Conn, userID, courseID string) (bool, error) {
	_, err := c.Do("ZSCORE", key("c", courseID, "u"), userID)
	if err == redis.ErrNil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
