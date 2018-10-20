package service

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"github.com/moonrhythm/dispatcher"

	"github.com/acoshift/acourse/context/sqlctx"
	"github.com/acoshift/acourse/entity"
	"github.com/acoshift/acourse/model/file"
	"github.com/acoshift/acourse/model/firebase"
	"github.com/acoshift/acourse/model/image"
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

	authURI := firebase.CreateAuthURI{
		ProviderID:  provider,
		ContinueURI: s.BaseURL + s.OpenIDCallback,
		SessionID:   sessID,
	}
	err := dispatcher.Dispatch(ctx, &authURI)
	if err != nil {
		return "", "", err
	}

	return authURI.Result, sessID, nil
}

func (s *svc) SignInOpenIDCallback(ctx context.Context, uri string, state string) (string, error) {
	q := firebase.VerifyAuthCallbackURI{
		CallbackURI: s.BaseURL + uri,
		SessionID:   state,
	}
	err := dispatcher.Dispatch(ctx, &q)
	if err != nil {
		return "", newUIError(err.Error())
	}
	user := q.Result

	err = sqlctx.RunInTx(ctx, func(ctx context.Context) error {
		// check is user sign up
		exists, err := s.Repository.IsUserExists(ctx, user.UserID)
		if err != nil {
			return err
		}
		if !exists {
			// user not found, insert new user
			imageURL := s.uploadProfileFromURLAsync(user.PhotoURL)
			err = s.Repository.RegisterUser(ctx, &RegisterUser{
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

	// TODO: refactor
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return ""
	}
	filename := file.GenerateFilename() + ".jpg"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req = req.WithContext(ctx)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	buf := &bytes.Buffer{}
	if err := dispatcher.Dispatch(ctx, &image.JPEG{
		Writer:  buf,
		Reader:  resp.Body,
		Width:   500,
		Height:  500,
		Quality: 90,
		Crop:    true,
	}); err != nil {
		return ""
	}
	cancel()
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	store := file.Store{Reader: buf, Filename: filename, Async: true}
	err = dispatcher.Dispatch(ctx, &store)
	if err != nil {
		return ""
	}
	return store.Result
}
