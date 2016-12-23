package app

import (
	"fmt"

	"net/http"
)

// Error type
type Error struct {
	ID     string `json:"id"`
	Status int    `json:"status"`
	Code   string `json:"code"`
	Detail string `json:"detail"`
}

// ErrorCreateFunc is the creator function for create error
type ErrorCreateFunc func(string) error

// CreateError creates error type
func CreateError(status int, code, detail string) error {
	return &Error{
		Status: status,
		Code:   code,
		Detail: detail,
	}
}

// CreateErrors generate create error function used for create template error
func CreateErrors(status int, code string) ErrorCreateFunc {
	return func(detail string) error {
		return &Error{
			Status: status,
			Code:   code,
			Detail: detail,
		}
	}
}

func createInternalError(r error) error {
	return &Error{
		Status: http.StatusInternalServerError,
		Code:   "unknown",
		Detail: r.Error(),
	}
}

// Error implements error interface
func (x *Error) Error() string {
	return fmt.Sprintf("error: %d (%s), %s, %s", x.Status, x.ID, x.Code, x.Detail)
}
