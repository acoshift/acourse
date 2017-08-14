package model

import "time"

// UserAssignment model
type UserAssignment struct {
	ID          string
	UserID      string
	CourseID    string
	DownloadURL string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
