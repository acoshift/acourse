package image

import (
	"image"
	"image/jpeg"
	"io"

	"github.com/disintegration/imaging"
)

func JPEG(w io.Writer, r io.Reader, width, height, quality int, crop bool) error {
	img, _, err := image.Decode(r)
	if err != nil {
		return err
	}

	if crop {
		img = imaging.Thumbnail(img, width, height, imaging.Lanczos)
	} else {
		if height == 0 && img.Bounds().Dx() > width {
			img = imaging.Resize(img, width, 0, imaging.Lanczos)
		} else if width == 0 && img.Bounds().Dy() > height {
			img = imaging.Resize(img, 0, height, imaging.Lanczos)
		} else if width != 0 && height != 0 {
			img = imaging.Resize(img, width, height, imaging.Lanczos)
		}
	}

	return jpeg.Encode(w, img, &jpeg.Options{Quality: quality})
}
