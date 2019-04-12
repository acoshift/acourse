package user

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/acoshift/pgsql/pgctx"

	"github.com/acoshift/acourse/internal/pkg/app"
	"github.com/acoshift/acourse/internal/pkg/course"
	"github.com/acoshift/acourse/internal/pkg/file"
	"github.com/acoshift/acourse/internal/pkg/image"
	"github.com/acoshift/acourse/internal/pkg/notify"
	"github.com/acoshift/acourse/internal/pkg/payment"
)

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

	err = pgctx.RunInTx(ctx, func(ctx context.Context) error {
		if c.Price == 0 {
			return InsertEnroll(ctx, c.ID, userID)
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

		pgctx.Committed(ctx, func(ctx context.Context) {
			go notify.Admin(fmt.Sprintf("New payment for course %s, price %.2f", c.Title, price))
		})

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

// InsertEnroll inserts enroll
func InsertEnroll(ctx context.Context, id, userID string) error {
	_, err := pgctx.Exec(ctx, `
		insert into enrolls
			(user_id, course_id)
		values
			($1, $2)
	`, userID, id)
	return err
}
