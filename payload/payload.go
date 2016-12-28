package payload

import (
	"time"
)

// User type
type User struct {
	Username string
	Name     string
	Photo    string
	AboutMe  string
}

// Course type
type Course struct {
	Title            string
	ShortDescription string
	Description      string
	Photo            string
	Start            time.Time
	Video            string
	Contents         []*CourseContent
	Attend           bool
	Assignment       bool
}

// CourseContent type
type CourseContent struct {
	Title       string
	Description string
	Video       string
	DownloadURL string
}

// CourseEnroll type
type CourseEnroll struct {
	Code string
	URL  string
}
