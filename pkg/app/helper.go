package app

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"io/ioutil"
)

func generateSessionID() string {
	b := make([]byte, 24)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func back(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, r.RequestURI, http.StatusSeeOther)
}

func sendSlackMessage(ctx context.Context, message string) error {
	if len(slackURL) == 0 {
		return nil
	}

	payload := struct {
		Text string `json:"text"`
	}{message}
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(&payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, slackURL, &buf)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("response not ok")
	}
	return nil
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
