package model

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

// UserAssignment model
type UserAssignment struct {
	id          string
	UserID      string
	CourseID    string
	DownloadURL string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ID returns user assignment id
func (x *UserAssignment) ID() string {
	return x.id
}

// GetUserAssignments gets user assignments
func GetUserAssignments(c redis.Conn, userAssignmentIDs []string) ([]*UserAssignment, error) {
	return nil, nil
}
