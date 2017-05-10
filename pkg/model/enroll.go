package model

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

func enroll(c redis.Conn, userID, courseID string) {
	now := time.Now().UnixNano()
	c.Send("ZADD", key("u", userID, "e"), now, courseID)
	c.Send("ZADD", key("c", courseID, "u"), now, userID)
}

// Enroll an user to a course
func Enroll(c redis.Conn, userID, courseID string) error {
	c.Send("MULTI")
	enroll(c, userID, courseID)
	_, err := c.Do("EXEC")
	if err != nil {
		return err
	}
	return nil
}
