package session

import "net/http"

const (
	prefixWWW = "www."

	headerXForwardedProto = "X-Forwarded-Proto"
)

func isTLS(r *http.Request) bool {
	if r.TLS != nil {
		return true
	}
	if r.Header.Get(headerXForwardedProto) == "https" {
		return true
	}
	return false
}
