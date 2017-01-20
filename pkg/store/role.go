package store

import (
	"context"
	"log"
	"time"

	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/gotcha"
)

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
	// var x model.Role
	x, _ := cacheRole.Get(userID).(*model.Role)
	if x == nil {
		return &model.Role{}, nil
	}
	return x, nil
	// err := c.getByName(ctx, kindRole, userID, &x)
	// log.Println(err)
	// if err == ErrNotFound {
	// 	return &model.Role{}, nil
	// }
	// if err != nil {
	// 	return nil, err
	// }
	// return &x, nil
}

// RoleSave saves role for a user
func (c *DB) RoleSave(ctx context.Context, x *model.Role) error {
	if x.Key() == nil {
		return ErrInvalidID
	}
	err := c.client.Save(ctx, kindRole, x)
	if err != nil {
		return err
	}
	cacheRole.Unset(x.ID)
	return nil
}

// RoleList retrieves all role in database
func (c *DB) RoleList(ctx context.Context) ([]*model.Role, error) {
	var xs []*model.Role
	err := c.client.Query(ctx, kindRole, &xs)
	if err != nil {
		return nil, err
	}
	return xs, nil
}
