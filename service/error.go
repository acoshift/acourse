package service

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

func newUIError(msg string) error {
	return &uiError{msg}
}
