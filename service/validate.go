package service

import (
	"mime"
	"mime/multipart"

	"github.com/acoshift/header"
)

var allowImageType = map[string]bool{
	"image/jpg":  true,
	"image/jpeg": true,
	"image/png":  true,
}

func validateImage(img *multipart.FileHeader) (err error) {
	err = newUIError("รองรับไฟล์ jpeg และ png เท่านั้น")

	if img == nil || img.Header == nil {
		return
	}

	ct, _, _ := mime.ParseMediaType(img.Header.Get(header.ContentType))

	if !allowImageType[ct] {
		return
	}

	return nil
}
