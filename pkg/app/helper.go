package app

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"golang.org/x/net/xsrftoken"
)

func generateSessionID() string {
	b := make([]byte, 24)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func verifyXSRF(token, userID, action string) bool {
	return xsrftoken.Valid(token, xsrfSecret, userID, action)
}

func back(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, r.RequestURI, http.StatusSeeOther)
}

// func hashPassword(password string) (string, error) {
// 	hpwd, err := bcrypt.GenerateFromPassword([]byte(password), 10)
// 	if err != nil {
// 		return "", err
// 	}
// 	return string(hpwd), nil
// }

// func verifyPassword(hpwd, password string) bool {
// 	return bcrypt.CompareHashAndPassword([]byte(hpwd), []byte(password)) == nil
// }
