package service

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"strings"
	"time"

	"github.com/acoshift/acourse/context/appctx"
	"github.com/acoshift/acourse/context/sqlctx"
	"github.com/acoshift/acourse/entity"
	"github.com/acoshift/acourse/file"
	"github.com/acoshift/acourse/repository"
	"github.com/acoshift/header"
)

// CreateCourse type
type CreateCourse struct {
	Title     string
	ShortDesc string
	LongDesc  string
	Image     *multipart.FileHeader
	Start     time.Time
}

func (s *svc) CreateCourse(ctx context.Context, x *CreateCourse) (courseID string, err error) {
	// TODO: validate user role
	user := appctx.GetUser(ctx)

	if x.Title == "" {
		return "", newUIError("title required")
	}

	var imageURL string
	if x.Image != nil {
		// TODO: allow only jpeg, png
		if !strings.Contains(x.Image.Header.Get(header.ContentType), "image") {
			return "", newUIError("file is not an image")
		}

		image, err := x.Image.Open()
		if err != nil {
			return "", err
		}
		defer image.Close()

		imageURL, err = s.uploadCourseCoverImage(ctx, image)
		if err != nil {
			return "", newUIError(err.Error())
		}
	}

	err = sqlctx.RunInTx(ctx, func(ctx context.Context) error {
		var err error

		courseID, err = repository.RegisterCourse(ctx, &entity.RegisterCourse{
			UserID:    user.ID,
			Title:     x.Title,
			ShortDesc: x.ShortDesc,
			LongDesc:  x.LongDesc,
			Image:     imageURL,
			Start:     x.Start,
		})
		if err != nil {
			return err
		}

		return repository.SetCourseOption(ctx, courseID, &entity.CourseOption{})
	})

	return
}

// UpdateCourse type
type UpdateCourse struct {
	ID        string
	Title     string
	ShortDesc string
	LongDesc  string
	Image     *multipart.FileHeader
	Start     time.Time
}

func (s *svc) UpdateCourse(ctx context.Context, x *UpdateCourse) error {
	// TODO: validate user role
	// user := appctx.GetUser(ctx)

	if x.ID == "" {
		return newUIError("course id required")
	}

	if x.Title == "" {
		return newUIError("title required")
	}

	var imageURL string
	if x.Image != nil {
		// TODO: allow only jpeg, png
		if !strings.Contains(x.Image.Header.Get(header.ContentType), "image") {
			return newUIError("file is not an image")
		}

		image, err := x.Image.Open()
		if err != nil {
			return err
		}
		defer image.Close()

		imageURL, err = s.uploadCourseCoverImage(ctx, image)
		if err != nil {
			return newUIError(err.Error())
		}
	}

	err := sqlctx.RunInTx(ctx, func(ctx context.Context) error {
		err := repository.UpdateCourse(ctx, &entity.UpdateCourse{
			ID:        x.ID,
			Title:     x.Title,
			ShortDesc: x.ShortDesc,
			LongDesc:  x.LongDesc,
			Start:     x.Start,
		})
		if err != nil {
			return err
		}

		if imageURL != "" {
			err = repository.SetCourseImage(ctx, x.ID, imageURL)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

// UploadCourseCoverImage uploads course cover image
func (s *svc) uploadCourseCoverImage(ctx context.Context, r io.Reader) (string, error) {
	buf := &bytes.Buffer{}
	err := s.ImageResizeEncoder.ResizeEncode(buf, r, 1200, 0, 90, false)
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
