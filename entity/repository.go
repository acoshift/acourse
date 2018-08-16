package entity

import (
	"time"
)

// RegisterUser type
type RegisterUser struct {
	ID       string
	Username string
	Name     string
	Email    string
	Image    string
}

// UpdateUser type
type UpdateUser struct {
	ID       string
	Username string
	Name     string
	AboutMe  string
}

// RegisterCourse type
type RegisterCourse struct {
	UserID    string
	Title     string
	ShortDesc string
	LongDesc  string
	Image     string
	Start     time.Time
}

// UpdateCourse type
type UpdateCourse struct {
	ID        string
	Title     string
	ShortDesc string
	LongDesc  string
	Start     time.Time
}

// RegisterCourseContent type
type RegisterCourseContent struct {
	CourseID  string
	Title     string
	LongDesc  string
	VideoID   string
	VideoType int
}
