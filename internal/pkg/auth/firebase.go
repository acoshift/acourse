package auth

import (
	"context"

	"github.com/acoshift/go-firebase-admin"

	"github.com/acoshift/acourse/internal/pkg/config"
)

var firAuth interface {
	CreateAuthURI(ctx context.Context, providerID, continueURI, sessionID string) (string, error)
	VerifyAuthCallbackURI(ctx context.Context, callbackURI, sessionID string) (*firebase.UserInfo, error)
	GetUserByEmail(ctx context.Context, email string) (*firebase.UserRecord, error)
	SendPasswordResetEmail(ctx context.Context, email string) error
	VerifyPassword(ctx context.Context, email, password string) (string, error)
	CreateUser(ctx context.Context, user *firebase.User) (string, error)
}

var isInTest = false

func Init() {
	firAuth = config.FirebaseApp().Auth()
}
