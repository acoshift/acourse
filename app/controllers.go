package app

import "github.com/labstack/echo"

// UserController is the controller interface for the User actions
type UserController interface {
	Show(*UserShowContext) error
	Update(*UserUpdateContext) error
}

// MountUserController mounts a User resource controller on the given service
func MountUserController(service *echo.Echo, ctrl UserController) {
	service.GET("/api/user/:userID", func(ctx echo.Context) error {
		rctx, err := NewUserShowContext(ctx)
		if err != nil {
			return err
		}
		return ctrl.Show(rctx)
	})
	service.Logger.Info("Mount ctrl User action Show")

	service.PATCH("/api/user/:userID", func(ctx echo.Context) error {
		rctx, err := NewUserUpdateContext(ctx)
		if err != nil {
			return err
		}
		var rawPayload UserRawPayload
		err = ctx.Bind(&rawPayload)
		if err != nil {
			return rctx.BadRequest(err)
		}
		if err = rawPayload.Validate(); err != nil {
			return rctx.BadRequest(err)
		}
		rctx.Payload = rawPayload.Payload()
		return ctrl.Update(rctx)
	})
	service.Logger.Info("Mount ctrl User action Update")
}
