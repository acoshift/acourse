package controller

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Upload upload files
func (c *ctrl) upload(ctx context.Context, r io.Reader, filename string) error {
	if len(filename) == 0 {
		return fmt.Errorf("invalid filename")
	}
	obj := c.bucketHandle.Object(filename)
	writer := obj.NewWriter(ctx)
	defer writer.Close()
	writer.CacheControl = "public, max-age=31536000"
	_, err := io.Copy(writer, r)
	if err != nil {
		return err
	}
	return nil
}

func generateFilename() string {
	return "upload/" + uuid.New().String()
}

func (c *ctrl) generateDownloadURL(filename string) string {
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", c.bucketName, filename)
}

// UploadPaymentImage uploads payment image
func (c *ctrl) uploadPaymentImage(ctx context.Context, r io.Reader) (string, error) {
	buf := &bytes.Buffer{}
	err := resizeEncode(buf, r, 700, 0, 60)
	if err != nil {
		return "", err
	}
	filename := generateFilename() + ".jpg"
	downloadURL := c.generateDownloadURL(filename)
	err = c.upload(ctx, buf, filename)
	if err != nil {
		return "", err
	}
	return downloadURL, nil
}

// UploadProfileImage uploads profile image and return url
func (c *ctrl) uploadProfileImage(ctx context.Context, r io.Reader) (string, error) {
	buf := &bytes.Buffer{}
	err := resizeCropEncode(buf, r, 500, 500, 90)
	if err != nil {
		return "", err
	}
	filename := generateFilename() + ".jpg"
	downloadURL := c.generateDownloadURL(filename)
	err = c.upload(ctx, buf, filename)
	if err != nil {
		return "", err
	}
	return downloadURL, nil
}

// UploadProfileFromURLAsync copies data from given url and upload profile in background,
// returns url of destination file
func (c *ctrl) uploadProfileFromURLAsync(url string) string {
	if len(url) == 0 {
		return ""
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return ""
	}
	filename := generateFilename() + ".jpg"
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
		err = resizeCropEncode(buf, resp.Body, 500, 500, 90)
		if err != nil {
			return
		}
		cancel()
		ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err = c.upload(ctx, buf, filename)
		if err != nil {
			return
		}
	}()
	return c.generateDownloadURL(filename)
}

// UploadCourseCoverImage uploads course cover image
func (c *ctrl) uploadCourseCoverImage(ctx context.Context, r io.Reader) (string, error) {
	buf := &bytes.Buffer{}
	err := resizeEncode(buf, r, 1200, 0, 90)
	if err != nil {
		return "", err
	}
	filename := generateFilename() + ".jpg"
	downloadURL := c.generateDownloadURL(filename)
	err = c.upload(ctx, buf, filename)
	if err != nil {
		return "", err
	}
	return downloadURL, nil
}
