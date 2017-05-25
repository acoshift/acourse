package middleware

import "net/http"

// Middleware is the http middleware
type Middleware func(http.Handler) http.Handler

// Chain is the helper function for chain middlewares into one middleware
func Chain(hs ...Middleware) Middleware {
	return func(h http.Handler) http.Handler {
		for i := len(hs); i > 0; i-- {
			h = hs[i-1](h)
		}
		return h
	}
}
