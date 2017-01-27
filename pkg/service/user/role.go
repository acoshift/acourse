package user

import (
	"context"

	"github.com/acoshift/acourse/pkg/acourse"
)

func (s *service) GetRole(ctx context.Context, req *acourse.UserIDRequest) (*acourse.Role, error) {
	if req.GetUserId() == "" {
		return &acourse.Role{}, nil
	}

	var x role
	s.client.GetByName(ctx, kindRole, req.GetUserId(), &x)

	return &acourse.Role{
		Admin:      x.Admin,
		Instructor: x.Instructor,
	}, nil
}
