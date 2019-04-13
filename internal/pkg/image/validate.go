package image

import (
	"errors"
	"mime"
	"mime/multipart"

	"github.com/acoshift/header"
)

var ErrInvalidType = errors.New("image: invalid type")

var allowImageType = map[string]bool{
	"image/jpg":  true,
	"image/jpeg": true,
	"image/png":  true,
}

// Validate validates is file header an image
func Validate(img *multipart.FileHeader) error {
	if img == nil || img.Header == nil {
		return ErrInvalidType
	}

	ct, _, _ := mime.ParseMediaType(img.Header.Get(header.ContentType))
	if !allowImageType[ct] {
		return ErrInvalidType
	}

	return nil
}
