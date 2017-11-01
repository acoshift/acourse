package session

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
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

func generateID() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		// this should never happended
		// or something wrong with OS's crypto pseudorandom generator
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(b)
}
