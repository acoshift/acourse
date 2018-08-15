package notify

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// AdminNotifier type
type AdminNotifier interface {
	Notify(message string) error
}

// NewOutgoingWebhookAdminNotifier returns new admin notifier
func NewOutgoingWebhookAdminNotifier(webhookURL string) AdminNotifier {
	return &outgoingWebhookAdminNotifier{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		url: webhookURL,
	}
}

type outgoingWebhookAdminNotifier struct {
	client *http.Client
	url    string
}

func (n *outgoingWebhookAdminNotifier) Notify(message string) error {
	if n.url == "" {
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

	req, err := http.NewRequest(http.MethodPost, n.url, &buf)
	if err != nil {
		return err
	}

	resp, err := n.client.Do(req)
	if err != nil {
		return err
	}

	io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()

	return nil
}
