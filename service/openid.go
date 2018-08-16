package service

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"github.com/acoshift/acourse/context/sqlctx"
	"github.com/acoshift/acourse/entity"
	"github.com/acoshift/acourse/file"
	"github.com/acoshift/acourse/repository"
)

var allowProvider = map[string]bool{
	"google.com": true,
	"github.com": true,
}

func (s *svc) GenerateOpenIDURI(ctx context.Context, provider string) (string, string, error) {
	if !allowProvider[provider] {
		return "", "", newUIError("provider not allowed")
	}

	sessID := generateSessionID()

	redirect, err := s.Auth.CreateAuthURI(ctx, provider, s.BaseURL+s.OpenIDCallback, sessID)
	if err != nil {
		return "", "", err
	}

	return redirect, sessID, nil
}

func (s *svc) SignInOpenIDCallback(ctx context.Context, uri string, state string) (string, error) {
	user, err := s.Auth.VerifyAuthCallbackURI(ctx, s.BaseURL+uri, state)
	if err != nil {
		return "", newUIError(err.Error())
	}

	err = sqlctx.RunInTx(ctx, func(ctx context.Context) error {
		// check is user sign up
		exists, err := repository.IsUserExists(ctx, user.UserID)
		if err != nil {
			return err
		}
		if !exists {
			// user not found, insert new user
			imageURL := s.uploadProfileFromURLAsync(user.PhotoURL)
			err = repository.RegisterUser(ctx, &entity.RegisterUser{
				ID:       user.UserID,
				Name:     user.DisplayName,
				Username: user.UserID,
				Email:    user.Email,
				Image:    imageURL,
			})
			if err != nil {
				return err
			}
		}

		return nil
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

	return user.UserID, nil
}

func generateSessionID() string {
	return generateRandomString(24)
}

// uploadProfileFromURLAsync copies data from given url and upload profile in background,
// returns url of destination file
func (s *svc) uploadProfileFromURLAsync(url string) string {
	if len(url) == 0 {
		return ""
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return ""
	}
	filename := file.GenerateFilename() + ".jpg"
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		req = req.WithContext(ctx)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return
		}
		defer resp.Body.Close()
		buf := &bytes.Buffer{}
		err = s.ImageResizeEncoder.ResizeEncode(buf, resp.Body, 500, 500, 90, true)
		if err != nil {
			return
		}
		cancel()
		ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err = s.FileStorage.Store(ctx, buf, filename)
		if err != nil {
			return
		}
	}()
	return s.FileStorage.DownloadURL(filename)
}
