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
