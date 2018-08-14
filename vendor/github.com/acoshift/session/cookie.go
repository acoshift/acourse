package session

import (
	"net/http"
)

type cookie struct {
	http.Cookie
	SameSite SameSite
}

func setCookie(w http.ResponseWriter, c cookie) {
	if v := c.String(); v != "" {
		if len(c.SameSite) > 0 {
			v += "; SameSite=" + string(c.SameSite)
		}
		w.Header().Add("Set-Cookie", v)
	}
}
