package app

import (
	"github.com/acoshift/acourse/pkg/e"
	"github.com/acoshift/acourse/pkg/payload"
	"gopkg.in/gin-gonic/gin.v1"
)

// UserShowContext provides the user show action context
type UserShowContext struct {
	CurrentUserID string
	UserID        string
}

// NewUserShowContext parses the incoming request and create context
func NewUserShowContext(ctx *gin.Context) (*UserShowContext, error) {
	rctx := UserShowContext{}
	if v, ok := ctx.Get(keyCurrentUserID); ok {
		rctx.CurrentUserID, _ = v.(string)
	}
	rctx.UserID = ctx.Param("userID")
	return &rctx, nil
}

// UserUpdateContext provides the user update action context
type UserUpdateContext struct {
	CurrentUserID string
	UserID        string
	Payload       *payload.User
}

// NewUserUpdateContext parses the incoming request and create context
func NewUserUpdateContext(ctx *gin.Context) (*UserUpdateContext, error) {
	rctx := UserUpdateContext{}
	if v, ok := ctx.Get(keyCurrentUserID); ok {
		rctx.CurrentUserID, _ = v.(string)
	}
	rctx.UserID = ctx.Param("userID")
	var rp payload.RawUser
	err := ctx.BindJSON(&rp)
	if err != nil {
		return nil, e.BadRequest(err)
	}
	if err = rp.Validate(); err != nil {
		return nil, e.BadRequest(err)
	}
	rctx.Payload = rp.Payload()
	return &rctx, nil
}

// HealthHealthContext provides the health health action context
type HealthHealthContext struct{}

// NewHealthHealthContext parses the incoming request and create context
func NewHealthHealthContext(ctx *gin.Context) (*HealthHealthContext, error) {
	return &HealthHealthContext{}, nil
}

// CourseShowContext provides the course show action context
type CourseShowContext struct {
	CurrentUserID string
	CourseID      string
	Own           *bool
	Student       *string
}

// NewCourseShowContext parses the incoming request and create context
func NewCourseShowContext(ctx *gin.Context) (*CourseShowContext, error) {
	rctx := CourseShowContext{}
	if v, ok := ctx.Get(keyCurrentUserID); ok {
		rctx.CurrentUserID, _ = v.(string)
	}
	rctx.CourseID = ctx.Param("courseID")
	return &rctx, nil
}

// CourseCreateContext provides the course create action context
type CourseCreateContext struct {
	CurrentUserID string
	Payload       *payload.Course
}

// NewCourseCreateContext parses the incoming request and create context
func NewCourseCreateContext(ctx *gin.Context) (*CourseCreateContext, error) {
	rctx := CourseCreateContext{}
	if v, ok := ctx.Get(keyCurrentUserID); ok {
		rctx.CurrentUserID, _ = v.(string)
	}
	var rp payload.RawCourse
	err := ctx.Bind(&rp)
	if err != nil {
		return nil, e.BadRequest(err)
	}
	if err = rp.Validate(); err != nil {
		return nil, e.BadRequest(err)
	}
	rctx.Payload = rp.Payload()
	return &rctx, nil
}

// CourseUpdateContext provides the course update action context
type CourseUpdateContext struct {
	CurrentUserID string
	CourseID      string
	Payload       *payload.Course
}

// NewCourseUpdateContext parses the incoming request and create context
func NewCourseUpdateContext(ctx *gin.Context) (*CourseUpdateContext, error) {
	rctx := CourseUpdateContext{}
	if v, ok := ctx.Get(keyCurrentUserID); ok {
		rctx.CurrentUserID, _ = v.(string)
	}
	rctx.CourseID = ctx.Param("courseID")
	var rp payload.RawCourse
	err := ctx.Bind(&rp)
	if err != nil {
		return nil, e.BadRequest(err)
	}
	if err = rp.Validate(); err != nil {
		return nil, e.BadRequest(err)
	}
	rctx.Payload = rp.Payload()
	return &rctx, nil
}

// CourseListContext provides the course list action context
type CourseListContext struct {
	CurrentUserID string
	Owner         string
	Student       string
}

// NewCourseListContext parses the incoming request and create context
func NewCourseListContext(ctx *gin.Context) (*CourseListContext, error) {
	rctx := CourseListContext{}
	if v, ok := ctx.Get(keyCurrentUserID); ok {
		rctx.CurrentUserID, _ = v.(string)
	}
	rctx.Owner = ctx.Query("owner")
	rctx.Student = ctx.Query("student")
	return &rctx, nil
}

// CourseEnrollContext provides the course enroll action context
type CourseEnrollContext struct {
	CurrentUserID string
	CourseID      string
	Payload       *payload.CourseEnroll
}

// NewCourseEnrollContext parses the incoming request and create context
func NewCourseEnrollContext(ctx *gin.Context) (*CourseEnrollContext, error) {
	rctx := CourseEnrollContext{}
	if v, ok := ctx.Get(keyCurrentUserID); ok {
		rctx.CurrentUserID, _ = v.(string)
	}
	rctx.CourseID = ctx.Param("courseID")
	var rp payload.RawCourseEnroll
	err := ctx.Bind(&rp)
	if err != nil {
		return nil, e.BadRequest(err)
	}
	if err = rp.Validate(); err != nil {
		return nil, e.BadRequest(err)
	}
	rctx.Payload = rp.Payload()
	return &rctx, nil
}

// PaymentListContext provides the payment list action context
type PaymentListContext struct {
	CurrentUserID string
	History       bool
}

// NewPaymentListContext parses the incoming request and create context
func NewPaymentListContext(ctx *gin.Context) (*PaymentListContext, error) {
	rctx := PaymentListContext{}
	if v, ok := ctx.Get(keyCurrentUserID); ok {
		rctx.CurrentUserID, _ = v.(string)
	}
	if ctx.Query("history") == "true" {
		rctx.History = true
	}
	return &rctx, nil
}

// PaymentApproveContext provides the payment approve action context
type PaymentApproveContext struct {
	CurrentUserID string
	PaymentID     string
}

// NewPaymentApproveContext parses the incoming request and create context
func NewPaymentApproveContext(ctx *gin.Context) (*PaymentApproveContext, error) {
	rctx := PaymentApproveContext{}
	if v, ok := ctx.Get(keyCurrentUserID); ok {
		rctx.CurrentUserID, _ = v.(string)
	}
	rctx.PaymentID = ctx.Param("paymentID")
	return &rctx, nil
}

// PaymentRejectContext provides the payment reject action context
type PaymentRejectContext struct {
	CurrentUserID string
	PaymentID     string
}

// NewPaymentRejectContext parses the incoming request and create context
func NewPaymentRejectContext(ctx *gin.Context) (*PaymentRejectContext, error) {
	rctx := PaymentRejectContext{}
	if v, ok := ctx.Get(keyCurrentUserID); ok {
		rctx.CurrentUserID, _ = v.(string)
	}
	rctx.PaymentID = ctx.Param("paymentID")
	return &rctx, nil
}

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
