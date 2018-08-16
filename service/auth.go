package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/url"
	"unicode/utf8"

	"github.com/acoshift/go-firebase-admin"
	"github.com/asaskevich/govalidator"

	"github.com/acoshift/acourse/entity"
	"github.com/acoshift/acourse/view"
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

	userID, err := s.Auth.CreateUser(ctx, &firebase.User{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return "", err
	}

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

	user, err := s.Auth.GetUserByEmail(ctx, email)
	if err != nil {
		// don't send any error back to user
		return nil
	}

	err = s.Auth.SendPasswordResetEmail(ctx, user.Email)
	if err != nil {
		return newUIError(err.Error())
	}

	return nil
}

func (s *svc) SendSignInMagicLinkEmail(ctx context.Context, email string) error {
	if email == "" {
		return newUIError("email required")
	}
	email, err := govalidator.NormalizeEmail(email)
	if err != nil {
		return newUIError("invalid email")
	}

	ok, err := s.Repository.CanAcquireMagicLink(ctx, email)
	if err != nil {
		return err
	}
	if !ok {
		return newUIError("อีเมลของคุณได้ขอ Magic Link จากเราไปแล้ว กรุณาตรวจสอบอีเมล")
	}

	user, err := s.Repository.GetEmailSignInUserByEmail(ctx, email)
	if err == entity.ErrNotFound {
		return nil
	}
	if err != nil {
		return err
	}

	linkID := generateMagicLinkID()

	err = s.Repository.StoreMagicLink(ctx, linkID, user.ID)
	if err != nil {
		return err
	}

	linkQuery := make(url.Values)
	linkQuery.Set("id", linkID)

	message := fmt.Sprintf(`สวัสดีครับคุณ %s,


ตามที่ท่านได้ขอ Magic Link เพื่อเข้าสู่ระบบสำหรับ acourse.io นั้นท่านสามารถเข้าได้ผ่าน Link ข้างล่างนี้ ภายใน 1 ชม.

%s

ทีมงาน acourse.io
	`, user.Name, s.BaseURL+s.MagicLinkCallback+"?"+linkQuery.Encode())

	go s.EmailSender.Send(user.Email, "Magic Link Request", view.MarkdownEmail(message))

	return nil
}

func (s *svc) SignInPassword(ctx context.Context, email, password string) (string, error) {
	if email == "" {
		return "", newUIError("email required")
	}
	if password == "" {
		return "", newUIError("password required")
	}

	userID, err := s.Auth.VerifyPassword(ctx, email, password)
	if err != nil {
		return "", newUIError(err.Error())
	}

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

func (s *svc) SignInMagicLink(ctx context.Context, link string) (string, error) {
	if link == "" {
		return "", newUIError("ไม่พบ Magic Link ของคุณ")
	}

	userID, err := s.Repository.FindMagicLink(ctx, link)
	if err == entity.ErrNotFound {
		return "", newUIError("ไม่พบ Magic Link ของคุณ")
	}
	return userID, err
}

func generateRandomString(n int) string {
	b := make([]byte, n)
	io.ReadFull(rand.Reader, b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func generateMagicLinkID() string {
	return generateRandomString(64)
}
