package ctrl

import (
	"github.com/acoshift/acourse/pkg/app"
	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/acourse/pkg/store"
	"github.com/acoshift/e"
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
func (c *UserController) Show(ctx *app.UserShowContext) (interface{}, error) {
	if ctx.CurrentUserID == ctx.UserID {
		// show me view
		role, err := c.db.RoleGet(ctx.UserID)
		if err != nil {
			return nil, err
		}
		x, err := c.db.UserMustGet(ctx.UserID)

		return ToUserMeView(x, ToRoleView(role)), nil
	}

	x, err := c.db.UserGet(ctx.UserID)
	if err != nil {
		return nil, err
	}
	if x == nil {
		return nil, e.NotFound
	}

	return ToUserView(x), nil
}

// Update runs update action
func (c *UserController) Update(ctx *app.UserUpdateContext) error {
	role, err := c.db.RoleGet(ctx.CurrentUserID)
	if err != nil {
		return err
	}

	if !role.Admin && ctx.CurrentUserID != ctx.UserID {
		return e.Forbidden
	}

	var user *model.User

	if ctx.CurrentUserID == ctx.UserID {
		user, err = c.db.UserMustGet(ctx.UserID)
	} else {
		user, err = c.db.UserGet(ctx.UserID)
	}
	if err != nil {
		return err
	}
	if user == nil {
		return e.NotFound
	}

	user.Name = ctx.Payload.Name
	user.Username = ctx.Payload.Username
	user.Photo = ctx.Payload.Photo
	user.AboutMe = ctx.Payload.AboutMe

	err = c.db.UserSave(user)
	if err != nil {
		return err
	}

	return nil
}
