package user

import (
	"context"

	"github.com/acoshift/acourse/pkg/acourse"
	"github.com/acoshift/acourse/pkg/internal"
	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/ds"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// New creates new User service server
func New(store Store, client *ds.Client) acourse.UserServiceServer {
	s := &service{client, store}
	go s.startCacheRole()
	return s
}

// Store is the store interface for user service
type Store interface {
	UserGetMulti(context.Context, []string) (model.Users, error)
	UserMustGet(context.Context, string) (*model.User, error)
	UserSave(context.Context, *model.User) error
}

type service struct {
	client *ds.Client
	store  Store
}

func (s *service) GetUser(ctx context.Context, req *acourse.GetUserRequest) (*acourse.GetUserResponse, error) {
	users, err := s.store.UserGetMulti(ctx, req.GetUserIds())
	if err != nil {
		return nil, err
	}
	return &acourse.GetUserResponse{Users: acourse.ToUsers(users)}, nil
}

func (s *service) GetMe(ctx context.Context, req *acourse.Empty) (*acourse.GetMeResponse, error) {
	userID := internal.GetUserID(ctx)
	if userID == "" {
		return nil, grpc.Errorf(codes.Unauthenticated, "authorization required")
	}
	user, err := s.store.UserMustGet(ctx, userID)
	if err != nil {
		return nil, err
	}
	role, err := s.GetRole(ctx, &acourse.UserIDRequest{UserId: userID})
	if err != nil {
		return nil, err
	}
	return &acourse.GetMeResponse{
		User: acourse.ToUser(user),
		Role: role,
	}, nil
}

func (s *service) UpdateMe(ctx context.Context, req *acourse.User) (*acourse.Empty, error) {
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
