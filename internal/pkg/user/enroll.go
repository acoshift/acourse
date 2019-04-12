package user

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/acoshift/acourse/internal/pkg/app"
	"github.com/acoshift/acourse/internal/pkg/context/sqlctx"
	"github.com/acoshift/acourse/internal/pkg/course"
	"github.com/acoshift/acourse/internal/pkg/file"
	"github.com/acoshift/acourse/internal/pkg/image"
	"github.com/acoshift/acourse/internal/pkg/notify"
	"github.com/acoshift/acourse/internal/pkg/payment"
)

// IsEnroll checks is user enrolled a course
func IsEnroll(ctx context.Context, userID, courseID string) (bool, error) {
	var b bool
	err := sqlctx.QueryRow(ctx, `
		select exists (
			select 1
			from enrolls
			where user_id = $1 and course_id = $2
		)
	`, userID, courseID).Scan(&b)
	return b, err
}

// TODO: move enroll to course

// Enroll enrolls a course
func Enroll(ctx context.Context, userID, courseID string, price float64, paymentImage *multipart.FileHeader) error {
	c, err := course.Get(ctx, courseID)
	if err == app.ErrNotFound {
		return app.ErrNotFound
	}
	if err != nil {
		return err
	}

	// is owner
	if userID == c.UserID {
		return nil
	}

	// is enrolled
	{
		isEnrolled, err := IsEnroll(ctx, userID, courseID)
		if err != nil {
			return err
		}
		if isEnrolled {
			return nil
		}
	}

	// has pending enroll
	{
		hasPending, err := payment.HasPending(ctx, userID, courseID)
		if err != nil {
			return err
		}
		if hasPending {
			return nil
		}
	}

	originalPrice := c.Price
	if c.Option.Discount {
		originalPrice = c.Discount
	}

	if price < 0 {
		return app.NewUIError("จำนวนเงินติดลบไม่ได้")
	}

	var imageURL string
	if originalPrice != 0 {
		if paymentImage == nil {
			return app.NewUIError("กรุณาอัพโหลดรูปภาพ")
		}

		err := image.Validate(paymentImage)
		if err != nil {
			return err
		}

		img, err := paymentImage.Open()
		if err != nil {
			return app.NewUIError(err.Error())
		}
		defer img.Close()

		imageURL, err = uploadPaymentImage(ctx, img)
		img.Close()
		if err != nil {
			return app.NewUIError(err.Error())
		}
	}

	newPayment := false

	err = sqlctx.RunInTx(ctx, func(ctx context.Context) error {
		if c.Price == 0 {
			return course.InsertEnroll(ctx, c.ID, userID)
		}

		newPayment = true

		return registerPayment(ctx, &registerPaymentArgs{
			CourseID:      c.ID,
			UserID:        userID,
			Image:         imageURL,
			Price:         price,
			OriginalPrice: originalPrice,
			Status:        payment.Pending,
		})
	})
	if err != nil {
		return err
	}

	if newPayment {
		go notify.Admin(fmt.Sprintf("New payment for course %s, price %.2f", c.Title, price))
	}

	return nil
}

func uploadPaymentImage(ctx context.Context, r io.Reader) (string, error) {
	buf := &bytes.Buffer{}

	err := image.JPEG(buf, r, 700, 0, 60, false)
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

type registerPaymentArgs struct {
	UserID        string
	CourseID      string
	Image         string
	Price         float64
	OriginalPrice float64
	Code          string
	Status        int
}

func registerPayment(ctx context.Context, x *registerPaymentArgs) error {
	_, err := sqlctx.Exec(ctx, `
		insert into payments
			(user_id, course_id, image, price, original_price, code, status)
		values
			($1, $2, $3, $4, $5, $6, $7)
		returning id
	`, x.UserID, x.CourseID, x.Image, x.Price, x.OriginalPrice, x.Code, x.Status)
	if err != nil {
		return err
	}
	return nil
}
