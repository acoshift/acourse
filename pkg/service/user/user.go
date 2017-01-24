package user

import (
	"context"

	"github.com/acoshift/acourse/pkg/acourse"
	"github.com/acoshift/acourse/pkg/internal"
	"github.com/acoshift/acourse/pkg/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// New creates new User service server
func New(store Store) acourse.UserServiceServer {
	return &userServiceServer{store}
}

// Store is the store interface for user service
type Store interface {
	UserGetMulti(context.Context, []string) (model.Users, error)
	UserMustGet(context.Context, string) (*model.User, error)
	UserSave(context.Context, *model.User) error
	RoleGet(context.Context, string) (*model.Role, error)
}

type userServiceServer struct {
	store Store
}

func (s *userServiceServer) GetUser(ctx context.Context, req *acourse.GetUserRequest) (*acourse.GetUserResponse, error) {
	users, err := s.store.UserGetMulti(ctx, req.GetUserIds())
	if err != nil {
		return nil, err
	}
	return &acourse.GetUserResponse{Users: acourse.ToUsers(users)}, nil
}

func (s *userServiceServer) GetMe(ctx context.Context, req *acourse.Empty) (*acourse.GetMeResponse, error) {
	userID := internal.GetUserID(ctx)
	if userID == "" {
		return nil, grpc.Errorf(codes.Unauthenticated, "authorization required")
	}
	user, err := s.store.UserMustGet(ctx, userID)
	if err != nil {
		return nil, err
	}
	role, err := s.store.RoleGet(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &acourse.GetMeResponse{
		User: acourse.ToUser(user),
		Role: acourse.ToRole(role),
	}, nil
}

func (s *userServiceServer) UpdateMe(ctx context.Context, req *acourse.User) (*acourse.Empty, error) {
	userID := internal.GetUserID(ctx)
	if userID == "" {
		return nil, grpc.Errorf(codes.Unauthenticated, "authorization required")
	}
	user, err := s.store.UserMustGet(ctx, userID)
	if err != nil {
		return nil, err
	}
	user.Username = req.GetUsername()
	user.Name = req.GetName()
	user.Photo = req.GetPhoto()
	user.AboutMe = req.GetAboutMe()

	err = s.store.UserSave(ctx, user)
	if err != nil {
		return nil, err
	}

	return new(acourse.Empty), nil
}
