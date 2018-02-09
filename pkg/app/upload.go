package app

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
func upload(ctx context.Context, r io.Reader, filename string) error {
	if len(filename) == 0 {
		return fmt.Errorf("invalid filename")
	}
	obj := bucketHandle.Object(filename)
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

func generateDownloadURL(filename string) string {
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, filename)
}

// UploadPaymentImage uploads payment image
func uploadPaymentImage(ctx context.Context, r io.Reader) (string, error) {
	buf := &bytes.Buffer{}
	err := resizeEncode(buf, r, 700, 0, 60)
	if err != nil {
		return "", err
	}
	filename := generateFilename() + ".jpg"
	downloadURL := generateDownloadURL(filename)
	err = upload(ctx, buf, filename)
	if err != nil {
		return "", err
	}
	return downloadURL, nil
}

// uploadProfileImage uploads profile image and return url
func uploadProfileImage(ctx context.Context, r io.Reader) (string, error) {
	buf := &bytes.Buffer{}
	err := resizeCropEncode(buf, r, 500, 500, 90)
	if err != nil {
		return "", err
	}
	filename := generateFilename() + ".jpg"
	downloadURL := generateDownloadURL(filename)
	err = upload(ctx, buf, filename)
	if err != nil {
		return "", err
	}
	return downloadURL, nil
}

// uploadProfileFromURLAsync copies data from given url and upload profile in background,
// returns url of destination file
func uploadProfileFromURLAsync(url string) string {
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
		err = upload(ctx, buf, filename)
		if err != nil {
			return
		}
	}()
	return generateDownloadURL(filename)
}

// UploadCourseCoverImage uploads course cover image
func uploadCourseCoverImage(ctx context.Context, r io.Reader) (string, error) {
	buf := &bytes.Buffer{}
	err := resizeEncode(buf, r, 1200, 0, 90)
	if err != nil {
		return "", err
	}
	filename := generateFilename() + ".jpg"
	downloadURL := generateDownloadURL(filename)
	err = upload(ctx, buf, filename)
	if err != nil {
		return "", err
	}
	return downloadURL, nil
}
