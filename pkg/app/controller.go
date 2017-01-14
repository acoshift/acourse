package app

import (
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
