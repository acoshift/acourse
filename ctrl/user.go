package ctrl

import (
	"acourse/app"
	"acourse/store"
)

// UserController type
type UserController struct {
	db *store.DB
}

// NewUserController creates new controller
func NewUserController(db *store.DB) *UserController {
	return &UserController{db}
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
		role, err := c.db.RoleFindByUserID(ctx.UserID)
		if err != nil {
			return ctx.InternalServerError()
		}

		return ctx.OKMe(ToUserMeView(x, ToRoleView(role)))
	}

	return ctx.OK(ToUserView(x))
}

// Update runs update action
func (c *UserController) Update(ctx *app.UserUpdateContext) error {
	return ctx.NoContent()
}
