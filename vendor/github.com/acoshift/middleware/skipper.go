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

// AlwaysSkip always return true
func AlwaysSkip(*http.Request) bool {
	return true
}

// SkipHTTP skips http request
func SkipHTTP(r *http.Request) bool {
	return !isTLS(r)
}

// SkipHTTPS skips https request
func SkipHTTPS(r *http.Request) bool {
	return isTLS(r)
}

// SkipIf skips if b is true
func SkipIf(b bool) Skipper {
	return func(*http.Request) bool {
		return b
	}
}

// SkipUnless skips if b is false
func SkipUnless(b bool) Skipper {
	return SkipIf(!b)
}
