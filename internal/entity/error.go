package entity

import "errors"

// Errors
var (
	ErrNotFound             = errors.New("acourse: not found")
	ErrUsernameNotAvailable = errors.New("acourse: username not available")
	ErrEmailNotAvailable    = errors.New("acourse: email not available")
)
