package internal

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/google/uuid"
)

// Upload upload files
func Upload(ctx context.Context, r io.Reader, filename string) error {
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
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucket, filename)
}

// UploadProfileImage uploads profile image and return url
func UploadProfileImage(ctx context.Context, r io.Reader) (string, error) {
	buf := &bytes.Buffer{}
	err := resizeCropEncode(buf, r, 500, 500, 90)
	if err != nil {
		return "", err
	}
	filename := generateFilename() + ".jpg"
	downloadURL := generateDownloadURL(filename)
	err = Upload(ctx, buf, filename)
	if err != nil {
		return "", err
	}
	return downloadURL, nil
}

// UploadCourseCoverImage uploads course cover image
func UploadCourseCoverImage(ctx context.Context, r io.Reader) (string, error) {
	buf := &bytes.Buffer{}
	err := resizeCropEncode(buf, r, 1200, 700, 90)
	if err != nil {
		return "", err
	}
	filename := generateFilename() + ".jpg"
	downloadURL := generateDownloadURL(filename)
	err = Upload(ctx, buf, filename)
	if err != nil {
		return "", err
	}
	return downloadURL, nil
}
