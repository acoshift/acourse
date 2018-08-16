package auth

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"github.com/acoshift/acourse/file"
)

// uploadProfileFromURLAsync copies data from given url and upload profile in background,
// returns url of destination file
func (c *ctrl) uploadProfileFromURLAsync(url string) string {
	if len(url) == 0 {
		return ""
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return ""
	}
	filename := file.GenerateFilename() + ".jpg"
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		req = req.WithContext(ctx)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return
		}
		defer resp.Body.Close()
		buf := &bytes.Buffer{}
		err = c.ImageResizeEncoder.ResizeEncode(buf, resp.Body, 500, 500, 90, true)
		if err != nil {
			return
		}
		cancel()
		ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err = c.FileStorage.Store(ctx, buf, filename)
		if err != nil {
			return
		}
	}()
	return c.FileStorage.DownloadURL(filename)
}
