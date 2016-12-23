package app

import "github.com/labstack/echo"
import "net/http"

// Errors
var (
	ErrPayload = CreateErrors(http.StatusBadRequest, "payload")
)

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
			return handleError(ctx, err)
		}
		var rawPayload UserRawPayload
		err = ctx.Bind(&rawPayload)
		if err != nil {
			return handleError(ctx, ErrPayload(err.Error()))
		}
		if err = rawPayload.Validate(); err != nil {
			return handleError(ctx, ErrPayload(err.Error()))
		}
		rctx.Payload = rawPayload.Payload()
		return ctrl.Update(rctx)
	})
	service.Logger.Info("Mount ctrl User action Update")
}

// HealthController is the controller interface for the Health actions
type HealthController interface {
	Health(*HealthHealthContext) error
}

// MountHealthController mounts a Health resource controller on the given service
func MountHealthController(service *echo.Echo, ctrl HealthController) {
	service.GET("/_ah/health", func(ctx echo.Context) error {
		rctx, err := NewHealthHealthContext(ctx)
		if err != nil {
			return handleError(ctx, err)
		}
		return ctrl.Health(rctx)
	})
}

// CourseController is the controller interface for course actions
type CourseController interface {
	Show(*CourseShowContext) error
	Update(*CourseUpdateContext) error
}

// MountCourseController mounts a Course resource controller on the given service
func MountCourseController(service *echo.Echo, ctrl CourseController) {
	service.GET("/api/course/:courseID", func(ctx echo.Context) error {
		rctx, err := NewCourseShowContext(ctx)
		if err != nil {
			return handleError(ctx, err)
		}
		return ctrl.Show(rctx)
	})

	service.PATCH("/api/course/:courseID", func(ctx echo.Context) error {
		rctx, err := NewCourseUpdateContext(ctx)
		if err != nil {
			return handleError(ctx, err)
		}
		var rawPayload CourseRawPayload
		err = ctx.Bind(&rawPayload)
		if err != nil {
			return handleError(ctx, ErrPayload(err.Error()))
		}
		if err = rawPayload.Validate(); err != nil {
			return handleError(ctx, ErrPayload(err.Error()))
		}
		rctx.Payload = rawPayload.Payload()
		return ctrl.Update(rctx)
	})
}
