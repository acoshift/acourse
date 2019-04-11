package app

import (
	"time"

	"github.com/acoshift/acourse/internal/pkg/course"
)

// PublicCourse type
type PublicCourse struct {
	ID       string
	Option   course.Option
	Title    string
	Desc     string
	Image    string
	Start    time.Time
	URL      string
	Type     int
	Price    float64
	Discount float64
}

// Link returns course link
func (x *PublicCourse) Link() string {
	if x.URL != "" {
		return x.URL
	}
	return x.ID
}

// ShowStart returns true if course should show start date
func (x *PublicCourse) ShowStart() bool {
	return x.Type == course.Live && !x.Start.IsZero()
}

// EnrolledCourse type
type EnrolledCourse struct {
	ID    string
	Title string
	Desc  string
	Image string
	Start time.Time
	URL   string
	Type  int
}

// Link returns course link
func (x *EnrolledCourse) Link() string {
	if x.URL != "" {
		return x.URL
	}
	return x.ID
}

// ShowStart returns true if course should show start date
func (x *EnrolledCourse) ShowStart() bool {
	return x.Type == course.Live && !x.Start.IsZero()
}

// OwnCourse type
type OwnCourse struct {
	ID          string
	Title       string
	Desc        string
	Image       string
	Start       time.Time
	URL         string
	Type        int
	EnrollCount int
}

// Link returns course link
func (x *OwnCourse) Link() string {
	if x.URL != "" {
		return x.URL
	}
	return x.ID
}

// ShowStart returns true if course should show start date
func (x *OwnCourse) ShowStart() bool {
	return x.Type == course.Live && !x.Start.IsZero()
}

// Course type
type Course struct {
	ID           string
	Option       course.Option
	Owner        CourseOwner
	Title        string
	ShortDesc    string
	Desc         string
	Image        string
	Start        time.Time
	URL          string
	Type         int
	Price        float64
	Discount     float64
	EnrollDetail string
}

// Link returns course link
func (x *Course) Link() string {
	if x.URL != "" {
		return x.URL
	}
	return x.ID
}

// CourseOwner type
type CourseOwner struct {
	ID    string
	Name  string
	Image string
}
