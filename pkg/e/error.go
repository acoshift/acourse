package e

import (
	"fmt"
	"net/http"
)

// Error is the error type
type Error struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ErrorFunc is the error generator function which created for each error code
type ErrorFunc func(error) error

// New is the helper function for create error func
func New(status int, code string) ErrorFunc {
	return func(err error) error {
		return &Error{Status: status, Code: code, Message: err.Error()}
	}
}

// Error predefined
var (
	ErrNotFound     = &Error{http.StatusNotFound, "not_found", "Not found"}
	ErrUnauthorized = &Error{http.StatusUnauthorized, "unauthorized", ""}
	ErrForbidden    = &Error{http.StatusForbidden, "forbidden", ""}
	ErrBadRequest   = &Error{http.StatusBadRequest, "bad_request", ""}
)

// Error implement error interface
func (x *Error) Error() string {
	return fmt.Sprintf("error: [%d] %s | %s", x.Status, x.Code, x.Message)
}

// BadRequest builds error for bad request
func BadRequest(err error) error {
	return &Error{http.StatusBadRequest, "bad_request", err.Error()}
}

// InternalError builds error for internal server error
func InternalError(err error) error {
	return &Error{http.StatusInternalServerError, "internal_server_error", err.Error()}
}

// Conflict builds error for conflict
func Conflict(err error) error {
	return &Error{http.StatusConflict, "conflict", err.Error()}
}
