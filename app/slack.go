package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

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
