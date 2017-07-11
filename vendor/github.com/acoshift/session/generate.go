package session

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
)

func generateID(entropy int) string {
	b := make([]byte, entropy)
	if _, err := rand.Read(b); err != nil {
		// this should never happended
		// or something wrong with OS's crypto pseudorandom generator
		panic(err)
	}
	return strings.TrimRight(base64.URLEncoding.EncodeToString(b), "=")
}
