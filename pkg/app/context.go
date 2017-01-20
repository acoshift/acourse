package app

import (
	"net/http"
)

// RenderIndexContext provides the render index action context
// use for render static file
type RenderIndexContext struct{}

// NewRenderIndexContext parses the incoming request and create context
func NewRenderIndexContext(r *http.Request) (*RenderIndexContext, error) {
	return &RenderIndexContext{}, nil
}

// RenderCourseContext provides the render course action context
type RenderCourseContext struct {
	CourseID string
}

// NewRenderCourseContext parses the incoming request and create context
func NewRenderCourseContext(r *http.Request) (*RenderCourseContext, error) {
	rctx := RenderCourseContext{}
	rctx.CourseID = r.URL.Path
	return &rctx, nil
}
