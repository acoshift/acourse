package editor

import (
	"bytes"
	"context"
	"io"

	"github.com/acoshift/acourse/file"
)

// UploadCourseCoverImage uploads course cover image
func (c *ctrl) uploadCourseCoverImage(ctx context.Context, r io.Reader) (string, error) {
	buf := &bytes.Buffer{}
	err := c.ImageResizeEncoder.ResizeEncode(buf, r, 1200, 0, 90, false)
	if err != nil {
		return "", err
	}
	filename := file.GenerateFilename() + ".jpg"
	downloadURL := c.FileStorage.DownloadURL(filename)
	err = c.FileStorage.Store(ctx, buf, filename)
	if err != nil {
		return "", err
	}
	return downloadURL, nil
}
