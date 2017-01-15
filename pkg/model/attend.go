package model

// Attend model
type Attend struct {
	Base
	Stampable
	UserID   string
	CourseID string
}

// Attends type
type Attends []*Attend
