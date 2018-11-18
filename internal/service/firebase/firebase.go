package firebase

import (
	"context"

	admin "github.com/acoshift/go-firebase-admin"

	"github.com/acoshift/acourse/internal/model/firebase"
	"github.com/acoshift/acourse/internal/pkg/dispatcher"
)

// Init inits firebase
func Init(auth *admin.Auth) {
	s := svc{auth}

	dispatcher.Register(s.createAuthURI)
	dispatcher.Register(s.verifyAuthCallbackURI)
	dispatcher.Register(s.getUserByEmail)
	dispatcher.Register(s.sendPasswordResetEmail)
	dispatcher.Register(s.verifyPassword)
	dispatcher.Register(s.createUser)
}

type svc struct {
	auth *admin.Auth
}

func (s *svc) createAuthURI(ctx context.Context, m *firebase.CreateAuthURI) error {
	var err error
	m.Result, err = s.auth.CreateAuthURI(ctx, m.ProviderID, m.ContinueURI, m.SessionID)
	return err
}

func (s *svc) verifyAuthCallbackURI(ctx context.Context, m *firebase.VerifyAuthCallbackURI) error {
	var err error
	m.Result, err = s.auth.VerifyAuthCallbackURI(ctx, m.CallbackURI, m.SessionID)
	return err
}

func (s *svc) getUserByEmail(ctx context.Context, m *firebase.GetUserByEmail) error {
	var err error
	m.Result, err = s.auth.GetUserByEmail(ctx, m.Email)
	return err
}

func (s *svc) sendPasswordResetEmail(ctx context.Context, m *firebase.SendPasswordResetEmail) error {
	return s.auth.SendPasswordResetEmail(ctx, m.Email)
}

func (s *svc) verifyPassword(ctx context.Context, m *firebase.VerifyPassword) error {
	var err error
	m.Result, err = s.auth.VerifyPassword(ctx, m.Email, m.Password)
	return err
}

func (s *svc) createUser(ctx context.Context, m *firebase.CreateUser) error {
	var err error
	m.Result, err = s.auth.CreateUser(ctx, m.User)
	return err
}
