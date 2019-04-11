package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"io"
	"unicode/utf8"

	admin "github.com/acoshift/go-firebase-admin"
	"github.com/asaskevich/govalidator"

	"github.com/acoshift/acourse/internal/pkg/bus"
	"github.com/acoshift/acourse/internal/pkg/model/app"
	"github.com/acoshift/acourse/internal/pkg/model/auth"
	"github.com/acoshift/acourse/internal/pkg/model/user"
)

// Init inits auth service
func Init() {
	s := svc{}

	bus.Register(s.signUp)
	bus.Register(s.sendPasswordResetEmail)
	bus.Register(s.signInPassword)
	bus.Register(s.generateOpenIDURI)
	bus.Register(s.signInOpenIDCallback)
}

type svc struct{}

func (s *svc) signUp(ctx context.Context, m *auth.SignUp) error {
	if m.Email == "" {
		return app.NewUIError("email required")
	}

	email, err := govalidator.NormalizeEmail(m.Email)
	if err != nil {
		return app.NewUIError("invalid email")
	}

	if m.Password == "" {
		return app.NewUIError("password required")
	}
	if n := utf8.RuneCountInString(m.Password); n < 6 || n > 64 {
		return app.NewUIError("password must have 6 to 64 characters")
	}

	userID, err := firAuth.CreateUser(ctx, &admin.User{
		Email:    email,
		Password: m.Password,
	})
	if err != nil {
		return app.NewUIError(err.Error())
	}

	err = bus.Dispatch(ctx, &user.Create{
		ID:       userID,
		Username: userID,
		Email:    email,
	})
	if err == user.ErrEmailNotAvailable {
		return app.NewUIError("อีเมลนี้ถูกสมัครแล้ว")
	}
	if err == user.ErrUsernameNotAvailable {
		return app.NewUIError("username นี้ถูกใช้งานแล้ว")
	}
	if err != nil {
		return err
	}

	m.Result = userID

	return nil
}

func (s *svc) sendPasswordResetEmail(ctx context.Context, m *auth.SendPasswordResetEmail) error {
	if m.Email == "" {
		return app.NewUIError("email required")
	}

	u, err := firAuth.GetUserByEmail(ctx, m.Email)
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

func (s *svc) signInPassword(ctx context.Context, m *auth.SignInPassword) error {
	if m.Email == "" {
		return app.NewUIError("email required")
	}
	if m.Password == "" {
		return app.NewUIError("password required")
	}

	userID, err := firAuth.VerifyPassword(ctx, m.Email, m.Password)
	if err != nil {
		return app.NewUIError(err.Error())
	}

	// if user not found in our database, insert new user
	// this happened when database out of sync with firebase authentication
	{
		exists := user.IsExists{ID: userID}
		err = bus.Dispatch(ctx, &exists)
		if err != nil {
			return err
		}

		if !exists.Result {
			email, _ := govalidator.NormalizeEmail(m.Email)
			err = bus.Dispatch(ctx, &user.Create{
				ID:       userID,
				Username: userID,
				Name:     userID,
				Email:    email,
			})
			if err != nil {
				return err
			}
		}
	}

	m.Result = userID

	return nil
}

func generateRandomString(n int) string {
	b := make([]byte, n)
	io.ReadFull(rand.Reader, b)
	return base64.RawURLEncoding.EncodeToString(b)
}
