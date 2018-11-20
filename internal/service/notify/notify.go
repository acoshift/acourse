package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/acoshift/acourse/internal/pkg/bus"
	"github.com/acoshift/acourse/internal/pkg/model/notify"
)

// Init registers notify service
func Init(webhookURL string) {
	s := outgoingWebhookNotifier{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		url: webhookURL,
	}

	bus.Register(s.admin)
}

type outgoingWebhookNotifier struct {
	client *http.Client
	url    string
}

func (n *outgoingWebhookNotifier) admin(_ context.Context, m *notify.Admin) error {
	if n.url == "" {
		return nil
	}

	payload := struct {
		Text string `json:"text"`
	}{m.Message}
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
