package model

// Favorite model
type Favorite struct {
	Base
	Stampable
	UserID   string
	CourseID string
}
