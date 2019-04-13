package me

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"unicode/utf8"

	"github.com/acoshift/pgsql/pgctx"
	"github.com/asaskevich/govalidator"
	"github.com/moonrhythm/validator"

	"github.com/acoshift/acourse/internal/pkg/context/appctx"
	"github.com/acoshift/acourse/internal/pkg/file"
	"github.com/acoshift/acourse/internal/pkg/image"
	"github.com/acoshift/acourse/internal/pkg/user"
)

type UpdateProfileArgs struct {
	Username string
	Name     string
	AboutMe  string
	Image    *multipart.FileHeader
}

func UpdateProfile(ctx context.Context, m *UpdateProfileArgs) error {
	userID := appctx.GetUserID(ctx)

	v := validator.New()
	v.Must(govalidator.IsAlphanumeric(m.Username), "username allow only a-z, A-Z, and 0-9")
	{
		n := utf8.RuneCountInString(m.Username)
		v.Must(n >= 4 && n <= 32, "username must have 4 - 32 characters")
	}
	{
		n := utf8.RuneCountInString(m.Name)
		v.Must(n >= 4 && n <= 40, "name must have 4 - 40 characters")
	}
	{
		n := utf8.RuneCountInString(m.AboutMe)
		v.Must(n <= 256, "about me must have lower than 256 characters")
	}
	if !v.Valid() {
		return v.Error()
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
			return err
		}
	}

	err := pgctx.RunInTx(ctx, func(ctx context.Context) error {
		if imageURL != "" {
			err := user.SetImage(ctx, userID, imageURL)
			if err != nil {
				return err
			}
		}

		return user.Update(ctx, &user.UpdateArgs{
			ID:       userID,
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
