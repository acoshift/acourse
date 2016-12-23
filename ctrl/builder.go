package ctrl

import (
	"acourse/app"
	"acourse/store"
)

// ToUserView builds a UserView from a User model
func ToUserView(m *store.User) *app.UserView {
	return &app.UserView{
		ID:       m.ID,
		Name:     m.Name,
		Username: m.Username,
		Photo:    m.Photo,
		AboutMe:  m.AboutMe,
	}
}

// ToUserMeView builds a UserMeView from a User model
func ToUserMeView(m *store.User, role *app.RoleView) *app.UserMeView {
	return &app.UserMeView{
		ID:       m.ID,
		Name:     m.Name,
		Username: m.Username,
		Photo:    m.Photo,
		AboutMe:  m.AboutMe,
		Role:     role,
	}
}

// ToUserTinyView builds a UserTinyView from a User model
func ToUserTinyView(m *store.User) *app.UserTinyView {
	return &app.UserTinyView{
		ID:       m.ID,
		Name:     m.Name,
		Username: m.Username,
		Photo:    m.Photo,
	}
}

// ToRoleView builds a RoleView fromn a Role model
func ToRoleView(m *store.Role) *app.RoleView {
	return &app.RoleView{
		Admin:      m.Admin,
		Instructor: m.Instructor,
	}
}

// ToCourseView builds a CourseView from a Course model
func ToCourseView(m *store.Course, owner *app.UserTinyView, student int, enroll bool) *app.CourseView {
	return &app.CourseView{}
}

// ToCoursePublicView builds a CourseView from a Course model
func ToCoursePublicView(m *store.Course, owner *app.UserTinyView, student int) *app.CoursePublicView {
	return &app.CoursePublicView{}
}

// ToCourseTinyView builds a CourseTinyView from a Course model
func ToCourseTinyView(m *store.Course, owner *app.UserTinyView, student int) *app.CourseTinyView {
	return &app.CourseTinyView{
		ID:               m.ID,
		Owner:            owner,
		Title:            m.Title,
		ShortDescription: m.ShortDescription,
		Photo:            m.Photo,
		Start:            m.Start,
		URL:              m.URL,
		Type:             string(m.Type),
		Price:            m.Price,
		DiscountedPrice:  m.DiscountedPrice,
		Student:          student,
	}
}
