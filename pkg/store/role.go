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
func (c *DB) RoleGet(userID string) (*model.Role, error) {
	if userID == "" {
		return &model.Role{}, nil
	}

	ctx, cancel := getContext()
	defer cancel()

	var x model.Role
	err := c.get(ctx, datastore.NameKey(kindRole, userID, nil), &x)
	if notFound(err) {
		return &model.Role{}, nil
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
