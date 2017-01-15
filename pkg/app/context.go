package app

import (
	"gopkg.in/gin-gonic/gin.v1"
)

// RenderIndexContext provides the render index action context
// use for render static file
type RenderIndexContext struct{}

// NewRenderIndexContext parses the incoming request and create context
func NewRenderIndexContext(ctx *gin.Context) (*RenderIndexContext, error) {
	return &RenderIndexContext{}, nil
}

// RenderCourseContext provides the render course action context
type RenderCourseContext struct {
	CourseID string
}

// NewRenderCourseContext parses the incoming request and create context
func NewRenderCourseContext(ctx *gin.Context) (*RenderCourseContext, error) {
	rctx := RenderCourseContext{}
	rctx.CourseID = ctx.Param("courseID")
	return &rctx, nil
}
