package service

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"unicode/utf8"

	"github.com/asaskevich/govalidator"
	"github.com/moonrhythm/dispatcher"

	"github.com/acoshift/acourse/context/sqlctx"
	"github.com/acoshift/acourse/model/file"
	"github.com/acoshift/acourse/model/image"
	"github.com/acoshift/acourse/model/user"
)

// Profile type
type Profile struct {
	Username string
	Name     string
	AboutMe  string
	Image    *multipart.FileHeader
}

func (s *svc) updateProfile(ctx context.Context, m *user.UpdateProfile) error {
	if !govalidator.IsAlphanumeric(m.Username) {
		return newUIError("username allow only a-z, A-Z, and 0-9")
	}
	if n := utf8.RuneCountInString(m.Username); n < 4 || n > 32 {
		return newUIError("username must have 4 - 32 characters")
	}
	if n := utf8.RuneCountInString(m.Name); n < 4 || n > 40 {
		return newUIError("name must have 4 - 40 characters")
	}
	if n := utf8.RuneCountInString(m.AboutMe); n > 256 {
		return newUIError("about me must have lower than 256 characters")
	}

	var imageURL string
	if m.Image != nil {
		err := ValidateImage(m.Image)
		if err != nil {
			return err
		}

		image, err := m.Image.Open()
		if err != nil {
			return err
		}
		defer image.Close()

		imageURL, err = s.uploadProfileImage(ctx, image)
		if err != nil {
			return newUIError(err.Error())
		}
	}

	err := sqlctx.RunInTx(ctx, func(ctx context.Context) error {
		if imageURL != "" {
			err := s.Repository.SetUserImage(ctx, m.ID, imageURL)
			if err != nil {
				return err
			}
		}

		return s.Repository.UpdateUser(ctx, &UpdateUser{
			ID:       m.ID,
			Username: m.Username,
			Name:     m.Name,
			AboutMe:  m.AboutMe,
		})
	})

	return err
}

// uploadProfileImage uploads profile image and return url
func (s *svc) uploadProfileImage(ctx context.Context, r io.Reader) (string, error) {
	buf := &bytes.Buffer{}

	if err := dispatcher.Dispatch(ctx, &image.JPEG{
		Writer:  buf,
		Reader:  r,
		Width:   500,
		Height:  500,
		Quality: 90,
		Crop:    true,
	}); err != nil {
		return "", err
	}

	filename := file.GenerateFilename() + ".jpg"
	store := file.Store{Reader: buf, Filename: filename}
	if err := dispatcher.Dispatch(ctx, &store); err != nil {
		return "", err
	}
	return store.Result, nil
}
