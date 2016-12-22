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

// ToRoleView builds a RoleView fromn a Role model
func ToRoleView(m *store.Role) *app.RoleView {
	return &app.RoleView{
		Admin:      m.Admin,
		Instructor: m.Instructor,
	}
}
