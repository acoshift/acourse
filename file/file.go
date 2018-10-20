package file

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
	"github.com/satori/go.uuid"
)

// Storage is file storage
type Storage interface {
	Store(ctx context.Context, r io.Reader, filename string) error
	DownloadURL(filename string) string
}

// NewGCS creates new gcs
func NewGCS(client *storage.Client, bucket string) Storage {
	return &gcs{
		Bucket:     client.Bucket(bucket),
		BucketName: bucket,
	}
}

// GenerateFilename generates new filename
func GenerateFilename() string {
	return "upload/" + uuid.NewV4().String()
}

type gcs struct {
	Bucket     *storage.BucketHandle
	BucketName string
}

func (s *gcs) Store(ctx context.Context, r io.Reader, filename string) (err error) {
	if len(filename) == 0 {
		return fmt.Errorf("invalid filename")
	}

	obj := s.Bucket.Object(filename)
	w := obj.NewWriter(ctx)
	defer func() {
		if err != nil {
			w.CloseWithError(err)
			return
		}
		err = w.Close()
	}()

	w.CacheControl = "public, max-age=31536000"

	_, err = io.Copy(w, r)
	return
}

func (s *gcs) DownloadURL(filename string) string {
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", s.BucketName, filename)
}
