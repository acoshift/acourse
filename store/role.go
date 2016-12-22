package store

import (
	"cloud.google.com/go/datastore"
)

// Role model store user's role
type Role struct {
	Base
	Timestampable

	// roles
	Admin      bool
	Instructor bool
}

const kindRole = "Role"

// RoleGet retrieves role by id
func (c *DB) RoleGet(roleID string) (*Role, error) {
	id := idInt(roleID)
	if id == 0 {
		return nil, ErrInvalidID
	}

	ctx, cancel := getContext()
	defer cancel()

	var x Role
	err := c.get(ctx, datastore.IDKey(kindRole, id, nil), &x)
	if notFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &x, nil
}

// RoleFindByUserID retrieves role by user id
func (c *DB) RoleFindByUserID(userID string) (*Role, error) {
	if userID == "" {
		return nil, ErrInvalidID
	}

	ctx, cancel := getContext()
	defer cancel()

	pID := datastore.NameKey(kindUser, userID, nil)

	var x Role
	q := datastore.
		NewQuery(kindRole).
		Ancestor(pID).
		Limit(1)

	err := c.findFirst(ctx, q, &x)
	if notFound(err) {
		x.setKey(datastore.IncompleteKey(kindRole, pID))
		return &x, nil
	}
	if err != nil {
		return nil, err
	}
	return &x, nil
}

// RoleSave saves role for a user
func (c *DB) RoleSave(x *Role) error {
	if x.key == nil {
		return ErrInvalidID
	}
	ctx, cancel := getContext()
	defer cancel()

	x.Stamp()
	_, err := c.client.Put(ctx, x.key, x)
	return err
}

// RolePurge purges all users
func (c *DB) RolePurge() error {
	return c.purge(kindRole)
}
