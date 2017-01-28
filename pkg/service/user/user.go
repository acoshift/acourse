package user

import (
	"context"
	"log"
	"time"

	"github.com/acoshift/acourse/pkg/acourse"
	"github.com/acoshift/acourse/pkg/internal"
	"github.com/acoshift/ds"
	"github.com/acoshift/gotcha"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// New creates new User service server
func New(client *ds.Client) acourse.UserServiceServer {
	s := &service{client}
	go s.startCacheUser()
	return s
}

type service struct {
	client *ds.Client
}

var cacheUser = gotcha.New()

func (s *service) startCacheUser() {
	ctx := context.Background()
	for {
		var xs []*user
		err := s.client.Query(ctx, kindUser, &xs)
		if err != nil {
			time.Sleep(time.Minute * 10)
			continue
		}
		cacheUser.Purge()
		for _, x := range xs {
			cacheUser.Set(x.ID(), x)
		}
		log.Println("Cached Users")
		time.Sleep(time.Hour * 2)
	}
}

func (s *service) getUser(ctx context.Context, userID string) (*user, error) {
	if c := cacheUser.Get(userID); c != nil {
		return c.(*user), nil
	}

	var x user
	err := s.client.GetByName(ctx, kindUser, userID, &x)
	err = ds.IgnoreFieldMismatch(err)
	if err != nil {
		return nil, err
	}
	cacheUser.Set(userID, &x)
	return &x, nil
}

func (s *service) mustGetUser(ctx context.Context, userID string) (*user, error) {
	x, err := s.getUser(ctx, userID)
	if ds.NotFound(err) {
		x.SetNameID(kindUser, userID)
	}
	err = ds.IgnoreNotFound(err)
	if err != nil {
		return nil, err
	}
	return x, nil
}

func (s *service) findUser(ctx context.Context, username string) (*user, error) {
	var x user
	err := s.client.QueryFirst(ctx, kindUser, &x, ds.Filter("Username =", username))
	err = ds.IgnoreFieldMismatch(err)
	if err != nil {
		return nil, err
	}
	return &x, nil
}

func (s *service) saveUser(ctx context.Context, x *user) error {
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
	cacheUser.Set(x.ID(), x)
	return nil
}

func (s *service) GetUser(ctx context.Context, req *acourse.UserIDRequest) (*acourse.User, error) {
	if req.GetUserId() == "" {
		return nil, ErrUserIDRequired
	}

	x, err := s.mustGetUser(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}
	return toUser(x), nil
}

func (s *service) GetUsers(ctx context.Context, req *acourse.UserIDsRequest) (*acourse.UsersResponse, error) {
	userIDs := req.GetUserIds()
	l := len(userIDs)

	if l == 0 {
		return &acourse.UsersResponse{}, nil
	}

	xs := make([]*user, 0, l)
	ids := make([]string, 0, l)

	// try get in cache first
	for _, id := range userIDs {
		if c := cacheUser.Get(id); c != nil {
			xs = append(xs, c.(*user))
		} else {
			ids = append(ids, id)
		}
	}

	if len(ids) == 0 {
		return &acourse.UsersResponse{Users: toUsers(xs)}, nil
	}

	var ts []*user
	err := s.client.GetByNames(ctx, kindUser, ids, &ts)
	ds.SetNameIDs(kindUser, ids, ts)
	err = ds.IgnoreNotFound(err)
	err = ds.IgnoreFieldMismatch(err)
	if err != nil {
		return nil, err
	}

	for i, x := range ts {
		if x == nil {
			x = &user{}
			x.SetNameID(kindUser, ids[i])
		}
		xs = append(xs, x)
		cacheUser.Set(x.ID(), x)
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
	x.Username = req.GetUsername()
	x.Name = req.GetName()
	x.Photo = req.GetPhoto()
	x.AboutMe = req.GetAboutMe()

	err = s.saveUser(ctx, x)
	if err != nil {
		return nil, err
	}

	return new(acourse.Empty), nil
}
