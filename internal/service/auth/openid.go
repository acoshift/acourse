package auth

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"github.com/moonrhythm/hime"

	"github.com/acoshift/acourse/internal/pkg/bus"
	"github.com/acoshift/acourse/internal/pkg/context/sqlctx"
	"github.com/acoshift/acourse/internal/pkg/file"
	"github.com/acoshift/acourse/internal/pkg/image"
	"github.com/acoshift/acourse/internal/pkg/model/app"
	"github.com/acoshift/acourse/internal/pkg/model/auth"
	"github.com/acoshift/acourse/internal/pkg/model/user"
)

var allowProvider = map[string]bool{
	"google.com": true,
	"github.com": true,
}

func (s *svc) generateOpenIDURI(ctx context.Context, m *auth.GenerateOpenIDURI) error {
	if !allowProvider[m.Provider] {
		return app.NewUIError("provider not allowed")
	}

	sessID := generateSessionID()

	redirectURI, err := firAuth.CreateAuthURI(ctx,
		m.Provider,
		hime.Global(ctx, "baseURL").(string)+hime.Route(ctx, "auth.openid.callback"),
		sessID,
	)
	if err != nil {
		return err
	}

	m.Result.RedirectURI = redirectURI
	m.Result.State = sessID

	return nil
}

func (s *svc) signInOpenIDCallback(ctx context.Context, m *auth.SignInOpenIDCallback) error {
	u, err := firAuth.VerifyAuthCallbackURI(ctx,
		hime.Global(ctx, "baseURL").(string)+m.URI,
		m.State,
	)
	if err != nil {
		return app.NewUIError(err.Error())
	}

	err = sqlctx.RunInTx(ctx, func(ctx context.Context) error {
		// check is user sign up
		exists := user.IsExists{ID: u.UserID}
		err = bus.Dispatch(ctx, &exists)
		if err != nil {
			return err
		}

		if !exists.Result {
			// user not found, insert new user
			imageURL := s.uploadProfileFromURLAsync(u.PhotoURL)
			err = bus.Dispatch(ctx, &user.Create{
				ID:       u.UserID,
				Name:     u.DisplayName,
				Username: u.UserID,
				Email:    u.Email,
				Image:    imageURL,
			})
			if err != nil {
				return err
			}
		}

		return nil
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

	m.Result = u.UserID

	return nil
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
	err = image.JPEG(buf, resp.Body, 500, 500, 90, true)
	if err != nil {
		return ""
	}
	cancel()
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	downloadURL, err := file.Store(ctx, buf, filename, true)
	if err != nil {
		return ""
	}
	return downloadURL
}
