package service

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"unicode/utf8"

	"github.com/asaskevich/govalidator"
	"github.com/moonrhythm/dispatcher"

	"github.com/acoshift/acourse/context/appctx"
	"github.com/acoshift/acourse/context/sqlctx"
	"github.com/acoshift/acourse/file"
)

// Profile type
type Profile struct {
	Username string
	Name     string
	AboutMe  string
	Image    *multipart.FileHeader
}

func (s *svc) UpdateProfile(ctx context.Context, x *Profile) error {
	user := appctx.GetUser(ctx)

	if !govalidator.IsAlphanumeric(x.Username) {
		return newUIError("username allow only a-z, A-Z, and 0-9")
	}
	if n := utf8.RuneCountInString(x.Username); n < 4 || n > 32 {
		return newUIError("username must have 4 - 32 characters")
	}
	if n := utf8.RuneCountInString(x.Name); n < 4 || n > 40 {
		return newUIError("name must have 4 - 40 characters")
	}
	if n := utf8.RuneCountInString(x.AboutMe); n > 256 {
		return newUIError("about me must have lower than 256 characters")
	}

	var imageURL string
	if x.Image != nil {
		err := ValidateImage(x.Image)
		if err != nil {
			return err
		}

		image, err := x.Image.Open()
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
			err := s.Repository.SetUserImage(ctx, user.ID, imageURL)
			if err != nil {
				return err
			}
		}

		return s.Repository.UpdateUser(ctx, &UpdateUser{
			ID:       user.ID,
			Username: x.Username,
			Name:     x.Name,
			AboutMe:  x.AboutMe,
		})
	})

	return err
}

// uploadProfileImage uploads profile image and return url
func (s *svc) uploadProfileImage(ctx context.Context, r io.Reader) (string, error) {
	buf := &bytes.Buffer{}
	err := s.ImageResizeEncoder.ResizeEncode(buf, r, 500, 500, 90, true)
	if err != nil {
		return "", err
	}
	filename := file.GenerateFilename() + ".jpg"
	store := file.Store{Reader: buf, Filename: filename}
	if err = dispatcher.Dispatch(ctx, &store); err != nil {
		return "", err
	}
	return store.Result, nil
}
