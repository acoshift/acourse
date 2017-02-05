package user

import (
	"github.com/acoshift/acourse/pkg/acourse"
	context "golang.org/x/net/context"
)

func (s *service) GetRole(ctx context.Context, req *acourse.UserIDRequest) (*acourse.Role, error) {
	if req.GetUserId() == "" {
		return &acourse.Role{}, nil
	}

	var x roleModel
	s.client.GetByName(ctx, kindRole, req.GetUserId(), &x)

	return &acourse.Role{
		Admin:      x.Admin,
		Instructor: x.Instructor,
	}, nil
}
