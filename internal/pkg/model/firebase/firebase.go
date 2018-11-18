package firebase

import (
	"github.com/acoshift/go-firebase-admin"
)

// CreateAuthURI calls firebase.CreateAuthURI
type CreateAuthURI struct {
	ProviderID  string
	ContinueURI string
	SessionID   string

	Result string
}

// VerifyAuthCallbackURI calls firebase.VerifyAuthCallbackURI
type VerifyAuthCallbackURI struct {
	CallbackURI string
	SessionID   string

	Result *firebase.UserInfo
}

// GetUserByEmail calls firebase.GetUserByEmail
type GetUserByEmail struct {
	Email string

	Result *firebase.UserRecord
}

// SendPasswordResetEmail calls firebase.SendPasswordResetEmail
type SendPasswordResetEmail struct {
	Email string
}

// VerifyPassword calls firebase.VerifyPassword
type VerifyPassword struct {
	Email    string
	Password string

	Result string
}

// CreateUser calls firebase.CreateUser
type CreateUser struct {
	User *firebase.User

	Result string
}
