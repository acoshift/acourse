package store

import (
	"context"
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
			xs, err := c.RoleList(context.Background())
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
func (c *DB) RoleGet(ctx context.Context, userID string) (*model.Role, error) {
	if userID == "" {
		return &model.Role{}, nil
	}

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
func (c *DB) RoleSave(ctx context.Context, x *model.Role) error {
	if x.Key() == nil {
		return ErrInvalidID
	}

	x.Stamp()
	_, err := c.client.Put(ctx, x.Key(), x)
	if err != nil {
		return err
	}
	return nil
}

// RoleList retrieves all role in database
func (c *DB) RoleList(ctx context.Context) ([]*model.Role, error) {
	var xs []*model.Role
	keys, err := c.getAll(ctx, datastore.NewQuery(kindRole), &xs)
	if err != nil {
		return nil, err
	}
	for i := range keys {
		xs[i].SetKey(keys[i])
	}
	return xs, nil
}
