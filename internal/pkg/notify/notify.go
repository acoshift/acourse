package notify

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/acoshift/acourse/internal/pkg/config"
)

var (
	client = &http.Client{
		Timeout: 10 * time.Second,
	}
	webhookURL = config.String("slack_url")
)

// Admin sends notify to admin
func Admin(message string) error {
	if webhookURL == "" {
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

	req, err := http.NewRequest(http.MethodPost, webhookURL, &buf)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()

	return nil
}
