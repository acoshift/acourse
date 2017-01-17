package model

// Assignment model
type Assignment struct {
	Base
	Stampable
	UserID   string
	CourseID string
}
