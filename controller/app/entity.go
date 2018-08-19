package app

import (
	"time"

	"github.com/acoshift/acourse/entity"
)

// PublicCourse type
type PublicCourse struct {
	ID       string
	Option   entity.CourseOption
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
	return x.Type == entity.Live && !x.Start.IsZero()
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
	return x.Type == entity.Live && !x.Start.IsZero()
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
	return x.Type == entity.Live && !x.Start.IsZero()
}
