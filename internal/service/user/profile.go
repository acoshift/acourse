package user

import (
	"bytes"
	"context"
	"io"
	"unicode/utf8"

	"github.com/asaskevich/govalidator"

	"github.com/acoshift/acourse/internal/pkg/bus"
	"github.com/acoshift/acourse/internal/pkg/context/sqlctx"
	"github.com/acoshift/acourse/internal/pkg/file"
	"github.com/acoshift/acourse/internal/pkg/image"
	"github.com/acoshift/acourse/internal/pkg/model/app"
	"github.com/acoshift/acourse/internal/pkg/model/user"
)

func updateProfile(ctx context.Context, m *user.UpdateProfile) error {
	if !govalidator.IsAlphanumeric(m.Username) {
		return app.NewUIError("username allow only a-z, A-Z, and 0-9")
	}
	if n := utf8.RuneCountInString(m.Username); n < 4 || n > 32 {
		return app.NewUIError("username must have 4 - 32 characters")
	}
	if n := utf8.RuneCountInString(m.Name); n < 4 || n > 40 {
		return app.NewUIError("name must have 4 - 40 characters")
	}
	if n := utf8.RuneCountInString(m.AboutMe); n > 256 {
		return app.NewUIError("about me must have lower than 256 characters")
	}

	var imageURL string
	if m.Image != nil {
		err := image.Validate(m.Image)
		if err != nil {
			return err
		}

		img, err := m.Image.Open()
		if err != nil {
			return err
		}
		defer img.Close()

		imageURL, err = uploadProfileImage(ctx, img)
		if err != nil {
			return app.NewUIError(err.Error())
		}
	}

	err := sqlctx.RunInTx(ctx, func(ctx context.Context) error {
		if imageURL != "" {
			err := bus.Dispatch(ctx, &user.SetImage{ID: m.ID, Image: imageURL})
			if err != nil {
				return err
			}
		}

		return bus.Dispatch(ctx, &user.Update{
			ID:       m.ID,
			Username: m.Username,
			Name:     m.Name,
			AboutMe:  m.AboutMe,
		})
	})

	return err
}

// uploadProfileImage uploads profile image and return url
func uploadProfileImage(ctx context.Context, r io.Reader) (string, error) {
	buf := &bytes.Buffer{}
	err := image.JPEG(buf, r, 500, 500, 90, true)
	if err != nil {
		return "", err
	}

	filename := file.GenerateFilename() + ".jpg"
	downloadURL, err := file.Store(ctx, buf, filename, false)
	if err != nil {
		return "", err
	}
	return downloadURL, nil
}
