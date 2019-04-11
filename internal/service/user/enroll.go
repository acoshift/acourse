package user

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/acoshift/acourse/internal/pkg/bus"
	"github.com/acoshift/acourse/internal/pkg/context/appctx"
	"github.com/acoshift/acourse/internal/pkg/context/sqlctx"
	"github.com/acoshift/acourse/internal/pkg/course"
	"github.com/acoshift/acourse/internal/pkg/file"
	"github.com/acoshift/acourse/internal/pkg/image"
	"github.com/acoshift/acourse/internal/pkg/model"
	"github.com/acoshift/acourse/internal/pkg/model/app"
	"github.com/acoshift/acourse/internal/pkg/model/user"
	"github.com/acoshift/acourse/internal/pkg/notify"
	"github.com/acoshift/acourse/internal/pkg/payment"
)

func enroll(ctx context.Context, m *user.Enroll) error {
	u := appctx.GetUser(ctx)

	c, err := course.Get(ctx, m.CourseID)
	if err == model.ErrNotFound {
		return model.ErrNotFound
	}
	if err != nil {
		return err
	}

	// is owner
	if u.ID == c.UserID {
		return nil
	}

	// is enrolled
	{
		q := user.IsEnroll{ID: u.ID, CourseID: m.CourseID}
		err = bus.Dispatch(ctx, &q)
		if err != nil {
			return err
		}
		if q.Result {
			return nil
		}
	}

	// has pending enroll
	{
		hasPending, err := payment.HasPending(ctx, u.ID, m.CourseID)
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

	if m.Price < 0 {
		return app.NewUIError("จำนวนเงินติดลบไม่ได้")
	}

	var imageURL string
	if originalPrice != 0 {
		if m.PaymentImage == nil {
			return app.NewUIError("กรุณาอัพโหลดรูปภาพ")
		}

		err := image.Validate(m.PaymentImage)
		if err != nil {
			return err
		}

		img, err := m.PaymentImage.Open()
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
			return course.InsertEnroll(ctx, c.ID, u.ID)
		}

		newPayment = true

		return registerPayment(ctx, &RegisterPayment{
			CourseID:      c.ID,
			UserID:        u.ID,
			Image:         imageURL,
			Price:         m.Price,
			OriginalPrice: originalPrice,
			Status:        payment.Pending,
		})
	})
	if err != nil {
		return err
	}

	if newPayment {
		go notify.Admin(fmt.Sprintf("New payment for course %s, price %.2f", c.Title, m.Price))
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

// RegisterPayment type
type RegisterPayment struct {
	UserID        string
	CourseID      string
	Image         string
	Price         float64
	OriginalPrice float64
	Code          string
	Status        int
}

func registerPayment(ctx context.Context, x *RegisterPayment) error {
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
