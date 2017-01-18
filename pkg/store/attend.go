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
		Filter("CreatedAt >=", time.Now().Add(-6*time.Hour))

	var x model.Attend
	err := c.getFirst(ctx, q, &x)
	if err == ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &x, nil
}

// AttendSave saves attend to database
func (c *DB) AttendSave(ctx context.Context, x *model.Attend) error {
	x.Stamp()
	return c.save(ctx, kindAttend, x)
}

// AttendCreateAll creates all attend
func (c *DB) AttendCreateAll(ctx context.Context, xs []*model.Attend) error {
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
