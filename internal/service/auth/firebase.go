package auth

import (
	"context"

	"github.com/acoshift/go-firebase-admin"
)

type FirebaseAuth interface {
	CreateAuthURI(ctx context.Context, providerID, continueURI, sessionID string) (string, error)
	VerifyAuthCallbackURI(ctx context.Context, callbackURI, sessionID string) (*firebase.UserInfo, error)
	GetUserByEmail(ctx context.Context, email string) (*firebase.UserRecord, error)
	SendPasswordResetEmail(ctx context.Context, email string) error
	VerifyPassword(ctx context.Context, email, password string) (string, error)
	CreateUser(ctx context.Context, user *firebase.User) (string, error)
}

func SetFirebaseAuth(client FirebaseAuth) {
	firAuth = client
}

var (
	firAuth FirebaseAuth
)
