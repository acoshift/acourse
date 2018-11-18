package entity

import (
	"time"

	"github.com/lib/pq"

	"github.com/acoshift/acourse/internal/pkg/model/course"
	"github.com/acoshift/acourse/internal/pkg/model/user"
)

// Payment model
type Payment struct {
	ID            string
	UserID        string
	CourseID      string
	Image         string
	Price         float64
	OriginalPrice float64
	Code          string
	Status        int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	At            pq.NullTime

	User   user.User
	Course course.Course
}

// Assignment model
type Assignment struct {
	ID    string
	Title string
	Desc  string
	Open  bool
}
