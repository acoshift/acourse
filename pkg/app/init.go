package app

import (
	"errors"
	"net/http"
	"strings"

	"github.com/acoshift/go-firebase-admin"
	"github.com/acoshift/httperror"
)

// ContextKey is the key for app's context
type ContextKey int

// Predefined context keys
const (
	KeyCurrentUserID ContextKey = iota
	KeyRequestID
)

var (
	firAuth *admin.Auth

	tokenError = httperror.New(http.StatusUnauthorized, "token")
)

// InitService inits service
func InitService(auth *admin.Auth) {
	firAuth = auth
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
