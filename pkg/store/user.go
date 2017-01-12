package store

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/gotcha"
)

const kindUser = "User"

var cacheUser = gotcha.New()

func (c *DB) initUser() {
	go func() {
		for {
			xs, err := c.UserList()
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
func (c *DB) UserGet(userID string) (*model.User, error) {
	if cache := cacheUser.Get(userID); cache != nil {
		return cache.(*model.User), nil
	}

	ctx, cancel := getContext()
	defer cancel()

	var err error
	var x model.User

	key := datastore.NameKey(kindUser, userID, nil)
	err = c.client.Get(ctx, key, &x)
	if notFound(err) {
		return nil, nil
	}
	if datastoreError(err) {
		return nil, err
	}
	x.SetKey(key)

	cacheUser.Set(userID, &x)
	return &x, nil
}

// UserGetMulti retrieves multiple users from database
func (c *DB) UserGetMulti(ctx context.Context, userIDs []string) (model.Users, error) {
	if len(userIDs) == 0 {
		return []*model.User{}, nil
	}

	users := make([]*model.User, 0, len(userIDs))
	keys := make([]*datastore.Key, 0, len(userIDs))

	// try get in cache first
	for _, id := range userIDs {
		if c := cacheUser.Get(id); c != nil {
			users = append(users, c.(*model.User))
		} else {
			keys = append(keys, datastore.NameKey(kindUser, id, nil))
		}
	}

	if len(keys) == 0 {
		return users, nil
	}

	xs := make([]*model.User, len(keys))
	err := c.client.GetMulti(ctx, keys, xs)
	if multiError(err) {
		return nil, err
	}

	for i, x := range xs {
		if x == nil {
			x = &model.User{}
		}
		x.SetKey(keys[i])
		users = append(users, x)
		cacheUser.Set(x.ID, x)
	}
	return users, nil
}

// UserMustGet retrieves user from database
// if not exists return empty user with given id
func (c *DB) UserMustGet(userID string) (*model.User, error) {
	x, err := c.UserGet(userID)
	if err != nil {
		return nil, err
	}
	if x == nil {
		x = &model.User{}
		x.SetKey(datastore.NameKey(kindUser, userID, nil))
	}
	return x, nil
}

// UserFindUsername retrieves user from username from database
func (c *DB) UserFindUsername(username string) (*model.User, error) {
	ctx, cancel := getContext()
	defer cancel()

	var x model.User

	q := datastore.
		NewQuery(kindUser).
		Filter("Username =", username).
		Limit(1)

	err := c.findFirst(ctx, q, &x)
	if notFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &x, nil
}

// UserSave saves user to database
func (c *DB) UserSave(x *model.User) error {
	if x.Key() == nil {
		return ErrInvalidID
	}

	// check duplicated username
	if x.Username != "" {
		u, err := c.UserFindUsername(x.Username)
		if err != nil {
			return err
		}
		if u != nil && x.ID != u.ID {
			return ErrConflict("username already exists")
		}
	}

	ctx, cancel := getContext()
	defer cancel()
	x.Stamp()
	_, err := c.client.Put(ctx, x.Key(), x)
	cacheUser.Unset(x.ID)
	return err
}

// UserCreateAll creates new users to database
func (c *DB) UserCreateAll(userIDs []string, xs []*model.User) error {
	if len(userIDs) != len(xs) {
		return ErrConflict("user id count not match user count")
	}

	// validate keys, stamp model, and get keys
	keys := make([]*datastore.Key, len(xs))
	for i, x := range xs {
		if userIDs[i] == "" {
			return ErrInvalidID
		}
		x.Stamp()
		k := datastore.NameKey(kindUser, userIDs[i], nil)
		keys[i] = k
		x.SetKey(k)
	}

	// TODO: check duplicated username

	ctx, cancel := getContext()
	defer cancel()
	_, err := c.client.PutMulti(ctx, keys, xs)
	return err
}

// UserCreate creates new user
func (c *DB) UserCreate(userID string, x *model.User) error {
	if userID == "" {
		return ErrInvalidID
	}
	x.SetKey(datastore.NameKey(kindUser, userID, nil))
	err := c.UserSave(x)
	if err != nil {
		return err
	}
	cacheUser.Set(x.ID, x)
	return nil
}

// UserPurge purges all users
func (c *DB) UserPurge() error {
	return c.purge(kindUser)
}

// UserList retrieves all users
func (c *DB) UserList() ([]*model.User, error) {
	var xs []*model.User
	ctx, cancel := getLongContext()
	defer cancel()
	keys, err := c.getAll(ctx, datastore.NewQuery(kindUser), &xs)
	if err != nil {
		return nil, err
	}
	for i := range keys {
		xs[i].SetKey(keys[i])
	}
	return xs, nil
}
