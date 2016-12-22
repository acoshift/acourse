package main

import (
	"acourse/app"
	"acourse/store"
)

// UserController type
type UserController struct {
	db *store.DB
}

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
func ToUserMeView(m *store.User) *app.UserMeView {
	return &app.UserMeView{
		ID:       m.ID,
		Name:     m.Name,
		Username: m.Username,
		Photo:    m.Photo,
		AboutMe:  m.AboutMe,
	}
}

// ToRoleView builds a RoleView fromn a Role model
func ToRoleView(m *store.Role) *app.RoleView {
	return &app.RoleView{
		Admin:      m.Admin,
		Instructor: m.Instructor,
	}
}

// Show runs show action
func (c *UserController) Show(ctx *app.UserShowContext) error {
	x, err := c.db.UserFindUsername(ctx.UserID)
	if err != nil {
		return ctx.InternalServerError()
	}
	if x == nil {
		return ctx.NotFound()
	}

	if ctx.CurrentUserID == ctx.UserID {
		// show me view
		res := ToUserMeView(x)
		role, err := c.db.RoleFindByUserID(ctx.UserID)
		if err != nil {
			return ctx.InternalServerError()
		}
		res.Role = ToRoleView(role)

		return ctx.OKMe(res)
	}
	res := ToUserView(x)
	return ctx.OK(res)
}

// Update runs update action
func (c *UserController) Update(ctx *app.UserUpdateContext) error {
	return nil
}
