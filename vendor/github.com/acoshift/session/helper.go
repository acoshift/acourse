package session

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strings"

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

func generateID() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		// this should never happended
		// or something wrong with OS's crypto pseudorandom generator
		panic(err)
	}
	return strings.TrimRight(base64.URLEncoding.EncodeToString(b), "=")
}
