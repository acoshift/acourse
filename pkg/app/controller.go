package app

import (
	"github.com/acoshift/acourse/pkg/payload"
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
		rctx := NewUserShowContext(ctx)
		res, err := ctrl.Show(rctx)
		if err != nil {
			handleError(ctx, err)
		} else {
			handleOK(ctx, res)
		}
	})

	service.PATCH("/:userID", func(ctx *gin.Context) {
		rctx := NewUserUpdateContext(ctx)
		var rawPayload payload.RawUser
		err := ctx.BindJSON(&rawPayload)
		if err != nil {
			handleBadRequest(ctx, err)
			return
		}
		if err = rawPayload.Validate(); err != nil {
			handleBadRequest(ctx, err)
			return
		}
		rctx.Payload = rawPayload.Payload()
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
		rctx := NewHealthHealthContext(ctx)
		err := ctrl.Health(rctx)
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
		rctx := NewCourseListContext(ctx)
		res, err := ctrl.List(rctx)
		if err != nil {
			handleError(ctx, err)
		} else {
			handleOK(ctx, res)
		}
	})

	service.POST("", func(ctx *gin.Context) {
		rctx := NewCourseCreateContext(ctx)
		if rctx.CurrentUserID == "" {
			handleUnauthorized(ctx)
			return
		}
		var rawPayload payload.RawCourse
		err := ctx.Bind(&rawPayload)
		if err != nil {
			handleBadRequest(ctx, err)
			return
		}
		if err = rawPayload.Validate(); err != nil {
			handleBadRequest(ctx, err)
			return
		}
		rctx.Payload = rawPayload.Payload()
		res, err := ctrl.Create(rctx)
		if err != nil {
			handleError(ctx, err)
		} else {
			handleOK(ctx, res)
		}
	})

	service.GET("/:courseID", func(ctx *gin.Context) {
		rctx := NewCourseShowContext(ctx)
		res, err := ctrl.Show(rctx)
		if err != nil {
			handleError(ctx, err)
		} else {
			handleOK(ctx, res)
		}
	})

	service.PATCH("/:courseID", func(ctx *gin.Context) {
		rctx := NewCourseUpdateContext(ctx)
		if rctx.CurrentUserID == "" {
			handleUnauthorized(ctx)
			return
		}
		var rawPayload payload.RawCourse
		err := ctx.Bind(&rawPayload)
		if err != nil {
			handleBadRequest(ctx, err)
			return
		}
		if err = rawPayload.Validate(); err != nil {
			handleBadRequest(ctx, err)
			return
		}
		rctx.Payload = rawPayload.Payload()
		err = ctrl.Update(rctx)
		if err != nil {
			handleError(ctx, err)
		} else {
			handleSuccess(ctx)
		}
	})

	service.PUT("/:courseID/enroll", func(ctx *gin.Context) {
		rctx := NewCourseEnrollContext(ctx)
		if rctx.CurrentUserID == "" {
			handleUnauthorized(ctx)
			return
		}
		var rawPayload payload.RawCourseEnroll
		err := ctx.Bind(&rawPayload)
		if err != nil {
			handleBadRequest(ctx, err)
			return
		}
		if err = rawPayload.Validate(); err != nil {
			handleBadRequest(ctx, err)
			return
		}
		rctx.Payload = rawPayload.Payload()
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
		rctx := NewPaymentListContext(ctx)
		if rctx.CurrentUserID == "" {
			handleUnauthorized(ctx)
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
		rctx := NewPaymentApproveContext(ctx)
		if rctx.CurrentUserID == "" {
			handleUnauthorized(ctx)
			return
		}
		err := ctrl.Approve(rctx)
		if err != nil {
			handleError(ctx, err)
		} else {
			handleSuccess(ctx)
		}
	})

	service.PUT("/:paymentID/reject", func(ctx *gin.Context) {
		rctx := NewPaymentRejectContext(ctx)
		if rctx.CurrentUserID == "" {
			handleUnauthorized(ctx)
			return
		}
		err := ctrl.Reject(rctx)
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
		rctx := NewRenderCourseContext(ctx)
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
		rctx := NewRenderIndexContext(ctx)
		res, err := ctrl.Index(rctx)
		if err != nil {
			handleError(ctx, err)
		} else {
			handleHTML(ctx, "index", res)
		}
	}

	service.Use(h)
}
