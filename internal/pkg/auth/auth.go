package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"unicode/utf8"

	admin "github.com/acoshift/go-firebase-admin"
	"github.com/asaskevich/govalidator"
	"google.golang.org/api/googleapi"

	"github.com/acoshift/acourse/internal/pkg/user"
)

func SignUp(ctx context.Context, email, password string) (string, error) {
	if email == "" {
		return "", ErrEmailRequired
	}

	var err error
	email, err = govalidator.NormalizeEmail(email)
	if err != nil {
		return "", ErrEmailInvalid
	}

	if password == "" {
		return "", ErrPasswordRequired
	}
	if n := utf8.RuneCountInString(password); n < 6 || n > 64 {
		return "", ErrPasswordInvalid
	}

	userID, err := firAuth.CreateUser(ctx, &admin.User{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return "", err
	}

	err = userSvc.Create(ctx, &user.CreateArgs{
		ID:       userID,
		Username: userID,
		Email:    email,
	})
	if err == user.ErrEmailNotAvailable {
		return "", ErrEmailNotAvailable
	}
	if err == user.ErrUsernameNotAvailable {
		return "", ErrUsernameNotAvailable
	}
	if err != nil {
		return "", err
	}

	return userID, nil
}

func SendPasswordResetEmail(ctx context.Context, email string) error {
	if email == "" {
		return ErrEmailRequired
	}

	u, err := firAuth.GetUserByEmail(ctx, email)
	if err != nil {
		// don't send any error back to user
		return nil
	}

	err = firAuth.SendPasswordResetEmail(ctx, u.Email)
	if err != nil {
		return err
	}

	return nil
}

func SignInPassword(ctx context.Context, email, password string) (string, error) {
	if email == "" {
		return "", fmt.Errorf("email required")
	}
	if password == "" {
		return "", fmt.Errorf("password required")
	}

	userID, err := firAuth.VerifyPassword(ctx, email, password)
	var googErr googleapi.Error
	var retryNormalize bool
	if errors.As(err, &googErr) {
		switch googErr.Message {
		default:
			fallthrough
		case "INVALID_PASSWORD":
			return "", fmt.Errorf("invalid email or password")
		case "EMAIL_NOT_FOUND":
			retryNormalize = true
		}
	} else if err != nil {
		return "", err
	}
	if retryNormalize {
		normalizedEmail, err := govalidator.NormalizeEmail(email)
		if err != nil {
			return "", fmt.Errorf("invalid email or password")
		}
		userID, err = firAuth.VerifyPassword(ctx, normalizedEmail, password)
		if errors.As(err, &googErr) {
			switch googErr.Message {
			case "EMAIL_NOT_FOUND", "INVALID_PASSWORD":
				return "", fmt.Errorf("invalid email or password")
			default:
				return "", err
			}
		} else if err != nil {
			return "", err
		}
	}

	// if user not found in our database, insert new user
	// this happened when database out of sync with firebase authentication
	{
		exists, err := userSvc.IsExists(ctx, userID)
		if err != nil {
			return "", err
		}

		if !exists {
			email, _ := govalidator.NormalizeEmail(email)
			err = userSvc.Create(ctx, &user.CreateArgs{
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
