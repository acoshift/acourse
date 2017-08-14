package app

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"gopkg.in/gomail.v2"
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

func sendEmail(to string, subject, body string) error {
	if len(to) == 0 {
		return fmt.Errorf("invalid to")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", emailFrom)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	err := emailDialer.DialAndSend(m)
	if err != nil {
		return err
	}
	return nil
}

// markdown converts markdown to html for "email"
func markdown(s string) string {
	renderer := blackfriday.HtmlRenderer(
		0|
			blackfriday.HTML_USE_XHTML|
			blackfriday.HTML_USE_SMARTYPANTS|
			blackfriday.HTML_SMARTYPANTS_FRACTIONS|
			blackfriday.HTML_SMARTYPANTS_DASHES|
			blackfriday.HTML_SMARTYPANTS_LATEX_DASHES,
		"", "")
	md := blackfriday.MarkdownOptions([]byte(s), renderer, blackfriday.Options{
		Extensions: 0 |
			blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
			blackfriday.EXTENSION_FENCED_CODE |
			blackfriday.EXTENSION_AUTOLINK |
			blackfriday.EXTENSION_STRIKETHROUGH |
			blackfriday.EXTENSION_SPACE_HEADERS |
			blackfriday.EXTENSION_HEADER_IDS |
			blackfriday.EXTENSION_BACKSLASH_LINE_BREAK |
			blackfriday.EXTENSION_DEFINITION_LISTS,
	})
	p := bluemonday.UGCPolicy()
	return string(p.SanitizeBytes(md))
}

func strID(id int64) string {
	return strconv.FormatInt(id, 10)
}

func intID(id string) int64 {
	r, _ := strconv.ParseInt(id, 10, 64)
	return r
}
