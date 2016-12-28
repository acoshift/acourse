package app

import (
	"net/http"

	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/unrolled/render"
)

// ErrorReply is the error response
type ErrorReply struct {
	Error *Error `json:"error"`
}

// Error is the error holder
type Error struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (x *Error) Error() string {
	return fmt.Sprintf("%s: %s", x.Code, x.Message)
}

// SuccessReply is the success response without any content
type SuccessReply struct {
	OK int `json:"ok"`
}

// Reply templates
var (
	ErrNotFound     error = &Error{http.StatusNotFound, "not_found", "resource not found"}
	ErrUnauthorized error = &Error{http.StatusUnauthorized, "unauthorized", ""}
	ErrForbidden    error = &Error{http.StatusForbidden, "forbidden", ""}
)

// ErrorFunc is the error generator function which created for each error code
type ErrorFunc func(error) error

// CreateErrors is the helper function for create error func
func CreateErrors(status int, code string) ErrorFunc {
	return func(err error) error {
		return &Error{Status: status, Code: code, Message: err.Error()}
	}
}

var success = &SuccessReply{1}

var rr = render.New()

func handleOK(ctx *gin.Context, r interface{}) error {
	ctx.JSON(http.StatusOK, r)
	return nil
}

func handleError(ctx *gin.Context, r error) error {
	if err, ok := r.(*Error); ok {
		ctx.JSON(err.Status, &ErrorReply{err})
	} else {
		handleInternalServerError(ctx, r)
	}
	return nil
}

func handleNotFound(ctx *gin.Context) error {
	return handleError(ctx, ErrNotFound)
}

func handleSuccess(ctx *gin.Context) error {
	ctx.JSON(http.StatusOK, success)
	return nil
}

func handleUnauthorized(ctx *gin.Context) error {
	return handleError(ctx, ErrUnauthorized)
}

func handleForbidden(ctx *gin.Context) error {
	return handleError(ctx, ErrForbidden)
}

func handleBadRequest(ctx *gin.Context, r error) error {
	return handleError(ctx, &Error{http.StatusBadRequest, "bad_request", r.Error()})
}

func handleInternalServerError(ctx *gin.Context, r error) error {
	return handleError(ctx, &Error{http.StatusInternalServerError, "internal_server_error", r.Error()})
}

func handleHTML(ctx *gin.Context, name string, binding interface{}) error {
	return rr.HTML(ctx.Writer, http.StatusOK, name, binding)
}
