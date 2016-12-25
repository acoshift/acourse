package app

import (
	"io"
	"net/http"

	"github.com/labstack/echo"
	"github.com/unrolled/render"
)

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
	Create(*CourseCreateContext) error
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

	service.POST("/api/course", func(ctx echo.Context) error {
		rctx, err := NewCourseCreateContext(ctx)
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
		return ctrl.Create(rctx)
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

	service.PUT("/api/payment/:paymentID/approve", func(ctx echo.Context) error {
		rctx, err := NewPaymentApproveContext(ctx)
		if err != nil {
			return handleError(ctx, err)
		}
		if rctx.CurrentUserID == "" {
			return handleUnauthorized(ctx)
		}
		return ctrl.Approve(rctx)
	})

	service.PUT("/api/payment/:paymentID/reject", func(ctx echo.Context) error {
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

// RenderController is the controller interface for render actions
type RenderController interface {
	Index(*RenderIndexContext) error
	Course(*RenderCourseContext) error
}

// Render wraps render.Render
type Render struct {
	r *render.Render
}

// Render implements echo.Renderer
func (r *Render) Render(w io.Writer, name string, data interface{}, ctx echo.Context) error {
	return r.r.HTML(w, http.StatusOK, name, data)
}

// MountRenderController mount a Render template controller on the given resource
func MountRenderController(service *echo.Echo, ctrl RenderController) {
	service.Renderer = &Render{r: render.New()}

	service.Static("/static", "public")

	service.File("/favicon.ico", "public/acourse-120.png")

	service.GET("/course/:courseID", func(ctx echo.Context) error {
		rctx, err := NewRenderCourseContext(ctx)
		if err != nil {
			return err
		}
		return ctrl.Course(rctx)
	})

	h := func(ctx echo.Context) error {
		rctx, err := NewRenderIndexContext(ctx)
		if err != nil {
			return err
		}
		return ctrl.Index(rctx)
	}

	service.GET("*", h)
	service.GET("/course/:courseID/edit", h)
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
