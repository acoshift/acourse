package model

import "errors"

type scanFunc func(...interface{}) error

// Errors
var (
	ErrNotFound = errors.New("not found")
)
