package service

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/moonrhythm/dispatcher"

	"github.com/acoshift/acourse/context/appctx"
	"github.com/acoshift/acourse/context/sqlctx"
	"github.com/acoshift/acourse/entity"
	"github.com/acoshift/acourse/model/app"
	"github.com/acoshift/acourse/model/course"
	"github.com/acoshift/acourse/model/file"
	"github.com/acoshift/acourse/model/image"
	"github.com/acoshift/acourse/model/notify"
	"github.com/acoshift/acourse/model/payment"
	"github.com/acoshift/acourse/model/user"
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

		m.Result, err = s.Repository.RegisterCourse(ctx, &RegisterCourse{
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
		err := s.Repository.UpdateCourse(ctx, &UpdateCourseModel{
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

func (s *svc) enrollCourse(ctx context.Context, m *course.Enroll) error {
	u := appctx.GetUser(ctx)

	course, err := s.Repository.GetCourse(ctx, m.ID)
	if err == entity.ErrNotFound {
		return entity.ErrNotFound
	}
	if err != nil {
		return err
	}

	// is owner
	if u.ID == course.UserID {
		return nil
	}

	// is enrolled
	{
		q := user.IsEnroll{ID: u.ID, CourseID: m.ID}
		err = dispatcher.Dispatch(ctx, &q)
		if err != nil {
			return err
		}
		if q.Result {
			return nil
		}
	}

	// has pending enroll
	{
		q := payment.HasPending{UserID: u.ID, CourseID: m.ID}
		err := dispatcher.Dispatch(ctx, &q)
		if err != nil {
			return err
		}
		if q.Result {
			return nil
		}
	}

	originalPrice := course.Price
	if course.Option.Discount {
		originalPrice = course.Discount
	}

	if m.Price < 0 {
		return app.NewUIError("จำนวนเงินติดลบไม่ได้")
	}

	var imageURL string
	if originalPrice != 0 {
		if m.PaymentImage == nil {
			return app.NewUIError("กรุณาอัพโหลดรูปภาพ")
		}

		err := ValidateImage(m.PaymentImage)
		if err != nil {
			return err
		}

		image, err := m.PaymentImage.Open()
		if err != nil {
			return app.NewUIError(err.Error())
		}
		defer image.Close()

		imageURL, err = s.uploadPaymentImage(ctx, image)
		image.Close()
		if err != nil {
			return app.NewUIError(err.Error())
		}
	}

	newPayment := false

	err = sqlctx.RunInTx(ctx, func(ctx context.Context) error {
		if course.Price == 0 {
			return s.Repository.RegisterEnroll(ctx, u.ID, course.ID)
		}

		newPayment = true

		return s.Repository.RegisterPayment(ctx, &RegisterPayment{
			CourseID:      course.ID,
			UserID:        u.ID,
			Image:         imageURL,
			Price:         m.Price,
			OriginalPrice: originalPrice,
			Status:        entity.Pending,
		})
	})
	if err != nil {
		return err
	}

	if newPayment {
		go dispatcher.Dispatch(ctx, &notify.Admin{Message: fmt.Sprintf("New payment for course %s, price %.2f", course.Title, m.Price)})
	}

	return nil
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

// UploadPaymentImage uploads payment image
func (s *svc) uploadPaymentImage(ctx context.Context, r io.Reader) (string, error) {
	buf := &bytes.Buffer{}

	if err := dispatcher.Dispatch(ctx, &image.JPEG{
		Writer:  buf,
		Reader:  r,
		Width:   700,
		Quality: 60,
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
