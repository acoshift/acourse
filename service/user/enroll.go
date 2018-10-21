package user

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

func enroll(ctx context.Context, m *user.Enroll) error {
	u := appctx.GetUser(ctx)

	getCourse := course.Get{ID: m.ID}
	err := dispatcher.Dispatch(ctx, &getCourse)
	if err == entity.ErrNotFound {
		return entity.ErrNotFound
	}
	if err != nil {
		return err
	}
	c := getCourse.Result

	// is owner
	if u.ID == c.UserID {
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

		image, err := m.PaymentImage.Open()
		if err != nil {
			return app.NewUIError(err.Error())
		}
		defer image.Close()

		imageURL, err = uploadPaymentImage(ctx, image)
		image.Close()
		if err != nil {
			return app.NewUIError(err.Error())
		}
	}

	newPayment := false

	err = sqlctx.RunInTx(ctx, func(ctx context.Context) error {
		if c.Price == 0 {
			return dispatcher.Dispatch(ctx, &course.InsertEnroll{ID: c.ID, UserID: u.ID})
		}

		newPayment = true

		return registerPayment(ctx, &RegisterPayment{
			CourseID:      c.ID,
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
		go dispatcher.Dispatch(ctx, &notify.Admin{Message: fmt.Sprintf("New payment for course %s, price %.2f", c.Title, m.Price)})
	}

	return nil
}

func uploadPaymentImage(ctx context.Context, r io.Reader) (string, error) {
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
