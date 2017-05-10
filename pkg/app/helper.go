package app

import (
	"net/url"
	"strings"

	"github.com/acoshift/acourse/pkg/internal"
	"golang.org/x/net/xsrftoken"
)

func extractPathID(u *url.URL) string {
	ps := strings.Split(u.Path[1:], "/")
	if len(ps) == 0 {
		return ""
	}
	return ps[0]
}

func verifyXSRF(token, userID, action string) bool {
	return xsrftoken.Valid(token, internal.GetXSRFSecret(), userID, action)
}
