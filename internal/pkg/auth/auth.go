package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"io"
	"unicode/utf8"

	admin "github.com/acoshift/go-firebase-admin"
	"github.com/asaskevich/govalidator"

	"github.com/acoshift/acourse/internal/pkg/app"
	"github.com/acoshift/acourse/internal/pkg/bus"
	"github.com/acoshift/acourse/internal/pkg/model/user"
)

func SignUp(ctx context.Context, email, password string) (string, error) {
	if email == "" {
		return "", app.NewUIError("email required")
	}

	var err error
	email, err = govalidator.NormalizeEmail(email)
	if err != nil {
		return "", app.NewUIError("invalid email")
	}

	if password == "" {
		return "", app.NewUIError("password required")
	}
	if n := utf8.RuneCountInString(password); n < 6 || n > 64 {
		return "", app.NewUIError("password must have 6 to 64 characters")
	}

	userID, err := firAuth.CreateUser(ctx, &admin.User{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return "", app.NewUIError(err.Error())
	}

	err = bus.Dispatch(ctx, &user.Create{
		ID:       userID,
		Username: userID,
		Email:    email,
	})
	if err == user.ErrEmailNotAvailable {
		return "", app.NewUIError("อีเมลนี้ถูกสมัครแล้ว")
	}
	if err == user.ErrUsernameNotAvailable {
		return "", app.NewUIError("username นี้ถูกใช้งานแล้ว")
	}
	if err != nil {
		return "", err
	}

	return userID, nil
}

func SendPasswordResetEmail(ctx context.Context, email string) error {
	if email == "" {
		return app.NewUIError("email required")
	}

	u, err := firAuth.GetUserByEmail(ctx, email)
	if err != nil {
		// don't send any error back to user
		return nil
	}

	err = firAuth.SendPasswordResetEmail(ctx, u.Email)
	if err != nil {
		return app.NewUIError(err.Error())
	}

	return nil
}

func SignInPassword(ctx context.Context, email, password string) (string, error) {
	if email == "" {
		return "", app.NewUIError("email required")
	}
	if password == "" {
		return "", app.NewUIError("password required")
	}

	userID, err := firAuth.VerifyPassword(ctx, email, password)
	if err != nil {
		return "", app.NewUIError(err.Error())
	}

	// if user not found in our database, insert new user
	// this happened when database out of sync with firebase authentication
	{
		exists := user.IsExists{ID: userID}
		err = bus.Dispatch(ctx, &exists)
		if err != nil {
			return "", err
		}

		if !exists.Result {
			email, _ := govalidator.NormalizeEmail(email)
			err = bus.Dispatch(ctx, &user.Create{
				ID:       userID,
				Username: userID,
				Name:     userID,
				Email:    email,
			})
			if err != nil {
				return "", err
			}
		}
	}

	return userID, nil
}

func generateRandomString(n int) string {
	b := make([]byte, n)
	io.ReadFull(rand.Reader, b)
	return base64.RawURLEncoding.EncodeToString(b)
}
