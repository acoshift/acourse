package image

import (
	"image"
	"image/jpeg"
	"io"

	"github.com/disintegration/imaging"
)

// JPEGResizeEncoder type
type JPEGResizeEncoder interface {
	ResizeEncode(w io.Writer, r io.Reader, width, height int, quality int, crop bool) error
}

// NewJPEGResizeEncoder creates new jpeg resize encoder
func NewJPEGResizeEncoder() JPEGResizeEncoder {
	return &jpegResizeEncoder{}
}

type jpegResizeEncoder struct {
}

func (jpegResizeEncoder) ResizeEncode(w io.Writer, r io.Reader, width, height int, quality int, crop bool) error {
	m, _, err := image.Decode(r)
	if err != nil {
		return err
	}

	if crop {
		m = imaging.Thumbnail(m, width, height, imaging.Lanczos)
	} else {
		if height == 0 && m.Bounds().Dx() > width {
			m = imaging.Resize(m, width, 0, imaging.Lanczos)
		} else if width == 0 && m.Bounds().Dy() > height {
			m = imaging.Resize(m, 0, height, imaging.Lanczos)
		} else if width != 0 && height != 0 {
			m = imaging.Resize(m, width, height, imaging.Lanczos)
		}
	}

	return jpeg.Encode(w, m, &jpeg.Options{Quality: quality})
}
