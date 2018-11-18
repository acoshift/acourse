package admin

import (
	"time"

	"github.com/acoshift/acourse/internal/model/course"
)

// UserItem type
type UserItem struct {
	ID        string
	Username  string
	Name      string
	Email     string
	Image     string
	CreatedAt time.Time
}

// CourseItem type
type CourseItem struct {
	ID        string
	Title     string
	Image     string
	Type      int
	Price     float64
	Discount  float64
	URL       string
	CreatedAt time.Time
	UpdatedAt time.Time
	Option    course.Option
	Owner     struct {
		ID       string
		Username string
		Image    string
	}
}

// Payment type
type Payment struct {
	ID            string
	Image         string
	Price         float64
	OriginalPrice float64
	Code          string
	Status        int
	CreatedAt     time.Time
	At            time.Time
	User          struct {
		ID       string
		Username string
		Name     string
		Email    string
		Image    string
	}
	Course struct {
		ID    string
		Title string
		Image string
		URL   string
	}
}

// CourseLink returns course link
func (x *Payment) CourseLink() string {
	if x.Course.URL == "" {
		return x.Course.ID
	}
	return x.Course.URL
}
