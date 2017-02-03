package app

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/logging"
	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"github.com/acoshift/go-firebase-admin"
	"github.com/acoshift/httperror"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	firAuth *admin.Auth

	tokenError = httperror.New(http.StatusUnauthorized, "token")
)

// InitService inits service
func InitService(auth *admin.Auth) {
	firAuth = auth
}

// MakeTokenSource creates token source from service account
func MakeTokenSource(serviceAccount []byte) (oauth2.TokenSource, error) {
	jwtConfig, err := google.JWTConfigFromJSON([]byte(serviceAccount),
		datastore.ScopeDatastore,
		pubsub.ScopePubSub,
		storage.ScopeFullControl,
		logging.WriteScope,
		"https://www.googleapis.com/auth/trace.append",
		"https://www.googleapis.com/auth/cloud-platform",
	)

	if err != nil {
		return nil, err
	}
	return jwtConfig.TokenSource(context.Background()), nil
}

func validateHeaderToken(header string) (string, error) {
	tk := strings.Split(header, " ")
	if len(tk) != 2 || strings.ToLower(tk[0]) != "bearer" {
		return "", errors.New("invalid authorization header")
	}
	claims, err := firAuth.VerifyIDToken(tk[1])
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}
