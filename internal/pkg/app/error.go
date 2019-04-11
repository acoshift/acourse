package app

// TODO: rename package

// IsUIError returns true if given error is ui error
func IsUIError(err error) bool {
	_, ok := err.(*uiError)
	return ok
}

type uiError struct {
	msg string
}

func (err *uiError) Error() string {
	return err.msg
}

// NewUIError creates new ui error
func NewUIError(msg string) error {
	return &uiError{msg}
}
