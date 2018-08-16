package app

import (
	"bytes"
	"context"
	"io"

	"github.com/acoshift/acourse/file"
)

// UploadPaymentImage uploads payment image
func uploadPaymentImage(ctx context.Context, r io.Reader) (string, error) {
	buf := &bytes.Buffer{}
	err := imageResizeEncoder.ResizeEncode(buf, r, 700, 0, 60, false)
	if err != nil {
		return "", err
	}
	filename := file.GenerateFilename() + ".jpg"
	downloadURL := fileStorage.DownloadURL(filename)
	err = fileStorage.Store(ctx, buf, filename)
	if err != nil {
		return "", err
	}
	return downloadURL, nil
}

// uploadProfileImage uploads profile image and return url
func uploadProfileImage(ctx context.Context, r io.Reader) (string, error) {
	buf := &bytes.Buffer{}
	err := imageResizeEncoder.ResizeEncode(buf, r, 500, 500, 90, true)
	if err != nil {
		return "", err
	}
	filename := file.GenerateFilename() + ".jpg"
	downloadURL := fileStorage.DownloadURL(filename)
	err = fileStorage.Store(ctx, buf, filename)
	if err != nil {
		return "", err
	}
	return downloadURL, nil
}

// UploadCourseCoverImage uploads course cover image
func uploadCourseCoverImage(ctx context.Context, r io.Reader) (string, error) {
	buf := &bytes.Buffer{}
	err := imageResizeEncoder.ResizeEncode(buf, r, 1200, 0, 90, false)
	if err != nil {
		return "", err
	}
	filename := file.GenerateFilename() + ".jpg"
	downloadURL := fileStorage.DownloadURL(filename)
	err = fileStorage.Store(ctx, buf, filename)
	if err != nil {
		return "", err
	}
	return downloadURL, nil
}
