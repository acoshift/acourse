package course

import (
	"mime/multipart"
	"time"
)

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

// Enroll enrolls a course
type Enroll struct {
	ID           string
	Price        float64
	PaymentImage *multipart.FileHeader
}
