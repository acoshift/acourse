package file

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"

	"github.com/acoshift/acourse/internal/pkg/dispatcher"
	"github.com/acoshift/acourse/internal/pkg/model/file"
)

// InitGCS registers file dispatcher with gcs strategy
func InitGCS(client *storage.Client, bucket string) {
	s := &gcs{
		Bucket:     client.Bucket(bucket),
		BucketName: bucket,
	}

	dispatcher.Register(s.store)
}

type gcs struct {
	Bucket     *storage.BucketHandle
	BucketName string
}

func (s *gcs) store(ctx context.Context, m *file.Store) error {
	if len(m.Filename) == 0 {
		return fmt.Errorf("invalid filename")
	}

	m.Result = s.downloadURL(m.Filename)

	f := func() error {
		obj := s.Bucket.Object(m.Filename)
		w := obj.NewWriter(ctx)
		defer w.Close()

		w.CacheControl = "max-age=31536000, immutable"

		_, err := io.Copy(w, m.Reader)
		return err
	}

	if m.Async {
		go f()
		return nil
	}

	return f()
}

func (s *gcs) downloadURL(filename string) string {
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", s.BucketName, filename)
}
