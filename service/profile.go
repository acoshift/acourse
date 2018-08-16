package service

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"

	"github.com/acoshift/acourse/context/appctx"
	"github.com/acoshift/acourse/context/sqlctx"
	"github.com/acoshift/acourse/entity"
	"github.com/acoshift/acourse/file"
	"github.com/acoshift/acourse/repository"
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

	// TODO: validate profile

	var imageURL string
	if x.Image != nil {
		err := validateImage(x.Image)
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
			err := repository.SetUserImage(ctx, user.ID, imageURL)
			if err != nil {
				return err
			}
		}

		return repository.UpdateUser(ctx, &entity.UpdateUser{
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
	downloadURL := s.FileStorage.DownloadURL(filename)
	err = s.FileStorage.Store(ctx, buf, filename)
	if err != nil {
		return "", err
	}
	return downloadURL, nil
}
