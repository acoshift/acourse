package user

import (
	"context"

	"github.com/acoshift/acourse/pkg/app"
	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/acourse/pkg/view"
	"github.com/acoshift/httperror"
)

// New creates new user service
func New(store Store) app.UserService {
	return &service{store}
}

// Store is the store interface for user service
type Store interface {
	UserGetMulti(context.Context, []string) ([]*model.User, error)
	UserMustGet(string) (*model.User, error)
	UserSave(*model.User) error
	RoleGet(string) (*model.Role, error)
}

type service struct {
	store Store
}

func (s *service) GetUsers(ctx context.Context, req *app.IDsRequest) (*app.UsersReply, error) {
	users, err := s.store.UserGetMulti(ctx, req.IDs)
	if err != nil {
		return nil, err
	}
	return &app.UsersReply{Users: view.ToUserCollection(users)}, nil
}

func (s *service) GetMe(ctx context.Context) (*app.UserMeReply, error) {
	currentUserID, ok := ctx.Value(app.KeyCurrentUserID).(string)
	if !ok {
		return nil, httperror.Unauthorized
	}
	user, err := s.store.UserMustGet(currentUserID)
	if err != nil {
		return nil, err
	}
	role, err := s.store.RoleGet(currentUserID)
	if err != nil {
		return nil, err
	}
	return &app.UserMeReply{User: view.ToUserMe(user, view.ToRole(role))}, nil
}

func (s *service) UpdateMe(ctx context.Context, req *app.UserRequest) error {
	currentUserID, ok := ctx.Value(app.KeyCurrentUserID).(string)
	if !ok {
		return httperror.Unauthorized
	}
	if req.User == nil {
		return httperror.BadRequest
	}
	if err := req.User.Validate(); err != nil {
		return httperror.BadRequestWith(err)
	}
	user, err := s.store.UserMustGet(currentUserID)
	if err != nil {
		return err
	}
	if req.User.Username != nil {
		user.Username = *req.User.Username
	}
	if req.User.Name != nil {
		user.Name = *req.User.Name
	}
	if req.User.Photo != nil {
		user.Photo = *req.User.Photo
	}
	if req.User.AboutMe != nil {
		user.AboutMe = *req.User.AboutMe
	}
	return s.store.UserSave(user)
}
