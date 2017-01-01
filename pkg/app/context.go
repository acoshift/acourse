package app

import (
	"net/http"

	"github.com/acoshift/acourse/pkg/payload"
	"github.com/acoshift/acourse/pkg/view"
	"gopkg.in/gin-gonic/gin.v1"
)

// UserShowContext provides the user show action context
type UserShowContext struct {
	context       *gin.Context
	CurrentUserID string
	UserID        string
}

// NewUserShowContext parses the incoming request and create context
func NewUserShowContext(ctx *gin.Context) *UserShowContext {
	rctx := UserShowContext{context: ctx}
	if v, ok := ctx.Get(keyCurrentUserID); ok {
		rctx.CurrentUserID, _ = v.(string)
	}
	rctx.UserID = ctx.Param("userID")
	return &rctx
}

// NotFound sends a HTTP response
func (ctx *UserShowContext) NotFound() error {
	return handleNotFound(ctx.context)
}

// OK sends a HTTP response
func (ctx *UserShowContext) OK(r *view.User) error {
	return handleOK(ctx.context, r)
}

// OKMe send a HTTP response
func (ctx *UserShowContext) OKMe(r *view.UserMe) error {
	return handleOK(ctx.context, r)
}

// UserUpdateContext provides the user update action context
type UserUpdateContext struct {
	context       *gin.Context
	CurrentUserID string
	UserID        string
	Payload       *payload.User
}

// NewUserUpdateContext parses the incoming request and create context
func NewUserUpdateContext(ctx *gin.Context) *UserUpdateContext {
	rctx := UserUpdateContext{context: ctx}
	if v, ok := ctx.Get(keyCurrentUserID); ok {
		rctx.CurrentUserID, _ = v.(string)
	}
	rctx.UserID = ctx.Param("userID")
	return &rctx
}

// OK sends a HTTP response
func (ctx *UserUpdateContext) OK() error {
	return handleSuccess(ctx.context)
}

// NotFound sends a HTTP response
func (ctx *UserUpdateContext) NotFound() error {
	return handleNotFound(ctx.context)
}

// Forbidden sends a HTTP response
func (ctx *UserUpdateContext) Forbidden() error {
	return handleForbidden(ctx.context)
}

// HealthHealthContext provides the health health action context
type HealthHealthContext struct {
	context *gin.Context
}

// NewHealthHealthContext parses the incoming request and create context
func NewHealthHealthContext(ctx *gin.Context) *HealthHealthContext {
	return &HealthHealthContext{ctx}
}

// OK sends HTTP response
func (ctx *HealthHealthContext) OK() error {
	return handleSuccess(ctx.context)
}

// CourseShowContext provides the course show action context
type CourseShowContext struct {
	context       *gin.Context
	CurrentUserID string
	CourseID      string
	Own           *bool
	Student       *string
}

// NewCourseShowContext parses the incoming request and create context
func NewCourseShowContext(ctx *gin.Context) *CourseShowContext {
	rctx := CourseShowContext{context: ctx}
	if v, ok := ctx.Get(keyCurrentUserID); ok {
		rctx.CurrentUserID, _ = v.(string)
	}
	rctx.CourseID = ctx.Param("courseID")
	return &rctx
}

// OK sends HTTP response
func (ctx *CourseShowContext) OK(r *view.Course) error {
	return handleOK(ctx.context, r)
}

// OKPublic sends HTTP response
func (ctx *CourseShowContext) OKPublic(r *view.CoursePublic) error {
	return handleOK(ctx.context, r)
}

// NotFound sends HTTP response
func (ctx *CourseShowContext) NotFound() error {
	return handleNotFound(ctx.context)
}

// CourseCreateContext provides the course create action context
type CourseCreateContext struct {
	context       *gin.Context
	CurrentUserID string
	Payload       *payload.Course
}

// NewCourseCreateContext parses the incoming request and create context
func NewCourseCreateContext(ctx *gin.Context) *CourseCreateContext {
	rctx := CourseCreateContext{context: ctx}
	if v, ok := ctx.Get(keyCurrentUserID); ok {
		rctx.CurrentUserID, _ = v.(string)
	}
	return &rctx
}

// OK sends HTTP response
func (ctx *CourseCreateContext) OK(r *view.Course) error {
	return handleOK(ctx.context, r)
}

// Forbidden sends HTTP response
func (ctx *CourseCreateContext) Forbidden() error {
	return handleForbidden(ctx.context)
}

// CourseUpdateContext provides the course update action context
type CourseUpdateContext struct {
	context       *gin.Context
	CurrentUserID string
	CourseID      string
	Payload       *payload.Course
}

// NewCourseUpdateContext parses the incoming request and create context
func NewCourseUpdateContext(ctx *gin.Context) *CourseUpdateContext {
	rctx := CourseUpdateContext{context: ctx}
	if v, ok := ctx.Get(keyCurrentUserID); ok {
		rctx.CurrentUserID, _ = v.(string)
	}
	rctx.CourseID = ctx.Param("courseID")
	return &rctx
}

// OK sends HTTP response
func (ctx *CourseUpdateContext) OK() error {
	return handleSuccess(ctx.context)
}

// NotFound sends HTTP response
func (ctx *CourseUpdateContext) NotFound() error {
	return handleNotFound(ctx.context)
}

// Forbidden sends HTTP response
func (ctx *CourseUpdateContext) Forbidden() error {
	return handleForbidden(ctx.context)
}

// CourseListContext provides the course list action context
type CourseListContext struct {
	context       *gin.Context
	CurrentUserID string
	Owner         string
	Student       string
}

