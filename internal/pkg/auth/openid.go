package auth

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"github.com/acoshift/pgsql/pgctx"
	"github.com/moonrhythm/hime"

	"github.com/acoshift/acourse/internal/pkg/file"
	"github.com/acoshift/acourse/internal/pkg/image"
	"github.com/acoshift/acourse/internal/pkg/user"
)

var allowProvider = map[string]bool{
	"google.com": true,
	"github.com": true,
}

func GenerateOpenIDURI(ctx context.Context, provider string) (redirectURI, state string, err error) {
	if !allowProvider[provider] {
		return "", "", ErrInvalidProvider
	}

	state = generateSessionID()

	redirectURI, err = firAuth.CreateAuthURI(ctx,
		provider,
		hime.Global(ctx, "baseURL").(string)+hime.Route(ctx, "auth.openid.callback"),
		state,
	)
	return
}

func SignInOpenIDCallback(ctx context.Context, uri, state string) (string, error) {
	u, err := firAuth.VerifyAuthCallbackURI(ctx,
		hime.Global(ctx, "baseURL").(string)+uri,
		state,
	)
	if err != nil {
		return "", ErrInvalidCallbackURI
	}

	err = pgctx.RunInTx(ctx, func(ctx context.Context) error {
		// check is user sign up
		exists, err := userSvc.IsExists(ctx, u.UserID)
		if err != nil {
			return err
		}

		if !exists {
			// user not found, insert new user
			imageURL := uploadProfileFromURLAsync(u.PhotoURL)
			err = userSvc.Create(ctx, &user.CreateArgs{
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
		return "", ErrEmailNotAvailable
	}
	if err == user.ErrUsernameNotAvailable {
		return "", ErrUsernameNotAvailable
	}
	if err != nil {
		return "", err
	}

	return u.UserID, nil
}

func generateSessionID() string {
	return generateRandomString(24)
}

// uploadProfileFromURLAsync copies data from given url and upload profile in background,
// returns url of destination file
func uploadProfileFromURLAsync(url string) string {
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
