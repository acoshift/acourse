package app

import (
	"acourse/view"
	"net/http"

	"github.com/labstack/echo"
)

func handleError(ctx echo.Context, r error) error {
	switch r := r.(type) {
	case *Error:
		r.ID = ctx.Response().Header().Get("X-Request-Id")
		return ctx.JSON(r.Status, r)
	default:
		return handleError(ctx, createInternalError(r, http.StatusInternalServerError, "unknown"))
	}
}

func handleOK(ctx echo.Context, r interface{}) error {
	return ctx.JSON(http.StatusOK, r)
}

func handleNotFound(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotFound)
}

func handleNoContent(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNoContent)
}

func handleUnauthorized(ctx echo.Context) error {
	return ctx.NoContent(http.StatusUnauthorized)
}

func handleForbidden(ctx echo.Context) error {
	return ctx.NoContent(http.StatusForbidden)
}

// UserShowContext provides the user show action context
type UserShowContext struct {
	context       echo.Context
	CurrentUserID string
	UserID        string
}

// NewUserShowContext parses the incoming request and create context
func NewUserShowContext(ctx echo.Context) (*UserShowContext, error) {
	rctx := UserShowContext{context: ctx}
	rctx.CurrentUserID, _ = ctx.Get(keyCurrentUserID).(string)
	rctx.UserID = ctx.Param("userID")
	return &rctx, nil
}

// NotFound sends a HTTP response
func (ctx *UserShowContext) NotFound() error {
	return handleNotFound(ctx.context)
}

// OK sends a HTTP response
func (ctx *UserShowContext) OK(r *view.UserView) error {
	return handleOK(ctx.context, r)
}

// OKMe send a HTTP response
func (ctx *UserShowContext) OKMe(r *view.UserMeView) error {
	return handleOK(ctx.context, r)
}

// UserUpdateContext provides the user update action context
type UserUpdateContext struct {
	context       echo.Context
	CurrentUserID string
	UserID        string
	Payload       *UserPayload
}

// NewUserUpdateContext parses the incoming request and create context
func NewUserUpdateContext(ctx echo.Context) (*UserUpdateContext, error) {
	rctx := UserUpdateContext{context: ctx}
	rctx.CurrentUserID, _ = ctx.Get(keyCurrentUserID).(string)
	rctx.UserID = ctx.Param("userID")
	return &rctx, nil
}

// NoContent sends a HTTP response
func (ctx *UserUpdateContext) NoContent() error {
	return handleNoContent(ctx.context)
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
	context echo.Context
}

// NewHealthHealthContext parses the incoming request and create context
func NewHealthHealthContext(ctx echo.Context) (*HealthHealthContext, error) {
	rctx := HealthHealthContext{context: ctx}
	return &rctx, nil
}

// OK sends HTTP response
func (ctx *HealthHealthContext) OK(r string) error {
	return ctx.context.String(http.StatusOK, r)
}

// CourseShowContext provides the course show action context
type CourseShowContext struct {
	context       echo.Context
	CurrentUserID string
	CourseID      string
	Own           *bool
	Student       *string
}

// NewCourseShowContext parses the incoming request and create context
func NewCourseShowContext(ctx echo.Context) (*CourseShowContext, error) {
	rctx := CourseShowContext{context: ctx}
	rctx.CurrentUserID, _ = ctx.Get(keyCurrentUserID).(string)
	rctx.CourseID = ctx.Param("courseID")
	return &rctx, nil
}

// OK sends HTTP response
func (ctx *CourseShowContext) OK(r *view.CourseView) error {
	return handleOK(ctx.context, r)
}

// OKPublic sends HTTP response
func (ctx *CourseShowContext) OKPublic(r *view.CoursePublicView) error {
	return handleOK(ctx.context, r)
}

// NotFound sends HTTP response
func (ctx *CourseShowContext) NotFound() error {
	return handleNotFound(ctx.context)
}

// CourseCreateContext provides the course create action context
type CourseCreateContext struct {
	context       echo.Context
	CurrentUserID string
	Payload       *CoursePayload
}

// NewCourseCreateContext parses the incoming request and create context
func NewCourseCreateContext(ctx echo.Context) (*CourseCreateContext, error) {
	rctx := CourseCreateContext{context: ctx}
	rctx.CurrentUserID, _ = ctx.Get(keyCurrentUserID).(string)
	return &rctx, nil
}

// OK sends HTTP response
func (ctx *CourseCreateContext) OK(r *view.CourseView) error {
	return handleOK(ctx.context, r)
}

// Forbidden sends HTTP response
func (ctx *CourseCreateContext) Forbidden() error {
	return handleForbidden(ctx.context)
}

// CourseUpdateContext provides the course update action context
type CourseUpdateContext struct {
	context       echo.Context
	CurrentUserID string
	CourseID      string
	Payload       *CoursePayload
}

// NewCourseUpdateContext parses the incoming request and create context
func NewCourseUpdateContext(ctx echo.Context) (*CourseUpdateContext, error) {
	rctx := CourseUpdateContext{context: ctx}
	rctx.CurrentUserID, _ = ctx.Get(keyCurrentUserID).(string)
	rctx.CourseID = ctx.Param("courseID")
	return &rctx, nil
}

