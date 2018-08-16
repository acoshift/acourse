package service

import (
	"mime/multipart"
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

// CreateCourse type
type CreateCourse struct {
	Title     string
	ShortDesc string
	LongDesc  string
	Image     *multipart.FileHeader
	Start     time.Time
}

// UpdateCourse type
type UpdateCourse struct {
	ID        string
	Title     string
	ShortDesc string
	LongDesc  string
	Image     *multipart.FileHeader
	Start     time.Time
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

// UpdateCourseModel type
type UpdateCourseModel struct {
	ID        string
	Title     string
	ShortDesc string
	LongDesc  string
	Start     time.Time
}

// RegisterPayment type
type RegisterPayment struct {
	UserID        string
	CourseID      string
	Image         string
	Price         float64
	OriginalPrice float64
	Code          string
	Status        int
}

// User type
type User struct {
	ID    string
	Name  string
	Email string
}
