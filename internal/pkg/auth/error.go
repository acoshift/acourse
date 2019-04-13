package auth

import (
	"errors"
)

var (
	ErrEmailRequired        = errors.New("auth: email required")
	ErrEmailInvalid         = errors.New("auth: email invalid")
	ErrPasswordRequired     = errors.New("auth: password required")
	ErrPasswordInvalid      = errors.New("auth: password invalid")
	ErrInvalidProvider      = errors.New("auth: invalid provider")
	ErrInvalidCallbackURI   = errors.New("auth: invalid callback uri")
	ErrEmailNotAvailable    = errors.New("auth: email not available")
	ErrUsernameNotAvailable = errors.New("auth: username not available")
)
