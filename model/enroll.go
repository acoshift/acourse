package model

// Enroll model
type Enroll struct {
	Base
	Stampable
	UserID   string
	CourseID string
}
