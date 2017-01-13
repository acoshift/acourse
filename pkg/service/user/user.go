package user

import (
	"context"

	"github.com/acoshift/acourse/pkg/acourse"
	"github.com/acoshift/acourse/pkg/model"
	rctx "golang.org/x/net/context"
)

// New creates new service
func New(store Store) acourse.UserServiceServer {
	return &userServiceServer{store}
}

// Store is the store interface for user service
type Store interface {
	UserGetMulti(context.Context, []string) (model.Users, error)
	UserMustGet(string) (*model.User, error)
	UserSave(*model.User) error
	RoleGet(string) (*model.Role, error)
}

type userServiceServer struct {
	store Store
}

func (s *userServiceServer) GetUser(ctx rctx.Context, req *acourse.GetUserRequest) (*acourse.GetUserResponse, error) {
	return nil, nil
}
