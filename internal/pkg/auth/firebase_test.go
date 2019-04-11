package auth_test

import (
	"context"

	"github.com/acoshift/go-firebase-admin"

	. "github.com/acoshift/acourse/internal/pkg/auth"
)

type fakeFirebaseAuth struct {
	error error
}

var firAuth = &fakeFirebaseAuth{}

func init() {
	SetFirebaseAuth(firAuth)
}

func (auth *fakeFirebaseAuth) CreateAuthURI(ctx context.Context, providerID, continueURI, sessionID string) (string, error) {
	if auth.error != nil {
		return "", auth.error
	}
	return "http://localhost:9000", nil
}

func (auth *fakeFirebaseAuth) VerifyAuthCallbackURI(ctx context.Context, callbackURI, sessionID string) (*firebase.UserInfo, error) {
	if auth.error != nil {
		return nil, auth.error
	}
	return &firebase.UserInfo{}, nil
}

func (auth *fakeFirebaseAuth) GetUserByEmail(ctx context.Context, email string) (*firebase.UserRecord, error) {
	if auth.error != nil {
		return nil, auth.error
	}
	return &firebase.UserRecord{
		UserID: "123",
		Email:  email,
	}, nil
}

func (auth *fakeFirebaseAuth) SendPasswordResetEmail(ctx context.Context, email string) error {
	if auth.error != nil {
		return auth.error
	}
	return nil
}

func (auth *fakeFirebaseAuth) VerifyPassword(ctx context.Context, email, password string) (string, error) {
	if auth.error != nil {
		return "", auth.error
	}
	return "123", nil
}

func (auth *fakeFirebaseAuth) CreateUser(ctx context.Context, user *firebase.User) (string, error) {
	if auth.error != nil {
		return "", auth.error
	}
	return "123", nil
}
