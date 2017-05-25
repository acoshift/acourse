package middleware

import (
	"net/http"
)

// Skipper is the function to skip middleware,
// return true will skip middleware
type Skipper func(*http.Request) bool

// DefaultSkipper always return false
func DefaultSkipper(*http.Request) bool {
	return false
}
