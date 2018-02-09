package app

import (
	"image"
	"image/jpeg"
	"io"

	"github.com/disintegration/imaging"
)

func encodeJPEG(w io.Writer, m image.Image, q int) error {
	return jpeg.Encode(w, m, &jpeg.Options{Quality: q})
}

func resizeCropEncode(w io.Writer, r io.Reader, width, height int, quality int) error {
	m, _, err := image.Decode(r)
	if err != nil {
		return err
	}
	result := imaging.Thumbnail(m, width, height, imaging.Lanczos)
	return encodeJPEG(w, result, quality)
}

func resizeEncode(w io.Writer, r io.Reader, width, height int, quality int) error {
	m, _, err := image.Decode(r)
	if err != nil {
		return err
	}
	if height == 0 && m.Bounds().Dx() > width {
		m = imaging.Resize(m, width, 0, imaging.Lanczos)
	} else if width == 0 && m.Bounds().Dy() > height {
		m = imaging.Resize(m, 0, height, imaging.Lanczos)
	} else if width != 0 && height != 0 {
		m = imaging.Resize(m, width, height, imaging.Lanczos)
	}
	return encodeJPEG(w, m, quality)
}
