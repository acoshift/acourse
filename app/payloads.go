package app

import (
	"time"
)

// UserPayload type
type UserPayload struct {
	Username string
	Name     string
	Photo    string
	AboutMe  string
}

// CoursePayload type
type CoursePayload struct {
	Title            string
	ShortDescription string
	Description      string
	Photo            string
	Start            time.Time
	Video            string
	Type             string
	Contents         []*CourseContentPayload
	Enroll           bool
	Public           bool
	Attend           bool
	Assignment       bool
	Purchase         bool
}

// CourseContentPayload type
type CourseContentPayload struct {
	Title       string
	Description string
	Video       string
	DownloadURL string
}

// CourseEnrollPayload type
type CourseEnrollPayload struct {
	Code string
	URL  string
}
