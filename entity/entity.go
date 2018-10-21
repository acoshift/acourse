package entity

import (
	"time"

	"github.com/lib/pq"

	"github.com/acoshift/acourse/model/course"
)

// User type
type User struct {
	ID       string
	Role     UserRole
	Username string
	Name     string
	Email    string
	AboutMe  string
	Image    string
}

// UserRole type
type UserRole struct {
	Admin      bool
	Instructor bool
}

// Course model
type Course struct {
	ID            string
	Option        course.Option
	Owner         *User
	EnrollCount   int64
	Title         string
	ShortDesc     string
	Desc          string
	Image         string
	UserID        string
	Start         pq.NullTime
	URL           string
	Type          int
	Price         float64
	Discount      float64
	Contents      []*course.Content
	EnrollDetail  string
	AssignmentIDs []string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Course type values
const (
	_ = iota
	Live
	Video
	EBook
)

// Video type values
const (
	_ = iota
	Youtube
)

// Link returns id if url is invalid
func (x *Course) Link() string {
	if x.URL == "" {
		return x.ID
	}
	return x.URL
}

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

	User   User
	Course Course
}

// PaymentStatus values
const (
	Pending = iota
	Accepted
	Rejected
	Refunded
)

// Assignment model
type Assignment struct {
	ID    string
	Title string
	Desc  string
	Open  bool
}

// UserAssignment model
type UserAssignment struct {
	ID          string
	UserID      string
	CourseID    string
	DownloadURL string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
