package app

import (
	"github.com/acoshift/httperror"
	"gopkg.in/gin-gonic/gin.v1"
)

// HealthController is the controller interface for health check
type HealthController interface {
	Check() error
}

// MountHealthController mounts a Health controller to the http server
func MountHealthController(server *gin.Engine, c HealthController) {
	server.GET("/_ah/health", func(ctx *gin.Context) {
		err := c.Check()
		if err != nil {
			handleError(ctx, err)
			return
		}
		handleSuccess(ctx)
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
			handleError(ctx, httperror.Unauthorized)
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
			handleError(ctx, httperror.Unauthorized)
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
			handleError(ctx, httperror.Unauthorized)
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
