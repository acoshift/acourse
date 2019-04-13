package me

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/acoshift/pgsql/pgctx"

	"github.com/acoshift/acourse/internal/pkg/context/appctx"
	"github.com/acoshift/acourse/internal/pkg/course"
	"github.com/acoshift/acourse/internal/pkg/file"
	"github.com/acoshift/acourse/internal/pkg/image"
	"github.com/acoshift/acourse/internal/pkg/notify"
	"github.com/acoshift/acourse/internal/pkg/payment"
)

var (
	ErrImageRequired = errors.New("me: image required")
)

// Enroll enrolls a course
func Enroll(ctx context.Context, courseID string, price float64, paymentImage *multipart.FileHeader) error {
	userID := appctx.GetUserID(ctx)

	c, err := course.Get(ctx, courseID)
	if err != nil {
		return err
	}

	// is owner
	if userID == c.Owner.ID {
		return nil
	}

	// is enrolled
	{
		isEnrolled, err := course.IsEnroll(ctx, userID, courseID)
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
		return fmt.Errorf("invalid price")
	}

	var imageURL string
	if originalPrice != 0 {
		if paymentImage == nil {
			return ErrImageRequired
		}

		err := image.Validate(paymentImage)
		if err != nil {
			return err
		}

		img, err := paymentImage.Open()
		if err != nil {
			return err
		}
		defer img.Close()

		imageURL, err = uploadPaymentImage(ctx, img)
		img.Close()
		if err != nil {
			return err
		}
	}

	err = pgctx.RunInTx(ctx, func(ctx context.Context) error {
		pgctx.Committed(ctx, func(ctx context.Context) {
			go notify.Admin(fmt.Sprintf("New payment for course %s, price %.2f", c.Title, price))
		})

		if c.Price == 0 {
			return course.InsertEnroll(ctx, c.ID, userID)
		}

		// language=SQL
		_, err := pgctx.Exec(ctx, `
			insert into payments
				(user_id, course_id, image, price, original_price, code, status)
			values
				($1, $2, $3, $4, $5, $6, $7)
			returning id
		`, userID, c.ID, imageURL, price, originalPrice, "", payment.Pending)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
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
