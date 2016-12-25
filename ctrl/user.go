package ctrl

import (
	"acourse/app"
	"acourse/store"
)

// UserController implements UserController interface
type UserController struct {
	db *store.DB
}

// NewUserController creates new controller
func NewUserController(db *store.DB) *UserController {
	return &UserController{db}
}

// Show runs show action
func (c *UserController) Show(ctx *app.UserShowContext) error {
	x, err := c.db.UserGet(ctx.UserID)
	if err != nil {
		return ctx.InternalServerError(err)
	}
	if x == nil {
		return ctx.NotFound()
	}

	if ctx.CurrentUserID == ctx.UserID {
		// show me view
		role, err := c.db.RoleFindByUserID(ctx.UserID)
		if err != nil {
			return ctx.InternalServerError(err)
		}

		return ctx.OKMe(ToUserMeView(x, ToRoleView(role)))
	}

	return ctx.OK(ToUserView(x))
}

// Update runs update action
func (c *UserController) Update(ctx *app.UserUpdateContext) error {
	role, err := c.db.RoleFindByUserID(ctx.CurrentUserID)
	if err != nil {
		return err
	}

	if !role.Admin && ctx.CurrentUserID != ctx.UserID {
		return ctx.Forbidden()
	}

	user, err := c.db.UserGet(ctx.UserID)
	if err != nil {
		return err
	}
	if user == nil {
		return ctx.NotFound()
	}

	user.Name = ctx.Payload.Name
	user.Username = ctx.Payload.Username
	user.Photo = ctx.Payload.Photo
	user.AboutMe = ctx.Payload.AboutMe

	err = c.db.UserSave(user)
	if err != nil {
		return err
	}

	return ctx.NoContent()
}
