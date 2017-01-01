package app

import (
	"github.com/acoshift/acourse/pkg/payload"
	"gopkg.in/gin-gonic/gin.v1"
)

// UserShowContext provides the user show action context
type UserShowContext struct {
	CurrentUserID string
	UserID        string
}

// NewUserShowContext parses the incoming request and create context
func NewUserShowContext(ctx *gin.Context) *UserShowContext {
	rctx := UserShowContext{}
	if v, ok := ctx.Get(keyCurrentUserID); ok {
		rctx.CurrentUserID, _ = v.(string)
	}
	rctx.UserID = ctx.Param("userID")
	return &rctx
}

// UserUpdateContext provides the user update action context
type UserUpdateContext struct {
	CurrentUserID string
	UserID        string
	Payload       *payload.User
}

// NewUserUpdateContext parses the incoming request and create context
func NewUserUpdateContext(ctx *gin.Context) *UserUpdateContext {
	rctx := UserUpdateContext{}
	if v, ok := ctx.Get(keyCurrentUserID); ok {
		rctx.CurrentUserID, _ = v.(string)
	}
	rctx.UserID = ctx.Param("userID")
	return &rctx
}

// HealthHealthContext provides the health health action context
type HealthHealthContext struct{}

// NewHealthHealthContext parses the incoming request and create context
func NewHealthHealthContext(ctx *gin.Context) *HealthHealthContext {
	return &HealthHealthContext{}
}

// CourseShowContext provides the course show action context
type CourseShowContext struct {
	CurrentUserID string
	CourseID      string
	Own           *bool
	Student       *string
}

// NewCourseShowContext parses the incoming request and create context
func NewCourseShowContext(ctx *gin.Context) *CourseShowContext {
	rctx := CourseShowContext{}
	if v, ok := ctx.Get(keyCurrentUserID); ok {
		rctx.CurrentUserID, _ = v.(string)
	}
	rctx.CourseID = ctx.Param("courseID")
	return &rctx
}

// CourseCreateContext provides the course create action context
type CourseCreateContext struct {
	CurrentUserID string
	Payload       *payload.Course
}

// NewCourseCreateContext parses the incoming request and create context
func NewCourseCreateContext(ctx *gin.Context) *CourseCreateContext {
	rctx := CourseCreateContext{}
	if v, ok := ctx.Get(keyCurrentUserID); ok {
		rctx.CurrentUserID, _ = v.(string)
	}
	return &rctx
}

// CourseUpdateContext provides the course update action context
type CourseUpdateContext struct {
	CurrentUserID string
	CourseID      string
	Payload       *payload.Course
}

// NewCourseUpdateContext parses the incoming request and create context
func NewCourseUpdateContext(ctx *gin.Context) *CourseUpdateContext {
	rctx := CourseUpdateContext{}
	if v, ok := ctx.Get(keyCurrentUserID); ok {
		rctx.CurrentUserID, _ = v.(string)
	}
	rctx.CourseID = ctx.Param("courseID")
	return &rctx
}

// CourseListContext provides the course list action context
type CourseListContext struct {
	CurrentUserID string
	Owner         string
	Student       string
}

// NewCourseListContext parses the incoming request and create context
func NewCourseListContext(ctx *gin.Context) *CourseListContext {
	rctx := CourseListContext{}
	if v, ok := ctx.Get(keyCurrentUserID); ok {
		rctx.CurrentUserID, _ = v.(string)
	}
	rctx.Owner = ctx.Query("owner")
	rctx.Student = ctx.Query("student")
	return &rctx
}

// CourseEnrollContext provides the course enroll action context
type CourseEnrollContext struct {
	CurrentUserID string
	CourseID      string
	Payload       *payload.CourseEnroll
}

// NewCourseEnrollContext parses the incoming request and create context
func NewCourseEnrollContext(ctx *gin.Context) *CourseEnrollContext {
	rctx := CourseEnrollContext{}
	if v, ok := ctx.Get(keyCurrentUserID); ok {
		rctx.CurrentUserID, _ = v.(string)
	}
	rctx.CourseID = ctx.Param("courseID")
	return &rctx
}

// PaymentListContext provides the payment list action context
type PaymentListContext struct {
	CurrentUserID string
	History       bool
}

// NewPaymentListContext parses the incoming request and create context
func NewPaymentListContext(ctx *gin.Context) *PaymentListContext {
	rctx := PaymentListContext{}
	if v, ok := ctx.Get(keyCurrentUserID); ok {
		rctx.CurrentUserID, _ = v.(string)
	}
	if ctx.Query("history") == "true" {
		rctx.History = true
	}
	return &rctx
}

// PaymentApproveContext provides the payment approve action context
type PaymentApproveContext struct {
	CurrentUserID string
	PaymentID     string
}

// NewPaymentApproveContext parses the incoming request and create context
func NewPaymentApproveContext(ctx *gin.Context) *PaymentApproveContext {
	rctx := PaymentApproveContext{}
	if v, ok := ctx.Get(keyCurrentUserID); ok {
		rctx.CurrentUserID, _ = v.(string)
	}
	rctx.PaymentID = ctx.Param("paymentID")
	return &rctx
}

// PaymentRejectContext provides the payment reject action context
type PaymentRejectContext struct {
	CurrentUserID string
	PaymentID     string
}

// NewPaymentRejectContext parses the incoming request and create context
func NewPaymentRejectContext(ctx *gin.Context) *PaymentRejectContext {
	rctx := PaymentRejectContext{}
	if v, ok := ctx.Get(keyCurrentUserID); ok {
		rctx.CurrentUserID, _ = v.(string)
	}
	rctx.PaymentID = ctx.Param("paymentID")
	return &rctx
}

// RenderIndexContext provides the render index action context
// use for render static file
type RenderIndexContext struct{}

// NewRenderIndexContext parses the incoming request and create context
func NewRenderIndexContext(ctx *gin.Context) *RenderIndexContext {
	return &RenderIndexContext{}
}

// RenderCourseContext provides the render course action context
type RenderCourseContext struct {
	CourseID string
}

// NewRenderCourseContext parses the incoming request and create context
func NewRenderCourseContext(ctx *gin.Context) *RenderCourseContext {
	rctx := RenderCourseContext{}
	rctx.CourseID = ctx.Param("courseID")
	return &rctx
}
