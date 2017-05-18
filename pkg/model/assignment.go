package model

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

// Assignment model
type Assignment struct {
	ID        int64
	Title     string
	Desc      string
	Open      bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Save saves assignment
func (x *Assignment) Save(c redis.Conn) error {
	// if len(x.id) == 0 {
	// 	id, err := redis.Int64(c.Do("INCR", key("id", "a")))
	// 	if err != nil {
	// 		return err
	// 	}
	// 	x.id = strconv.FormatInt(id, 10)
	// }

	// x.UpdatedAt = time.Now()
	// if x.CreatedAt.IsZero() {
	// 	x.CreatedAt = x.UpdatedAt
	// }
	// _, err := c.Do("HSET", key("a"), x.id, enc(x))
	// if err != nil {
	// 	return err
	// }
	return nil
}

// GetAssignments gets assignments
func GetAssignments(c redis.Conn, assignmentIDs []string) ([]*Assignment, error) {
	// xs := make([]*Assignment, len(assignmentIDs))
	// infs := make([]interface{}, len(assignmentIDs)+1)
	// infs[0] = key("a")
	// for i, assignmentID := range assignmentIDs {
	// 	infs[i+1] = assignmentID
	// }
	// bs, err := redis.ByteSlices(c.Do("HMGET", infs...))
	// if err != nil {
	// 	return nil, err
	// }
	// for i, b := range bs {
	// 	var x Assignment
	// 	err = dec(b, &x)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	x.id = assignmentIDs[i]
	// 	xs[i] = &x
	// }
	// return xs, nil
	return nil, nil
}
