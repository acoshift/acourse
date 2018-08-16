package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/acoshift/acourse/context/appctx"
	"github.com/acoshift/acourse/context/sqlctx"
	"github.com/acoshift/acourse/entity"
	"github.com/acoshift/acourse/file"
)

func (s *svc) CreateCourse(ctx context.Context, x *CreateCourse) (courseID string, err error) {
	// TODO: validate user role
	user := appctx.GetUser(ctx)

	if x.Title == "" {
		return "", newUIError("title required")
	}

	var imageURL string
	if x.Image != nil {
		err := validateImage(x.Image)
		if err != nil {
			return "", err
		}

		image, err := x.Image.Open()
		if err != nil {
			return "", err
		}
		defer image.Close()

		imageURL, err = s.uploadCourseCoverImage(ctx, image)
		image.Close()
		if err != nil {
			return "", newUIError(err.Error())
		}
	}

	err = sqlctx.RunInTx(ctx, func(ctx context.Context) error {
		var err error

		courseID, err = s.Repository.RegisterCourse(ctx, &RegisterCourse{
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

		return s.Repository.SetCourseOption(ctx, courseID, &entity.CourseOption{})
	})

	return
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
		err := validateImage(x.Image)
		if err != nil {
			return err
		}

		image, err := x.Image.Open()
		if err != nil {
			return err
		}
		defer image.Close()

		imageURL, err = s.uploadCourseCoverImage(ctx, image)
		image.Close()
		if err != nil {
			return newUIError(err.Error())
		}
	}

	err := sqlctx.RunInTx(ctx, func(ctx context.Context) error {
		err := s.Repository.UpdateCourse(ctx, &UpdateCourseModel{
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
			err = s.Repository.SetCourseImage(ctx, x.ID, imageURL)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func (s *svc) EnrollCourse(ctx context.Context, courseID string, price float64, paymentImage *multipart.FileHeader) error {
	user := appctx.GetUser(ctx)

	course, err := s.Repository.GetCourse(ctx, courseID)
	if err == entity.ErrNotFound {
		return entity.ErrNotFound
	}
	if err != nil {
		return err
	}

	// is owner
	if user.ID == course.UserID {
		return nil
	}

	// is enrolled
	enrolled, err := s.Repository.IsEnrolled(ctx, user.ID, courseID)
	if err != nil {
		return err
	}
	if enrolled {
		return nil
	}

	// has pending enroll
	pendingPayment, err := s.Repository.HasPendingPayment(ctx, user.ID, courseID)
	if err != nil {
		return err
	}
	if pendingPayment {
		return nil
	}

	originalPrice := course.Price
	if course.Option.Discount {
		originalPrice = course.Discount
	}

	if price < 0 {
		return newUIError("จำนวนเงินติดลบไม่ได้")
	}

	var imageURL string
	if originalPrice != 0 {
		if paymentImage == nil {
			return newUIError("กรุณาอัพโหลดรูปภาพ")
		}

		err := validateImage(paymentImage)
		if err != nil {
			return err
		}

		image, err := paymentImage.Open()
		if err != nil {
			return newUIError(err.Error())
		}
		defer image.Close()

		imageURL, err = s.uploadPaymentImage(ctx, image)
		image.Close()
		if err != nil {
			return newUIError(err.Error())
		}
	}

	newPayment := false

	err = sqlctx.RunInTx(ctx, func(ctx context.Context) error {
		if course.Price == 0 {
			return s.Repository.RegisterEnroll(ctx, user.ID, course.ID)
		}

		newPayment = true

		return s.Repository.RegisterPayment(ctx, &RegisterPayment{
			CourseID:      course.ID,
			UserID:        user.ID,
			Image:         imageURL,
			Price:         price,
			OriginalPrice: originalPrice,
			Status:        entity.Pending,
		})
	})
	if err != nil {
		return err
	}

	if newPayment {
		go s.AdminNotifier.Notify(fmt.Sprintf("New payment for course %s, price %.2f", course.Title, price))
	}

	return nil
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

// UploadPaymentImage uploads payment image
func (s *svc) uploadPaymentImage(ctx context.Context, r io.Reader) (string, error) {
	buf := &bytes.Buffer{}
	err := s.ImageResizeEncoder.ResizeEncode(buf, r, 700, 0, 60, false)
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
