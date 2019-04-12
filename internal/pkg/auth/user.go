package auth

import (
	"context"

	"github.com/acoshift/acourse/internal/pkg/user"
)

var userSvc interface {
	Create(ctx context.Context, m *user.CreateArgs) error
	IsExists(ctx context.Context, id string) (bool, error)
} = _userSvcImpl{}

type _userSvcImpl struct{}

func (_userSvcImpl) Create(ctx context.Context, m *user.CreateArgs) error {
	return user.Create(ctx, m)
}

func (_userSvcImpl) IsExists(ctx context.Context, id string) (bool, error) {
	return user.IsExists(ctx, id)
}
