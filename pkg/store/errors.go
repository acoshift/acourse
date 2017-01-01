package store

import (
	"errors"
	"fmt"
)

// Errors
var (
	ErrInvalidID  = errors.New("invalid id")
	ErrInvalidUID = errors.New("invalid uid")
)

// ErrConflict creates conflict error
func ErrConflict(detail string) error {
	return fmt.Errorf("conflict: %s", detail)
}
