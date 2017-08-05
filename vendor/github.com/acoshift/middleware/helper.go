package middleware

import (
	"net/http"

	"github.com/acoshift/header"
)

const (
	prefixWWW = "www."
)

func isTLS(r *http.Request) bool {
	if r.TLS != nil {
		return true
	}
	if r.Header.Get(header.XForwardedProto) == "https" {
		return true
	}
	return false
}

func scheme(r *http.Request) string {
	if isTLS(r) {
		return "https"
	}
	return "http"
}
