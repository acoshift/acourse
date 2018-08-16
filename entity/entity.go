package entity

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
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
	Option        CourseOption
	Owner         *User
	EnrollCount   int64
	Title         string
	ShortDesc     string
	Desc          string
	Image         string
	UserID        string
	Start         pq.NullTime
	URL           sql.NullString // MUST not parsable to int
	Type          int
	Price         float64
	Discount      float64
	Contents      []*CourseContent
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

// CourseContent type
type CourseContent struct {
	ID          string
	CourseID    string
	Title       string
	Desc        string
	VideoID     string
	VideoType   int
	DownloadURL string
}

// CourseOption type
type CourseOption struct {
	Public     bool
	Enroll     bool
	Attend     bool
	Assignment bool
	Discount   bool
}

// Link returns id if url is invalid
func (x *Course) Link() string {
	if !x.URL.Valid || len(x.URL.String) == 0 {
		return x.ID
	}
	return x.URL.String
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

// RegisterCourseContent type
type RegisterCourseContent struct {
	CourseID  string
	Title     string
	LongDesc  string
	VideoID   string
	VideoType int
}
