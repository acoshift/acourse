package store

import "cloud.google.com/go/datastore"

// User model
type User struct {
	Base
	Timestampable
	Username string
	Name     string `datastore:",noindex"`
	Photo    string `datastore:",noindex"`
	AboutMe  string `datastore:",noindex"`
}

const kindUser = "User"

// UserGet retrieves user from database
func (c *DB) UserGet(userID string) (*User, error) {
	ctx, cancel := getContext()
	defer cancel()

	var err error
	var x User

	key := datastore.NameKey(kindUser, userID, nil)
	err = c.client.Get(ctx, key, &x)
	if notFound(err) {
		return nil, nil
	}
	if datastoreError(err) {
		return nil, err
	}
	x.setKey(key)
	return &x, nil
}

// UserFindUsername retrieves user from username from database
func (c *DB) UserFindUsername(username string) (*User, error) {
	ctx, cancel := getContext()
	defer cancel()

	var x User

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
func (c *DB) UserSave(x *User) error {
	if x.key == nil {
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
	_, err := c.client.Put(ctx, x.key, x)
	return err
}

// UserCreateAll creates new users to database
func (c *DB) UserCreateAll(userIDs []string, xs []*User) error {
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
		x.setKey(k)
	}

	// TODO: check duplicated username

	ctx, cancel := getContext()
	defer cancel()
	_, err := c.client.PutMulti(ctx, keys, xs)
	return err
}

// UserCreate creates new user
func (c *DB) UserCreate(userID string, x *User) error {
	if userID == "" {
		return ErrInvalidID
	}
	x.setKey(datastore.NameKey(kindUser, userID, nil))
	return c.UserSave(x)
}

// UserPurge purges all users
func (c *DB) UserPurge() error {
	return c.purge(kindUser)
}
