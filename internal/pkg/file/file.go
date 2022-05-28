package file

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
	uuid "github.com/satori/go.uuid"

	"github.com/acoshift/acourse/internal/pkg/config"
)

var (
	bucketName string
	bucket     *storage.BucketHandle
)

func Init() {
	bucketName = config.String("bucket")
	bucket = config.StorageClient().Bucket(bucketName)
}

// GenerateFilename generates new filename
func GenerateFilename() string {
	return "upload/" + uuid.NewV4().String()
}

// Store stores file
func Store(ctx context.Context, r io.Reader, filename string, async bool) (string, error) {
	if len(filename) == 0 {
		return "", fmt.Errorf("invalid filename")
	}

	url := downloadURL(filename)

	f := func() error {
		obj := bucket.Object(filename)
		w := obj.NewWriter(ctx)
		defer w.Close()

		w.CacheControl = "max-age=31536000, immutable"

		_, err := io.Copy(w, r)
		return err
	}

	if async {
		go f()
		return url, nil
	}

	return url, f()
}

func downloadURL(filename string) string {
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, filename)
}
