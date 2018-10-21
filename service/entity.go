package service

import (
	"time"

	"github.com/lib/pq"
)

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

// User type
type User struct {
	ID    string
	Name  string
	Email string
}

// Payment type
type Payment struct {
	ID            string
	Image         string
	Price         float64
	OriginalPrice float64
	Code          string
	Status        int
	CreatedAt     time.Time
	At            pq.NullTime

	User struct {
		ID       string
		Username string
		Name     string
		Email    string
		Image    string
	}
	Course struct {
		ID    string
		Title string
		Image string
		URL   string
	}
}
