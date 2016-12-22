package app

import "github.com/labstack/echo"

// UserController is the controller interface for the User actions
type UserController interface {
	Show(*UserShowContext) error
	Update(*UserUpdateContext) error
}

// MountUserController mounts a User resource controller on the given service
func MountUserController(service *echo.Echo, ctrl UserController) {
	h := func(ctx echo.Context) error {
		rctx, err := NewUserShowContext(ctx)
		if err != nil {
			return err
		}
		return ctrl.Show(rctx)
	}
	service.GET("/api/user/:userID", h)
	service.Logger.Info("Mount ctrl User action Show")

	h = func(ctx echo.Context) error {
		rctx, err := NewUserUpdateContext(ctx)
		if err != nil {
			return err
		}
		return ctrl.Update(rctx)
	}
	service.PATCH("/api/user/:userID", h)
	service.Logger.Info("Mount ctrl User action Update")
}
