package main

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/acoshift/configfile"
	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

var (
	bucketHandle *storage.BucketHandle
)

func main() {
	ctx := context.Background()
	config := configfile.NewReader("config")
	serviceAccount := config.Bytes("service_account")
	sqlURL := config.String("sql_url")

	gconf, err := google.JWTConfigFromJSON(serviceAccount, storage.ScopeReadWrite)
	must(err)
	storageClient, err := storage.NewClient(ctx, option.WithTokenSource(gconf.TokenSource(ctx)))
	must(err)
	bucketHandle = storageClient.Bucket("acourse")

	db, err := sql.Open("postgres", sqlURL)
	must(err)

	rows, err := db.Query(`select id, download_url from user_assignments where download_url <> ''`)
	must(err)
	for rows.Next() {
		var (
			id    string
			image string
		)
		err = rows.Scan(&id, &image)
		must(err)
		if !strings.HasPrefix(image, "https://storage.googleapis.com/acourse/upload") {
			log.Println(image)
			resp, err := http.Get(image)
			if err == nil {
				if resp.StatusCode == 200 {
					p, _ := url.Parse(image)
					fn := generateFilename() + path.Ext(p.Path)
					err = Upload(ctx, resp.Body, fn)
					must(err)
					resImage := generateDownloadURL(fn)
					_, err = db.Exec(`update user_assignments set download_url = $2 where id = $1`, id, resImage)
					must(err)
					log.Println(resImage)
				} else {
					log.Println(resp.StatusCode)
				}
				resp.Body.Close()
			} else {
				log.Println(err)
			}
		}
	}
	must(rows.Err())
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func resizeCropEncode(w io.Writer, r io.Reader, width, height int, quality int) error {
	m, _, err := image.Decode(r)
	if err != nil {
		return err
	}
	result := imaging.Thumbnail(m, width, height, imaging.Lanczos)
	return jpeg.Encode(w, result, &jpeg.Options{Quality: quality})
}

func resizeEncode(w io.Writer, r io.Reader, width int, quality int) error {
	m, _, err := image.Decode(r)
	if err != nil {
		return err
	}
	if m.Bounds().Dx() > width {
		m = imaging.Resize(m, width, 0, imaging.Lanczos)
	}
	return jpeg.Encode(w, m, &jpeg.Options{Quality: quality})
}

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
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", "acourse", filename)
}

// UploadPaymentImage uploads payment image
func UploadPaymentImage(ctx context.Context, r io.Reader) (string, error) {
	buf := &bytes.Buffer{}
	err := resizeEncode(buf, r, 700, 60)
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

// UploadProfileFromURLAsync copies data from given url and upload profile in background,
// returns url of destination file
func UploadProfileFromURLAsync(url string) string {
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
		err = Upload(ctx, buf, filename)
		if err != nil {
			return
		}
	}()
	return generateDownloadURL(filename)
}

// UploadCourseCoverImage uploads course cover image
func UploadCourseCoverImage(ctx context.Context, r io.Reader) (string, error) {
	buf := &bytes.Buffer{}
	err := resizeEncode(buf, r, 1200, 90)
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
