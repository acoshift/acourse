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
			return handleError(ctx, ErrPayload(err))
		}
		if err = rawPayload.Validate(); err != nil {
			return handleError(ctx, ErrPayload(err))
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
	List(*CourseListContext) error
	Enroll(*CourseEnrollContext) error
}

// MountCourseController mounts a Course resource controller on the given service
func MountCourseController(service *echo.Echo, ctrl CourseController) {
	service.GET("/api/course", func(ctx echo.Context) error {
		rctx, err := NewCourseListContext(ctx)
		if err != nil {
			return handleError(ctx, err)
		}
		return ctrl.List(rctx)
	})

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
		if rctx.CurrentUserID == "" {
			return handleUnauthorized(ctx)
		}
		var rawPayload CourseRawPayload
		err = ctx.Bind(&rawPayload)
		if err != nil {
			return handleError(ctx, ErrPayload(err))
		}
		if err = rawPayload.Validate(); err != nil {
			return handleError(ctx, ErrPayload(err))
		}
		rctx.Payload = rawPayload.Payload()
		return ctrl.Update(rctx)
	})

	service.PUT("/api/course/:courseID/enroll", func(ctx echo.Context) error {
		rctx, err := NewCourseEnrollContext(ctx)
		if err != nil {
			return handleError(ctx, err)
		}
		if rctx.CurrentUserID == "" {
			return handleUnauthorized(ctx)
		}
		var rawPayload CourseEnrollRawPayload
		err = ctx.Bind(&rawPayload)
		if err != nil {
			return handleError(ctx, ErrPayload(err))
		}
		if err = rawPayload.Validate(); err != nil {
			return handleError(ctx, ErrPayload(err))
		}
		rctx.Payload = rawPayload.Payload()
		return ctrl.Enroll(rctx)
	})
}

// PaymentController is the controller interface for payment actions
type PaymentController interface {
	List(*PaymentListContext) error
	Approve(*PaymentApproveContext) error
	Reject(*PaymentRejectContext) error
}

// MountPaymentController mount a Payment resource controller on the given service
func MountPaymentController(service *echo.Echo, ctrl PaymentController) {
	service.GET("/api/payment", func(ctx echo.Context) error {
		rctx, err := NewPaymentListContext(ctx)
		if err != nil {
			return handleError(ctx, err)
		}
		if rctx.CurrentUserID == "" {
			return handleUnauthorized(ctx)
		}
		return ctrl.List(rctx)
	})

	service.PUT("/api/payment/approve", func(ctx echo.Context) error {
		rctx, err := NewPaymentApproveContext(ctx)
		if err != nil {
			return handleError(ctx, err)
		}
		if rctx.CurrentUserID == "" {
			return handleUnauthorized(ctx)
		}
		return ctrl.Approve(rctx)
	})

	service.PUT("/api/payment/reject", func(ctx echo.Context) error {
		rctx, err := NewPaymentRejectContext(ctx)
		if err != nil {
			return handleError(ctx, err)
		}
		if rctx.CurrentUserID == "" {
			return handleUnauthorized(ctx)
		}
		return ctrl.Reject(rctx)
	})
}

// UserIsAdminMiddleware is the middleware for autorization only admin user
// func UserIsAdminMiddleware(h echo.HandlerFunc) echo.HandlerFunc {
// 	return func(ctx echo.Context) error {
// 		userID, _ := ctx.Get(keyCurrentUserID).(string)
// 		if userID == "" {
// 			return handleForbidden(ctx)
// 		}

// 		return h(ctx)
// 	}
// }
