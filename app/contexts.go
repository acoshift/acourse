package app

import "github.com/labstack/echo"
import "net/http"

// UserShowContext provides the user show action context
type UserShowContext struct {
	echo.Context
	UserID        string
	CurrentUserID string
}

// NewUserShowContext parses the incoming request and create context
func NewUserShowContext(ctx echo.Context) (*UserShowContext, error) {
	var err error
	rctx := UserShowContext{Context: ctx}
	rctx.UserID = ctx.Param("userID")
	return &rctx, err
}

// NotFound sends a HTTP response
func (ctx *UserShowContext) NotFound() error {
	return ctx.Context.NoContent(http.StatusNotFound)
}

// InternalServerError sends a HTTP response
func (ctx *UserShowContext) InternalServerError() error {
	return ctx.Context.NoContent(http.StatusInternalServerError)
}

// OK sends a HTTP response
func (ctx *UserShowContext) OK(r *UserView) error {
	return ctx.JSON(http.StatusOK, r)
}

// OKMe send a HTTP response
func (ctx *UserShowContext) OKMe(r *UserMeView) error {
	return ctx.JSON(http.StatusOK, r)
}

// UserUpdateContext provides the user update action context
type UserUpdateContext struct {
	echo.Context
	UserID  string
	Payload *UserPayload
}

// NewUserUpdateContext parses the incoming request and create context
func NewUserUpdateContext(ctx echo.Context) (*UserUpdateContext, error) {
	var err error
	rctx := UserUpdateContext{Context: ctx}
	rctx.UserID = ctx.Param("userID")
	err = ctx.Bind(&rctx.Payload)
	return &rctx, err
}
