package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"io"
	"unicode/utf8"

	admin "github.com/acoshift/go-firebase-admin"
	"github.com/asaskevich/govalidator"
	"github.com/moonrhythm/dispatcher"

	"github.com/acoshift/acourse/entity"
	"github.com/acoshift/acourse/model/firebase"
)

func (s *svc) SignUp(ctx context.Context, email, password string) (string, error) {
	if email == "" {
		return "", newUIError("email required")
	}

	email, err := govalidator.NormalizeEmail(email)
	if err != nil {
		return "", newUIError("invalid email")
	}

	if password == "" {
		return "", newUIError("password required")
	}
	if n := utf8.RuneCountInString(password); n < 6 || n > 64 {
		return "", newUIError("password must have 6 to 64 characters")
	}

	createUser := firebase.CreateUser{
		User: &admin.User{
			Email:    email,
			Password: password,
		},
	}
	err = dispatcher.Dispatch(ctx, &createUser)
	if err != nil {
		return "", newUIError(err.Error())
	}
	userID := createUser.Result

	err = s.Repository.RegisterUser(ctx, &RegisterUser{
		ID:       userID,
		Username: userID,
		Email:    email,
	})
	if err == entity.ErrEmailNotAvailable {
		return "", newUIError("อีเมลนี้ถูกสมัครแล้ว")
	}
	if err == entity.ErrUsernameNotAvailable {
		return "", newUIError("username นี้ถูกใช้งานแล้ว")
	}
	if err != nil {
		return "", err
	}

	return userID, nil
}

func (s *svc) SendPasswordResetEmail(ctx context.Context, email string) error {
	if email == "" {
		return newUIError("email required")
	}

	getUser := firebase.GetUserByEmail{Email: email}
	err := dispatcher.Dispatch(ctx, &getUser)
	if err != nil {
		// don't send any error back to user
		return nil
	}

	err = dispatcher.Dispatch(ctx, &firebase.SendPasswordResetEmail{Email: getUser.Email})
	if err != nil {
		return newUIError(err.Error())
	}

	return nil
}

func (s *svc) SignInPassword(ctx context.Context, email, password string) (string, error) {
	if email == "" {
		return "", newUIError("email required")
	}
	if password == "" {
		return "", newUIError("password required")
	}

	q := firebase.VerifyPassword{Email: email, Password: password}
	err := dispatcher.Dispatch(ctx, &q)
	if err != nil {
		return "", newUIError(err.Error())
	}
	userID := q.Result

	// if user not found in our database, insert new user
	// this happend when database out of sync with firebase authentication
	{
		ok, err := s.Repository.IsUserExists(ctx, userID)
		if err != nil {
			return "", err
		}

		if !ok {
			email, _ = govalidator.NormalizeEmail(email)
			err = s.Repository.RegisterUser(ctx, &RegisterUser{
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

	return userID, err
}

func generateRandomString(n int) string {
	b := make([]byte, n)
	io.ReadFull(rand.Reader, b)
	return base64.RawURLEncoding.EncodeToString(b)
}
