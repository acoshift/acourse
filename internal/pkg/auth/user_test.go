package auth

import (
	"context"

	"github.com/acoshift/acourse/internal/pkg/user"
)

type fakeUserSvc struct{}

func init() {
	userSvc = fakeUserSvc{}
}

func (fakeUserSvc) Create(ctx context.Context, m *user.CreateArgs) error {
	if m.Email == "notavailable@test.com" {
		return user.ErrEmailNotAvailable
	}
	return nil
}

func (fakeUserSvc) IsExists(ctx context.Context, id string) (bool, error) {
	return true, nil
}
