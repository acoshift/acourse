package store

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/acoshift/acourse/pkg/model"
)

const kindAttend = "Attend"

// AttendFind finds attends for given user id and course id
func (c *DB) AttendFind(ctx context.Context, userID, courseID string) (*model.Attend, error) {
	q := datastore.
		NewQuery(kindAttend).
		Filter("UserID =", userID).
		Filter("CourseID =", courseID).
		Filter("CreatedAt >=", time.Now().Add(-6*time.Hour)).
		Limit(1)

	var x model.Attend
	err := c.findFirst(ctx, q, &x)
	if notFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &x, nil
}

// AttendSave saves attend to database
func (c *DB) AttendSave(ctx context.Context, attend *model.Attend) error {
	attend.Stamp()
	if attend.Key() == nil {
		attend.SetKey(datastore.IncompleteKey(kindAttend, nil))
	}

	key, err := c.client.Put(ctx, attend.Key(), attend)
	if err != nil {
		return err
	}
	attend.SetKey(key)
	return nil
}

// AttendCreateAll creates all attend
func (c *DB) AttendCreateAll(xs []*model.Attend) error {
	ctx, cancel := getContext()
	defer cancel()

	keys := make([]*datastore.Key, len(xs))
	for i, x := range xs {
		x.Stamp()
		keys[i] = datastore.IncompleteKey(kindAttend, nil)
	}
	var err error
	keys, err = c.client.PutMulti(ctx, keys, xs)
	if err != nil {
		return err
	}
	for i, x := range xs {
		x.SetKey(keys[i])
	}
	return nil
}

// AttendPurge purges all attends
func (c *DB) AttendPurge() error {
	return c.purge(kindAttend)
}
