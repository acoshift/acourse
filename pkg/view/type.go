package view

import (
	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/flash"
)

// CourseData type
type CourseData struct {
	Page          *Page
	Flash         flash.Flash
	Course        *model.Course
	Enrolled      bool
	Owned         bool
	PendingEnroll bool
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
