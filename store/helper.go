package store

import (
	"context"
	"strconv"
	"time"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
)

func getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Second*30)
}

func datastoreError(err error) bool {
	if err == nil {
		return false
	}
	if _, ok := err.(*datastore.ErrFieldMismatch); ok {
		// ignore field mismatch
		return false
	}
	return true
}

func notFound(err error) bool {
	return err == iterator.Done || err == datastore.ErrNoSuchEntity
}

func (c *DB) getAll(ctx context.Context, q *datastore.Query, dst interface{}) ([]*datastore.Key, error) {
	keys, err := c.client.GetAll(ctx, q, dst)
	if datastoreError(err) {
		return nil, err
	}
	return keys, nil
}

func (c *DB) get(ctx context.Context, key *datastore.Key, dst isBase) error {
	err := c.client.Get(ctx, key, dst)
	if datastoreError(err) {
		return err
	}
	dst.setKey(key)
	return nil
}

func (c *DB) findFirst(ctx context.Context, q *datastore.Query, dst isBase) error {
	key, err := c.client.Run(ctx, q).Next(dst)
	if datastoreError(err) {
		return err
	}
	dst.setKey(key)
	return nil
}

func (c *DB) purge(kind string) error {
	ctx, cancel := getContext()
	defer cancel()
	keys, err := c.client.GetAll(ctx, datastore.NewQuery(kind).KeysOnly(), nil)
	if err != nil {
		return err
	}
	return c.client.DeleteMulti(ctx, keys)
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
