package app

import (
	"github.com/acoshift/acourse/pkg/e"
	"gopkg.in/gin-gonic/gin.v1"
)

// UserController is the controller interface for the User actions
type UserController interface {
	Show(*UserShowContext) (interface{}, error)
	Update(*UserUpdateContext) error
}

// MountUserController mounts a User resource controller on the given service
func MountUserController(service *gin.RouterGroup, ctrl UserController) {
	service.GET("/:userID", func(ctx *gin.Context) {
		rctx, err := NewUserShowContext(ctx)
		if err != nil {
			handleError(ctx, err)
			return
		}
		res, err := ctrl.Show(rctx)
		if err != nil {
			handleError(ctx, err)
		} else {
			handleOK(ctx, res)
		}
	})

	service.PATCH("/:userID", func(ctx *gin.Context) {
		rctx, err := NewUserUpdateContext(ctx)
		if err != nil {
			handleError(ctx, err)
			return
		}
		err = ctrl.Update(rctx)
		if err != nil {
			handleError(ctx, err)
		} else {
			handleSuccess(ctx)
		}
	})
}

// HealthController is the controller interface for the Health actions
type HealthController interface {
	Health(*HealthHealthContext) error
}

// MountHealthController mounts a Health resource controller on the given service
func MountHealthController(service *gin.RouterGroup, ctrl HealthController) {
	service.GET("/health", func(ctx *gin.Context) {
		rctx, err := NewHealthHealthContext(ctx)
		if err != nil {
			handleError(ctx, err)
			return
		}
		err = ctrl.Health(rctx)
		if err != nil {
			handleError(ctx, err)
		} else {
			handleSuccess(ctx)
		}
	})
}

// CourseController is the controller interface for course actions
type CourseController interface {
	Show(*CourseShowContext) (interface{}, error)
	Create(*CourseCreateContext) (interface{}, error)
	Update(*CourseUpdateContext) error
	List(*CourseListContext) (interface{}, error)
	Enroll(*CourseEnrollContext) error
}

// MountCourseController mounts a Course resource controller on the given service
func MountCourseController(service *gin.RouterGroup, ctrl CourseController) {
	service.GET("", func(ctx *gin.Context) {
		rctx, err := NewCourseListContext(ctx)
		if err != nil {
			handleError(ctx, err)
			return
		}
		res, err := ctrl.List(rctx)
		if err != nil {
			handleError(ctx, err)
		} else {
			handleOK(ctx, res)
		}
	})

	service.POST("", func(ctx *gin.Context) {
		rctx, err := NewCourseCreateContext(ctx)
		if err != nil {
			handleError(ctx, err)
			return
		}
		if rctx.CurrentUserID == "" {
			handleError(ctx, e.ErrUnauthorized)
			return
		}
		res, err := ctrl.Create(rctx)
		if err != nil {
			handleError(ctx, err)
		} else {
			handleOK(ctx, res)
		}
	})

	service.GET("/:courseID", func(ctx *gin.Context) {
		rctx, err := NewCourseShowContext(ctx)
		if err != nil {
			handleError(ctx, err)
			return
		}
		res, err := ctrl.Show(rctx)
		if err != nil {
			handleError(ctx, err)
		} else {
			handleOK(ctx, res)
		}
	})

	service.PATCH("/:courseID", func(ctx *gin.Context) {
		rctx, err := NewCourseUpdateContext(ctx)
		if err != nil {
			handleError(ctx, err)
			return
		}
		if rctx.CurrentUserID == "" {
			handleError(ctx, e.ErrUnauthorized)
			return
		}
		err = ctrl.Update(rctx)
		if err != nil {
			handleError(ctx, err)
		} else {
			handleSuccess(ctx)
		}
	})

	service.PUT("/:courseID/enroll", func(ctx *gin.Context) {
		rctx, err := NewCourseEnrollContext(ctx)
		if err != nil {
			handleError(ctx, err)
			return
		}
		if rctx.CurrentUserID == "" {
			handleError(ctx, e.ErrUnauthorized)
			return
		}
		err = ctrl.Enroll(rctx)
		if err != nil {
			handleError(ctx, err)
		} else {
			handleSuccess(ctx)
		}
	})
}

// PaymentController is the controller interface for payment actions
type PaymentController interface {
	List(*PaymentListContext) (interface{}, error)
	Approve(*PaymentApproveContext) error
	Reject(*PaymentRejectContext) error
}

// MountPaymentController mount a Payment resource controller on the given service
func MountPaymentController(service *gin.RouterGroup, ctrl PaymentController) {
	service.GET("", func(ctx *gin.Context) {
		rctx, err := NewPaymentListContext(ctx)
		if err != nil {
			handleError(ctx, err)
			return
		}
		if rctx.CurrentUserID == "" {
			handleError(ctx, e.ErrUnauthorized)
			return
		}
		res, err := ctrl.List(rctx)
		if err != nil {
			handleError(ctx, err)
		} else {
			handleOK(ctx, res)
		}
	})

	service.PUT("/:paymentID/approve", func(ctx *gin.Context) {
		rctx, err := NewPaymentApproveContext(ctx)
		if err != nil {
			handleError(ctx, err)
			return
		}
		if rctx.CurrentUserID == "" {
			handleError(ctx, e.ErrUnauthorized)
			return
		}
		err = ctrl.Approve(rctx)
		if err != nil {
			handleError(ctx, err)
		} else {
			handleSuccess(ctx)
		}
	})

	service.PUT("/:paymentID/reject", func(ctx *gin.Context) {
		rctx, err := NewPaymentRejectContext(ctx)
		if err != nil {
			handleError(ctx, err)
			return
		}
		if rctx.CurrentUserID == "" {
			handleError(ctx, e.ErrUnauthorized)
			return
		}
		err = ctrl.Reject(rctx)
		if err != nil {
			handleError(ctx, err)
		} else {
			handleSuccess(ctx)
		}
	})
}

// RenderController is the controller interface for render actions
type RenderController interface {
	Index(*RenderIndexContext) (interface{}, error)
	Course(*RenderCourseContext) (interface{}, error)
}

// MountRenderController mount a Render template controller on the given resource
func MountRenderController(service *gin.Engine, ctrl RenderController) {
	cc := func(ctx *gin.Context) {
		ctx.Header("Cache-Control", "public, max-age=31536000")
		ctx.Next()
	}

	service.Group("/static", cc).Static("", "public")

	service.StaticFile("/favicon.ico", "public/acourse-120.png")

	service.GET("/course/:courseID", func(ctx *gin.Context) {
		rctx, err := NewRenderCourseContext(ctx)
		if err != nil {
			handleError(ctx, err)
			return
		}
		res, err := ctrl.Course(rctx)
		if err != nil {
			handleError(ctx, err)
		} else if res == nil {
			handleRedirect(ctx, "/")
		} else {
			handleHTML(ctx, "index", res)
		}
	})

	h := func(ctx *gin.Context) {
		rctx, err := NewRenderIndexContext(ctx)
		if err != nil {
			handleError(ctx, err)
			return
		}
		res, err := ctrl.Index(rctx)
		if err != nil {
			handleError(ctx, err)
		} else {
			handleHTML(ctx, "index", res)
		}
	}

	service.Use(h)
}
