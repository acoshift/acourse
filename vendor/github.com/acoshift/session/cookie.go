package session

import (
	"net/http"
)

type cookie struct {
	http.Cookie
	SameSite SameSite
}

func setCookie(w http.ResponseWriter, cookie *cookie) {
	if v := cookie.String(); v != "" {
		if len(cookie.SameSite) > 0 {
			v += "; SameSite=" + string(cookie.SameSite)
		}
		w.Header().Add("Set-Cookie", v)
	}
}
