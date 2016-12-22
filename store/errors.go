package store

import (
	"errors"

	"github.com/goadesign/goa"
)

// Errors
var (
	ErrInvalidID  = errors.New("invalid id")
	ErrInvalidUID = errors.New("invalid uid")
	ErrConflict   = goa.NewErrorClass("conflict", 409)
)