// NewCourseListContext parses the incoming request and create context
func NewCourseListContext(ctx *gin.Context) *CourseListContext {
	rctx := CourseListContext{context: ctx}
	if v, ok := ctx.Get(keyCurrentUserID); ok {
		rctx.CurrentUserID, _ = v.(string)
	}
	rctx.Owner = ctx.Query("owner")
	rctx.Student = ctx.Query("student")
	return &rctx
}

// OKTiny sends HTTP response
func (ctx *CourseListContext) OKTiny(r view.CourseTinyCollection) error {
	return handleOK(ctx.context, r)
}

// CourseEnrollContext provides the course enroll action context
type CourseEnrollContext struct {
	context       *gin.Context
	CurrentUserID string
	CourseID      string
	Payload       *payload.CourseEnroll
}

// NewCourseEnrollContext parses the incoming request and create context
func NewCourseEnrollContext(ctx *gin.Context) *CourseEnrollContext {
	rctx := CourseEnrollContext{context: ctx}
	if v, ok := ctx.Get(keyCurrentUserID); ok {
		rctx.CurrentUserID, _ = v.(string)
	}
	rctx.CourseID = ctx.Param("courseID")
	return &rctx
}

// OK sends HTTP response
func (ctx *CourseEnrollContext) OK() error {
	return handleSuccess(ctx.context)
}

// Forbidden sends HTTP response
func (ctx *CourseEnrollContext) Forbidden() error {
	return handleForbidden(ctx.context)
}

// NotFound sends HTTP response
func (ctx *CourseEnrollContext) NotFound() error {
	return handleNotFound(ctx.context)
}

// PaymentListContext provides the payment list action context
type PaymentListContext struct {
	context       *gin.Context
	CurrentUserID string
	History       bool
}

// NewPaymentListContext parses the incoming request and create context
func NewPaymentListContext(ctx *gin.Context) *PaymentListContext {
	rctx := PaymentListContext{context: ctx}
	if v, ok := ctx.Get(keyCurrentUserID); ok {
		rctx.CurrentUserID, _ = v.(string)
	}
	if ctx.Query("history") == "true" {
		rctx.History = true
	}
	return &rctx
}

// OK sends HTTP response
func (ctx *PaymentListContext) OK(r view.PaymentCollection) error {
	return handleOK(ctx.context, r)
}

// Forbidden sends HTTP response
func (ctx *PaymentListContext) Forbidden() error {
	return handleForbidden(ctx.context)
}

// PaymentApproveContext provides the payment approve action context
type PaymentApproveContext struct {
	context       *gin.Context
	CurrentUserID string
	PaymentID     string
}

// NewPaymentApproveContext parses the incoming request and create context
func NewPaymentApproveContext(ctx *gin.Context) *PaymentApproveContext {
	rctx := PaymentApproveContext{context: ctx}
	if v, ok := ctx.Get(keyCurrentUserID); ok {
		rctx.CurrentUserID, _ = v.(string)
	}
	rctx.PaymentID = ctx.Param("paymentID")
	return &rctx
}

// Forbidden sends HTTP response
func (ctx *PaymentApproveContext) Forbidden() error {
	return handleForbidden(ctx.context)
}

// NotFound sends HTTP response
func (ctx *PaymentApproveContext) NotFound() error {
	return handleNotFound(ctx.context)
}

// OK sends HTTP response
func (ctx *PaymentApproveContext) OK() error {
	return handleSuccess(ctx.context)
}

// PaymentRejectContext provides the payment reject action context
type PaymentRejectContext struct {
	context       *gin.Context
	CurrentUserID string
	PaymentID     string
}

// NewPaymentRejectContext parses the incoming request and create context
func NewPaymentRejectContext(ctx *gin.Context) *PaymentRejectContext {
	rctx := PaymentRejectContext{context: ctx}
	if v, ok := ctx.Get(keyCurrentUserID); ok {
		rctx.CurrentUserID, _ = v.(string)
	}
	rctx.PaymentID = ctx.Param("paymentID")
	return &rctx
}

// Forbidden sends HTTP response
func (ctx *PaymentRejectContext) Forbidden() error {
	return handleForbidden(ctx.context)
}

// NotFound sends HTTP response
func (ctx *PaymentRejectContext) NotFound() error {
	return handleNotFound(ctx.context)
}

// OK sends HTTP response
func (ctx *PaymentRejectContext) OK() error {
	return handleSuccess(ctx.context)
}

// RenderIndexContext provides the render index action context
// use for render static file
type RenderIndexContext struct {
	context *gin.Context
}

// NewRenderIndexContext parses the incoming request and create context
func NewRenderIndexContext(ctx *gin.Context) *RenderIndexContext {
	return &RenderIndexContext{ctx}
}

// OK sends HTTP response
func (ctx *RenderIndexContext) OK(r *view.RenderIndex) error {
	return handleHTML(ctx.context, "index", r)
}

// RenderCourseContext provides the render course action context
type RenderCourseContext struct {
	context  *gin.Context
	CourseID string
}

// NewRenderCourseContext parses the incoming request and create context
func NewRenderCourseContext(ctx *gin.Context) *RenderCourseContext {
	rctx := RenderCourseContext{context: ctx}
	rctx.CourseID = ctx.Param("courseID")
	return &rctx
}

// OK sends HTTP response
func (ctx *RenderCourseContext) OK(r *view.RenderIndex) error {
	return handleHTML(ctx.context, "index", r)
}

// NotFound redirect to home page
func (ctx *RenderCourseContext) NotFound() error {
	ctx.context.Redirect(http.StatusTemporaryRedirect, "/")
	return nil
}
