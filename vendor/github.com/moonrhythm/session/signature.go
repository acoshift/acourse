package session

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/subtle"
	"encoding/base64"
)

func sign(value string, key []byte) string {
	h := hmac.New(sha1.New, key)
	h.Write([]byte(value))
	digest := h.Sum(nil)
	return base64.RawURLEncoding.EncodeToString(digest)
}

func verify(value, digest string, keys [][]byte) bool {
	for _, k := range keys {
		if subtle.ConstantTimeCompare([]byte(digest), []byte(sign(value, k))) == 1 {
			return true
		}
	}
	return false
}
