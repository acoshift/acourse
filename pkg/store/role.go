package store

import (
	"log"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/gotcha"
)

const kindRole = "Role"

var cacheRole = gotcha.New()

func (c *DB) initRole() {
	go func() {
		for {
			xs, err := c.RoleList()
			if err != nil {
				time.Sleep(time.Minute * 10)
				continue
			}
			cacheRole.Purge()
			for _, x := range xs {
				cacheRole.Set(x.ID, x)
			}
			log.Println("Cached Roles")
			time.Sleep(time.Hour)
		}
	}()
}

// RoleGet retrieves role by id
func (c *DB) RoleGet(roleID string) (*model.Role, error) {
	id := idInt(roleID)
	if id == 0 {
		return &model.Role{}, nil
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
		return &model.Role{}, nil
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
	if err != nil {
		return err
	}
	return nil
}

// RolePurge purges all users
func (c *DB) RolePurge() error {
	return c.purge(kindRole)
}

// RoleList retrieves all role in database
func (c *DB) RoleList() ([]*model.Role, error) {
	var xs []*model.Role
	ctx, cancel := getLongContext()
	defer cancel()
	keys, err := c.getAll(ctx, datastore.NewQuery(kindRole), &xs)
	if err != nil {
		return nil, err
	}
	for i := range keys {
		xs[i].SetKey(keys[i])
	}
	return xs, nil
}
