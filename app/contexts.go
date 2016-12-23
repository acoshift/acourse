package app

import "github.com/labstack/echo"
import "net/http"

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

// UserShowContext provides the user show action context
type UserShowContext struct {
	context       echo.Context
	CurrentUserID string
	UserID        string
}

// NewUserShowContext parses the incoming request and create context
func NewUserShowContext(ctx echo.Context) (*UserShowContext, error) {
	var err error
	rctx := UserShowContext{context: ctx}
	rctx.UserID = ctx.Param("userID")
	return &rctx, err
}

// NotFound sends a HTTP response
func (ctx *UserShowContext) NotFound() error {
	return handleNotFound(ctx.context)
}

// InternalServerError sends a HTTP response
func (ctx *UserShowContext) InternalServerError(r error) error {
	return handleError(ctx.context, r)
}

// OK sends a HTTP response
func (ctx *UserShowContext) OK(r *UserView) error {
	return handleOK(ctx.context, r)
}

// OKMe send a HTTP response
func (ctx *UserShowContext) OKMe(r *UserMeView) error {
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
	var err error
	rctx := UserUpdateContext{context: ctx}
	rctx.UserID = ctx.Param("userID")
	return &rctx, err
}

// NoContent sends a HTTP response
func (ctx *UserUpdateContext) NoContent() error {
	return handleNoContent(ctx.context)
}

// NotFound sends a HTTP response
func (ctx *UserUpdateContext) NotFound() error {
	return handleNotFound(ctx.context)
}

// InternalServerError sends a HTTP response
func (ctx *UserUpdateContext) InternalServerError(r error) error {
	return handleError(ctx.context, r)
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
	rctx.CourseID = ctx.Param("courseID")
	rctx.CurrentUserID, _ = ctx.Get(keyCurrentUserID).(string)
	return &rctx, nil
}

// OK sends HTTP response
func (ctx *CourseShowContext) OK(r *CourseView) error {
	return handleOK(ctx.context, r)
}

// OKPublic sends HTTP response
func (ctx *CourseShowContext) OKPublic(r *CoursePublicView) error {
	return handleOK(ctx.context, r)
}

// NotFound sends HTTP response
func (ctx *CourseShowContext) NotFound() error {
	return handleNotFound(ctx.context)
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

// CourseListContext provides the course list action context
type CourseListContext struct {
	context echo.Context
}

// NewCourseListContext parses the incoming request and create context
func NewCourseListContext(ctx echo.Context) (*CourseListContext, error) {
	rctx := CourseListContext{context: ctx}
	return &rctx, nil
}

// OKTiny sends HTTP response
func (ctx *CourseListContext) OKTiny(r CourseTinyCollectionView) error {
	return handleOK(ctx.context, r)
}
