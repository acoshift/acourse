package view

import (
	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/flash"
)

// Page type provides layout data like title, description, and og
type Page struct {
	Title string
	Desc  string
	Image string
	URL   string
}

// IndexData type
type IndexData struct {
	*Page
	Courses []*model.Course
}

// AuthData type
type AuthData struct {
	*Page
	flash.Flash
}

// ProfileData type
type ProfileData struct {
	*Page
	flash.Flash
	OwnCourses      []*model.Course
	EnrolledCourses []*model.Course
}

// CourseData type
type CourseData struct {
	*Page
	Course   *model.Course
	Enrolled bool
}

// AdminUsersData type
type AdminUsersData struct {
	*Page
	Users []*model.User
}

// AdminCoursesData type
type AdminCoursesData struct {
	*Page
	Courses []*model.Course
}

// AdminPaymentsData type
type AdminPaymentsData struct {
	*Page
	Payments []*model.Payment
}
