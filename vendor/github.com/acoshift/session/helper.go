package session

import (
	"net/http"

	"github.com/acoshift/header"
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
