package user

import (
	"errors"
)

var (
	ErrNotFound = errors.New("user: not found")
)
