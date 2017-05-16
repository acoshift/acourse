package model

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

// UserAssignment model
type UserAssignment struct {
	ID          int64
	UserID      string
	CourseID    string
	DownloadURL string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// GetUserAssignments gets user assignments
func GetUserAssignments(c redis.Conn, userAssignmentIDs []string) ([]*UserAssignment, error) {
	return nil, nil
}
