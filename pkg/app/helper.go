package app

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/acoshift/acourse/pkg/internal"
	"golang.org/x/crypto/bcrypt"
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

func back(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, r.RequestURI, http.StatusSeeOther)
}

func hashPassword(password string) (string, error) {
	hpwd, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(hpwd), nil
}

func verifyPassword(hpwd, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hpwd), []byte(password)) == nil
}
