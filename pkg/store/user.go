package store

import (
	"context"
	"log"
	"time"

	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/ds"
	"github.com/acoshift/gotcha"
)

var cacheUser = gotcha.New()

func (c *DB) initUser() {
	go func() {
		for {
			xs, err := c.UserList(context.Background())
			if err != nil {
				time.Sleep(time.Minute * 10)
				continue
			}
			cacheUser.Purge()
			for _, x := range xs {
				cacheUser.Set(x.ID, x)
			}
			log.Println("Cached Users")
			time.Sleep(time.Hour * 2)
		}
	}()
}

// UserGet retrieves user from database
func (c *DB) UserGet(ctx context.Context, userID string) (*model.User, error) {
	if cache := cacheUser.Get(userID); cache != nil {
		return cache.(*model.User), nil
	}

	var err error
	var x model.User

	err = c.client.GetByName(ctx, kindUser, userID, &x)
	err = ds.IgnoreFieldMismatch(err)
	if ds.NotFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	cacheUser.Set(userID, &x)
	return &x, nil
}

// UserGetMulti retrieves multiple users from database
func (c *DB) UserGetMulti(ctx context.Context, userIDs []string) (model.Users, error) {
	if len(userIDs) == 0 {
		return []*model.User{}, nil
	}

	users := make([]*model.User, 0, len(userIDs))
	ids := make([]string, 0, len(userIDs))

	// try get in cache first
	for _, id := range userIDs {
		if c := cacheUser.Get(id); c != nil {
			users = append(users, c.(*model.User))
		} else {
			ids = append(ids, id)
		}
	}

	if len(ids) == 0 {
		return users, nil
	}

	xs := make([]*model.User, len(ids))
	err := c.client.GetByNames(ctx, kindUser, ids, xs)
	err = ds.IgnoreNotFound(err)
	err = ds.IgnoreFieldMismatch(err)
	if err != nil {
		return nil, err
	}

	for i, x := range xs {
		if x == nil {
			x = &model.User{}
			x.SetNameKey(kindUser, ids[i])
		}
		users = append(users, x)
		cacheUser.Set(x.ID, x)
	}
	return users, nil
}

// UserMustGet retrieves user from database
// if not exists return empty user with given id
func (c *DB) UserMustGet(ctx context.Context, userID string) (*model.User, error) {
	x, err := c.UserGet(ctx, userID)
	if err != nil {
		return nil, err
	}
	if x == nil {
		x = &model.User{}
		x.SetNameKey(kindUser, userID)
	}
	return x, nil
}

// UserFindUsername retrieves user from username from database
func (c *DB) UserFindUsername(ctx context.Context, username string) (*model.User, error) {
	var x model.User

	err := c.client.QueryFirst(ctx, kindUser, &x, ds.Filter("Username =", username))
	err = ds.IgnoreFieldMismatch(err)
	if ds.NotFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &x, nil
}

// UserSave saves user to database
func (c *DB) UserSave(ctx context.Context, x *model.User) error {
	if x.Key() == nil {
		return ErrInvalidID
	}

	// check duplicated username
	if x.Username != "" {
		u, err := c.UserFindUsername(ctx, x.Username)
		if err != nil {
			return err
		}
		if u != nil && x.ID != u.ID {
			return ErrConflict("username already exists")
		}
	}

	err := c.client.Save(ctx, kindUser, x)
	cacheUser.Unset(x.ID)
	return err
}

// UserCreateAll creates new users to database
func (c *DB) UserCreateAll(ctx context.Context, userIDs []string, xs []*model.User) error {
	if len(userIDs) != len(xs) {
		return ErrConflict("user id count not match user count")
	}

	// validate keys, stamp model, and get keys
	for i, x := range xs {
		if userIDs[i] == "" {
			return ErrInvalidID
		}
		x.SetNameKey(kindUser, userIDs[i])
	}

	// TODO: check duplicated username

	err := c.client.SaveMulti(ctx, kindUser, xs)
	if err != nil {
		return err
	}
	return nil
}

// UserCreate creates new user
func (c *DB) UserCreate(ctx context.Context, userID string, x *model.User) error {
	if userID == "" {
		return ErrInvalidID
	}
	x.SetNameKey(kindUser, userID)
	err := c.UserSave(ctx, x)
	if err != nil {
		return err
	}
	cacheUser.Set(x.ID, x)
	return nil
}

// UserList retrieves all users
func (c *DB) UserList(ctx context.Context) ([]*model.User, error) {
	var xs []*model.User
	err := c.client.Query(ctx, kindUser, &xs)
	err = ds.IgnoreFieldMismatch(err)
	if err != nil {
		return nil, err
	}
	return xs, nil
}
