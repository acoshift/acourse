package image

import (
	"context"
	"image"
	"image/jpeg"

	"github.com/disintegration/imaging"

	model "github.com/acoshift/acourse/internal/model/image"
	"github.com/acoshift/acourse/internal/pkg/dispatcher"
)

// Init inits image service
func Init() {
	dispatcher.Register(encodeJPEG)
}

func encodeJPEG(_ context.Context, m *model.JPEG) error {
	img, _, err := image.Decode(m.Reader)
	if err != nil {
		return err
	}

	if m.Crop {
		img = imaging.Thumbnail(img, m.Width, m.Height, imaging.Lanczos)
	} else {
		if m.Height == 0 && img.Bounds().Dx() > m.Width {
			img = imaging.Resize(img, m.Width, 0, imaging.Lanczos)
		} else if m.Width == 0 && img.Bounds().Dy() > m.Height {
			img = imaging.Resize(img, 0, m.Height, imaging.Lanczos)
		} else if m.Width != 0 && m.Height != 0 {
			img = imaging.Resize(img, m.Width, m.Height, imaging.Lanczos)
		}
	}

	return jpeg.Encode(m.Writer, img, &jpeg.Options{Quality: m.Quality})
}
