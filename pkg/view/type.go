package view

import (
	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/flash"
)

// Page type provides layout data like title, description, and og
type Page struct {
	Title        string
	Desc         string
	Image        string
	URL          string
	NavbarActive string
}

// IndexData type
type IndexData struct {
	Page    *Page
	Courses []*model.Course
}

// AuthData type
type AuthData struct {
	Page  *Page
	Flash flash.Flash
}

// ProfileData type
type ProfileData struct {
	Page            *Page
	Flash           flash.Flash
	OwnCourses      []*model.Course
	EnrolledCourses []*model.Course
}

// ProfileEditData type
type ProfileEditData struct {
	Page  *Page
	Flash flash.Flash
}

// CourseData type
type CourseData struct {
	Page          *Page
	Flash         flash.Flash
	Course        *model.Course
	Enrolled      bool
	Owned         bool
	PendingEnroll bool
}

// CourseCreateData type
type CourseCreateData struct {
	Page  *Page
	Flash flash.Flash
}

// CourseEditData type
type CourseEditData struct {
	Page   *Page
	Flash  flash.Flash
	Course *model.Course
}

// EditorContentCreateData type
type EditorContentCreateData struct {
	Page        *Page
	Flash       flash.Flash
	CourseTitle string
}

// AdminUsersData type
type AdminUsersData struct {
	Page        *Page
	Users       []*model.User
	CurrentPage int
	TotalPage   int
}

// AdminCoursesData type
type AdminCoursesData struct {
	Page    *Page
	Courses []*model.Course
}

// AdminPaymentsData type
type AdminPaymentsData struct {
	Page        *Page
	Payments    []*model.Payment
	CurrentPage int
	TotalPage   int
}