// NoContent sends HTTP response
func (ctx *CourseUpdateContext) NoContent() error {
	return handleNoContent(ctx.context)
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
	context       echo.Context
	CurrentUserID string
	Owner         string
	Student       string
}

// NewCourseListContext parses the incoming request and create context
func NewCourseListContext(ctx echo.Context) (*CourseListContext, error) {
	rctx := CourseListContext{context: ctx}
	rctx.CurrentUserID, _ = ctx.Get(keyCurrentUserID).(string)
	rctx.Owner = ctx.QueryParam("owner")
	rctx.Student = ctx.QueryParam("student")
	return &rctx, nil
}

// OKTiny sends HTTP response
func (ctx *CourseListContext) OKTiny(r view.CourseTinyCollectionView) error {
	return handleOK(ctx.context, r)
}

// CourseEnrollContext provides the course enroll action context
type CourseEnrollContext struct {
	context       echo.Context
	CurrentUserID string
	CourseID      string
	Payload       *CourseEnrollPayload
}

// NewCourseEnrollContext parses the incoming request and create context
func NewCourseEnrollContext(ctx echo.Context) (*CourseEnrollContext, error) {
	rctx := CourseEnrollContext{context: ctx}
	rctx.CurrentUserID, _ = ctx.Get(keyCurrentUserID).(string)
	rctx.CourseID = ctx.Param("courseID")
	return &rctx, nil
}

// NoContent sends HTTP response
func (ctx *CourseEnrollContext) NoContent() error {
	return handleNoContent(ctx.context)
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
	context       echo.Context
	CurrentUserID string
}

// NewPaymentListContext parses the incoming request and create context
func NewPaymentListContext(ctx echo.Context) (*PaymentListContext, error) {
	rctx := PaymentListContext{context: ctx}
	rctx.CurrentUserID, _ = ctx.Get(keyCurrentUserID).(string)
	return &rctx, nil
}

// OK sends HTTP response
func (ctx *PaymentListContext) OK(r view.PaymentCollectionView) error {
	return ctx.context.JSON(http.StatusOK, r)
}

// Forbidden sends HTTP response
func (ctx *PaymentListContext) Forbidden() error {
	return handleForbidden(ctx.context)
}

// PaymentApproveContext provides the payment approve action context
type PaymentApproveContext struct {
	context       echo.Context
	CurrentUserID string
	PaymentID     string
}

// NewPaymentApproveContext parses the incoming request and create context
func NewPaymentApproveContext(ctx echo.Context) (*PaymentApproveContext, error) {
	rctx := PaymentApproveContext{context: ctx}
	rctx.CurrentUserID, _ = ctx.Get(keyCurrentUserID).(string)
	rctx.PaymentID = ctx.Param("paymentID")
	return &rctx, nil
}

// Forbidden sends HTTP response
func (ctx *PaymentApproveContext) Forbidden() error {
	return handleForbidden(ctx.context)
}

// NotFound sends HTTP response
func (ctx *PaymentApproveContext) NotFound() error {
	return handleNotFound(ctx.context)
}

// NoContent sends HTTP response
func (ctx *PaymentApproveContext) NoContent() error {
	return handleNoContent(ctx.context)
}

// PaymentRejectContext provides the payment reject action context
type PaymentRejectContext struct {
	context       echo.Context
	CurrentUserID string
	PaymentID     string
}

// NewPaymentRejectContext parses the incoming request and create context
func NewPaymentRejectContext(ctx echo.Context) (*PaymentRejectContext, error) {
	rctx := PaymentRejectContext{context: ctx}
	rctx.CurrentUserID, _ = ctx.Get(keyCurrentUserID).(string)
	rctx.PaymentID = ctx.Param("paymentID")
	return &rctx, nil
}

// Forbidden sends HTTP response
func (ctx *PaymentRejectContext) Forbidden() error {
	return handleForbidden(ctx.context)
}

// NotFound sends HTTP response
func (ctx *PaymentRejectContext) NotFound() error {
	return handleNotFound(ctx.context)
}

// NoContent sends HTTP response
func (ctx *PaymentRejectContext) NoContent() error {
	return handleNoContent(ctx.context)
}

// RenderIndexContext provides the render index action context
// use for render static file
type RenderIndexContext struct {
	context echo.Context
}

// NewRenderIndexContext parses the incoming request and create context
func NewRenderIndexContext(ctx echo.Context) (*RenderIndexContext, error) {
	rctx := RenderIndexContext{context: ctx}
	return &rctx, nil
}

// OK sends HTTP response
func (ctx *RenderIndexContext) OK(r *view.RenderIndexView) error {
	return ctx.context.Render(http.StatusOK, "index", r)
}

// RenderCourseContext provides the render course action context
type RenderCourseContext struct {
	context  echo.Context
	CourseID string
}

// NewRenderCourseContext parses the incoming request and create context
func NewRenderCourseContext(ctx echo.Context) (*RenderCourseContext, error) {
	rctx := RenderCourseContext{context: ctx}
	rctx.CourseID = ctx.Param("courseID")
	return &rctx, nil
}

// OK sends HTTP response
func (ctx *RenderCourseContext) OK(r *view.RenderIndexView) error {
	return ctx.context.Render(http.StatusOK, "index", r)
}

// NotFound redirect to home page
func (ctx *RenderCourseContext) NotFound() error {
	return ctx.context.Redirect(http.StatusTemporaryRedirect, "/")
}
