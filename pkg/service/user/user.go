package user

import (
	"github.com/acoshift/acourse/pkg/acourse"
	"github.com/acoshift/acourse/pkg/app"
	"github.com/acoshift/acourse/pkg/internal"
	"github.com/acoshift/ds"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// New creates new User service server
func New(client *ds.Client) acourse.UserServiceServer {
	s := &service{client}
	return s
}

type service struct {
	client *ds.Client
}

func (s *service) getUser(ctx context.Context, userID string) (*userModel, error) {
	var x userModel
	err := s.client.GetByName(ctx, kindUser, userID, &x)
	err = ds.IgnoreFieldMismatch(err)
	if err != nil {
		return nil, err
	}
	return &x, nil
}

func (s *service) mustGetUser(ctx context.Context, userID string) (*userModel, error) {
	x, err := s.getUser(ctx, userID)
	if ds.NotFound(err) {
		x = &userModel{}
		x.SetNameID(kindUser, userID)
	}
	err = ds.IgnoreNotFound(err)
	if err != nil {
		return nil, err
	}
	return x, nil
}

func (s *service) findUser(ctx context.Context, username string) (*userModel, error) {
	var x userModel
	err := s.client.QueryFirst(ctx, kindUser, &x, ds.Filter("Username =", username))
	err = ds.IgnoreFieldMismatch(err)
	if err != nil {
		return nil, err
	}
	return &x, nil
}

func (s *service) saveUser(ctx context.Context, x *userModel) error {
	if x.GetKey() == nil {
		return ErrUserIDRequired
	}

	// check duplicated username
	if x.Username != "" {
		u, err := s.findUser(ctx, x.Username)
		if !ds.NotFound(err) && err != nil {
			return err
		}
		if u != nil && x.ID() != u.ID() {
			return ErrUserNameConflict
		}
	}

	err := s.client.SaveModel(ctx, "", x)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetUser(ctx context.Context, req *acourse.UserIDRequest) (*acourse.User, error) {
	if len(req.UserId) == 0 {
		return nil, ErrUserIDRequired
	}

	x, err := s.mustGetUser(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return toUser(x), nil
}

func (s *service) GetUsers(ctx context.Context, req *acourse.UserIDsRequest) (*acourse.UsersResponse, error) {
	ids := req.UserIds
	if len(ids) == 0 {
		return &acourse.UsersResponse{}, nil
	}

	ids = app.UniqueIDs(ids)

	var xs []*userModel
	err := s.client.GetByNames(ctx, kindUser, ids, &xs)
	err = ds.IgnoreNotFound(err)
	err = ds.IgnoreFieldMismatch(err)
	if err != nil {
		return nil, err
	}

	for i := range xs {
		if xs[i] == nil {
			xs[i] = &userModel{}
			xs[i].SetNameID(kindUser, ids[i])
		}
	}
	return &acourse.UsersResponse{Users: toUsers(xs)}, nil
}

func (s *service) GetMe(ctx context.Context, req *acourse.Empty) (*acourse.GetMeResponse, error) {
	userID := internal.GetUserID(ctx)
	if userID == "" {
		return nil, grpc.Errorf(codes.Unauthenticated, "authorization required")
	}
	user, err := s.mustGetUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	role, err := s.GetRole(ctx, &acourse.UserIDRequest{UserId: userID})
	if err != nil {
		return nil, err
	}
	return &acourse.GetMeResponse{
		User: toUser(user),
		Role: role,
	}, nil
}

func (s *service) UpdateMe(ctx context.Context, req *acourse.User) (*acourse.Empty, error) {
	userID := internal.GetUserID(ctx)
	if userID == "" {
		return nil, grpc.Errorf(codes.Unauthenticated, "authorization required")
	}
	x, err := s.mustGetUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	x.Username = req.Username
	x.Name = req.Name
	x.Photo = req.Photo
	x.AboutMe = req.AboutMe

	err = s.saveUser(ctx, x)
	if err != nil {
		return nil, err
	}

	return new(acourse.Empty), nil
}
