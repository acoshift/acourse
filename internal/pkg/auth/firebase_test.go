package auth

import (
	"context"
	"fmt"

	"github.com/acoshift/go-firebase-admin"
)

type fakeFirebaseAuth struct {
	error error
}

var _ = func() bool {
	isInTest = true
	return true
}()

func init() {
	firAuth = &fakeFirebaseAuth{}
}

func (auth *fakeFirebaseAuth) CreateAuthURI(ctx context.Context, providerID, continueURI, sessionID string) (string, error) {
	return "http://localhost:9000", nil
}

func (auth *fakeFirebaseAuth) VerifyAuthCallbackURI(ctx context.Context, callbackURI, sessionID string) (*firebase.UserInfo, error) {
	return &firebase.UserInfo{}, nil
}

func (auth *fakeFirebaseAuth) GetUserByEmail(ctx context.Context, email string) (*firebase.UserRecord, error) {
	if email == "notfound@test.com" {
		return nil, fmt.Errorf("not found")
	}

	return &firebase.UserRecord{
		UserID: "123",
		Email:  email,
	}, nil
}

func (auth *fakeFirebaseAuth) SendPasswordResetEmail(ctx context.Context, email string) error {
	if email == "notfound@test.com" {
		return fmt.Errorf("not found")
	}
	return nil
}

func (auth *fakeFirebaseAuth) VerifyPassword(ctx context.Context, email, password string) (string, error) {
	if password == "fakepass" {
		return "", fmt.Errorf("invalid password")
	}
	return "123", nil
}

func (auth *fakeFirebaseAuth) CreateUser(ctx context.Context, user *firebase.User) (string, error) {
	return "123", nil
}
