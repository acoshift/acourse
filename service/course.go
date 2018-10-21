package service

import (
	"bytes"
	"context"
	"io"

	"github.com/moonrhythm/dispatcher"

	"github.com/acoshift/acourse/context/appctx"
	"github.com/acoshift/acourse/context/sqlctx"
	"github.com/acoshift/acourse/model/app"
	"github.com/acoshift/acourse/model/course"
	"github.com/acoshift/acourse/model/file"
	"github.com/acoshift/acourse/model/image"
)

func (s *svc) createCourse(ctx context.Context, m *course.Create) error {
	// TODO: validate user role
	user := appctx.GetUser(ctx)

	if m.Title == "" {
		return app.NewUIError("title required")
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

		imageURL, err = s.uploadCourseCoverImage(ctx, image)
		image.Close()
		if err != nil {
			return app.NewUIError(err.Error())
		}
	}

	return sqlctx.RunInTx(ctx, func(ctx context.Context) error {
		var err error

		m.Result, err = registerCourse(ctx, &RegisterCourse{
			UserID:    user.ID,
			Title:     m.Title,
			ShortDesc: m.ShortDesc,
			LongDesc:  m.LongDesc,
			Image:     imageURL,
			Start:     m.Start,
		})
		if err != nil {
			return err
		}

		return dispatcher.Dispatch(ctx, &course.SetOption{ID: m.Result, Option: course.Option{}})
	})
}

func (s *svc) updateCourse(ctx context.Context, m *course.Update) error {
	// TODO: validate user role
	// user := appctx.GetUser(ctx)

	if m.ID == "" {
		return app.NewUIError("course id required")
	}

	if m.Title == "" {
		return app.NewUIError("title required")
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

		imageURL, err = s.uploadCourseCoverImage(ctx, image)
		image.Close()
		if err != nil {
			return app.NewUIError(err.Error())
		}
	}

	err := sqlctx.RunInTx(ctx, func(ctx context.Context) error {
		err := updateCourse(ctx, &UpdateCourseModel{
			ID:        m.ID,
			Title:     m.Title,
			ShortDesc: m.ShortDesc,
			LongDesc:  m.LongDesc,
			Start:     m.Start,
		})
		if err != nil {
			return err
		}

		if imageURL != "" {
			err = dispatcher.Dispatch(ctx, &course.SetImage{ID: m.ID, Image: imageURL})
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

	if err := dispatcher.Dispatch(ctx, &image.JPEG{
		Writer:  buf,
		Reader:  r,
		Width:   1200,
		Quality: 90,
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
