package app

import "fmt"

// Error type
type Error struct {
	ID     string `json:"id"`
	Status int    `json:"status"`
	Code   string `json:"code"`
	Detail string `json:"detail"`
}

// ErrorCreateFunc is the creator function for create error
type ErrorCreateFunc func(interface{}) error

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
	return func(detail interface{}) error {
		d := ""
		switch p := detail.(type) {
		case error:
			d = p.Error()
		case string:
			d = p
		default:
			d = "unknown"
		}
		return &Error{
			Status: status,
			Code:   code,
			Detail: d,
		}
	}
}

func createInternalError(r error, status int, code string) error {
	if _, ok := r.(*Error); ok {
		return r
	}
	if code == "" {
		code = "unknown"
	}
	return &Error{
		Status: status,
		Code:   code,
		Detail: r.Error(),
	}
}

// Error implements error interface
func (x *Error) Error() string {
	return fmt.Sprintf("error: %d (%s), %s, %s", x.Status, x.ID, x.Code, x.Detail)
}
