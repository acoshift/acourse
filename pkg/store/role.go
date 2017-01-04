package store

import (
	"cloud.google.com/go/datastore"
	"github.com/acoshift/acourse/pkg/model"
)

const kindRole = "Role"

// RoleGet retrieves role by user id
func (c *DB) RoleGet(userID string) (*model.Role, error) {
	if userID == "" {
		return &model.Role{}, nil
	}

	ctx, cancel := getContext()
	defer cancel()

	var x model.Role
	err := c.get(ctx, datastore.NameKey(kindRole, userID, nil), &x)
	if notFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &x, nil
}

// RoleSave saves role for a user
func (c *DB) RoleSave(x *model.Role) error {
	if x.Key() == nil {
		return ErrInvalidID
	}
	ctx, cancel := getContext()
	defer cancel()

	x.Stamp()
	_, err := c.client.Put(ctx, x.Key(), x)
	return err
}

// RolePurge purges all users
func (c *DB) RolePurge() error {
	return c.purge(kindRole)
}
