package course

import (
	"mime/multipart"
	"time"

	"github.com/lib/pq"

	"github.com/acoshift/acourse/internal/pkg/model/user"
)

// Course model
type Course struct {
	ID            string
	Option        Option
	Owner         *user.User
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
	Contents      []*Content
	EnrollDetail  string
	AssignmentIDs []string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Link returns id if url is invalid
func (x *Course) Link() string {
	if x.URL == "" {
		return x.ID
	}
	return x.URL
}

// Option type
type Option struct {
	Public     bool
	Enroll     bool
	Attend     bool
	Assignment bool
	Discount   bool
}

// Create creates new course
type Create struct {
	UserID    string
	Title     string
	ShortDesc string
	LongDesc  string
	Image     *multipart.FileHeader
	Start     time.Time

	Result string // course id
}

// Update updates course
type Update struct {
	ID        string
	Title     string
	ShortDesc string
	LongDesc  string
	Image     *multipart.FileHeader
	Start     time.Time
}

// SetOption sets course option
type SetOption struct {
	ID     string
	Option Option
}

// SetImage sets course image
type SetImage struct {
	ID    string
	Image string
}

// GetURL gets course url
type GetURL struct {
	ID string

	Result string
}

// GetUserID gets course user id
type GetUserID struct {
	ID string

	Result string
}

// Get gets course from id
type Get struct {
	ID string

	Result *Course
}

// InsertEnroll inserts enroll
type InsertEnroll struct {
	ID     string
	UserID string
}
