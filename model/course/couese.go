package course

import (
	"mime/multipart"
	"time"
)

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

// Enroll enrolls a course
type Enroll struct {
	ID           string
	Price        float64
	PaymentImage *multipart.FileHeader
}
