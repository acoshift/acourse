package store

import (
	"acourse/model"

	"cloud.google.com/go/datastore"
)

const kindRole = "Role"

// RoleGet retrieves role by id
func (c *DB) RoleGet(roleID string) (*model.Role, error) {
	id := idInt(roleID)
	if id == 0 {
		return nil, ErrInvalidID
	}

	ctx, cancel := getContext()
	defer cancel()

	var x model.Role
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
func (c *DB) RoleFindByUserID(userID string) (*model.Role, error) {
	if userID == "" {
		return nil, ErrInvalidID
	}

	ctx, cancel := getContext()
	defer cancel()

	pID := datastore.NameKey(kindUser, userID, nil)

	var x model.Role
	q := datastore.
		NewQuery(kindRole).
		Ancestor(pID).
		Limit(1)

	err := c.findFirst(ctx, q, &x)
	if notFound(err) {
		x.SetKey(datastore.IncompleteKey(kindRole, pID))
		return &x, nil
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
