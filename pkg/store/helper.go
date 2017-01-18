package store

import (
	"context"
	"strconv"

	"cloud.google.com/go/datastore"
	"github.com/acoshift/acourse/pkg/model"
	"google.golang.org/api/iterator"
)

func datastoreError(err error) bool {
	if err == nil {
		return false
	}
	if _, ok := err.(*datastore.ErrFieldMismatch); ok {
		// ignore field mismatch
		return false
	}
	// check multi errors
	if errs, ok := err.(datastore.MultiError); ok {
		hasError := false
		for _, err := range errs {
			if datastoreError(err) {
				hasError = true
				break
			}
		}
		return hasError
	}
	return true
}

func notFound(err error) bool {
	return err == iterator.Done || err == datastore.ErrNoSuchEntity
}

func multiError(err error) bool {
	if err == nil {
		return false
	}
	if errs, ok := err.(datastore.MultiError); ok {
		hasError := false
		for _, err := range errs {
			if _, ok := err.(*datastore.ErrFieldMismatch); err != datastore.ErrNoSuchEntity && !ok {
				hasError = true
				break
			}
		}
		if hasError {
			return true
		}
	}
	return false
}

func (c *DB) getAll(ctx context.Context, q *datastore.Query, dst interface{}) ([]*datastore.Key, error) {
	keys, err := c.client.GetAll(ctx, q, dst)
	if datastoreError(err) {
		return nil, err
	}
	return keys, nil
}

func (c *DB) get(ctx context.Context, key *datastore.Key, dst model.KeySetter) error {
	err := c.client.Get(ctx, key, dst)
	if datastoreError(err) {
		return err
	}
	dst.SetKey(key)
	return nil
}

func (c *DB) findFirst(ctx context.Context, q *datastore.Query, dst model.KeySetter) error {
	key, err := c.client.Run(ctx, q).Next(dst)
	if datastoreError(err) {
		return err
	}
	dst.SetKey(key)
	return nil
}

func (c *DB) put(ctx context.Context, src model.KeyGetSetter) error {
	key, err := c.client.Put(ctx, src.Key(), src)
	if err != nil {
		return err
	}
	src.SetKey(key)
	return nil
}

func (c *DB) deleteByID(ctx context.Context, kind string, id int64) error {
	return c.client.Delete(ctx, datastore.IDKey(kind, id, nil))
}

func (c *DB) deleteByIDStr(ctx context.Context, kind string, id string) error {
	return c.deleteByID(ctx, kind, idInt(id))
}

func (c *DB) deleteByName(ctx context.Context, kind string, name string) error {
	return c.client.Delete(ctx, datastore.NameKey(kind, name, nil))
}

func (c *DB) save(ctx context.Context, kind string, src model.KeyGetSetter) error {
	if src.Key() == nil {
		src.NewKey(kind)
	}
	return c.put(ctx, src)
}

func idStr(id int64) string {
	return strconv.FormatInt(id, 10)
}

func idInt(id string) int64 {
	r, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0
	}
	return r
}
