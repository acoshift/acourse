package middleware

import (
	"net/http"
	"net/textproto"
)

func isTLS(r *http.Request) bool {
	if r.TLS != nil {
		return true
	}
	if r.Header.Get("X-Forwarded-Proto") == "https" {
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

func addHeaderIfNotExists(h http.Header, key, value string) {
	key = textproto.CanonicalMIMEHeaderKey(key)
	for _, v := range h[key] {
		if v == value {
			return
		}
	}
	h.Add(key, value)
}
