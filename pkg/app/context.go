package app

import (
	"context"
	"net/http"
)

// RenderIndexContext provides the render index action context
// use for render static file
type RenderIndexContext struct {
	context.Context
}

// NewRenderIndexContext parses the incoming request and create context
func NewRenderIndexContext(r *http.Request) (*RenderIndexContext, error) {
	return &RenderIndexContext{r.Context()}, nil
}

// RenderCourseContext provides the render course action context
type RenderCourseContext struct {
	context.Context
	CourseID string
}

// NewRenderCourseContext parses the incoming request and create context
func NewRenderCourseContext(r *http.Request) (*RenderCourseContext, error) {
	rctx := RenderCourseContext{Context: r.Context()}
	rctx.CourseID = r.URL.Path
	return &rctx, nil
}
