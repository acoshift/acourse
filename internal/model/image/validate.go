package image

import (
	"mime"
	"mime/multipart"

	"github.com/acoshift/header"

	"github.com/acoshift/acourse/internal/model/app"
)

var allowImageType = map[string]bool{
	"image/jpg":  true,
	"image/jpeg": true,
	"image/png":  true,
}

// Validate validates is file header an image
func Validate(img *multipart.FileHeader) (err error) {
	err = app.NewUIError("รองรับไฟล์ jpeg และ png เท่านั้น")

	if img == nil || img.Header == nil {
		return
	}

	ct, _, _ := mime.ParseMediaType(img.Header.Get(header.ContentType))

	if !allowImageType[ct] {
		return
	}

	return nil
}
