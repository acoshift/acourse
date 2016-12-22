package store

import (
	"time"

	"cloud.google.com/go/datastore"
)

// Attend model
type Attend struct {
	Base
	Timestampable
	UserID   string
	CourseID string
	At       time.Time `datastore:",noindex"`
}

const kindAttend = "Attend"

// AttendFind finds attend for given user id and course id
func (c *DB) AttendFind(userID, courseID string) (*Attend, error) {
	ctx, cancel := getContext()
	defer cancel()

	var x Attend
	q := datastore.
		NewQuery(kindAttend).
		Filter("UserID =", userID).
		Filter("CourseID =", courseID).
		Limit(1)

	err := c.findFirst(ctx, q, &x)
	if notFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &x, nil
}

// AttendCreateAll creates all attend
func (c *DB) AttendCreateAll(xs []*Attend) error {
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
		x.setKey(keys[i])
	}
	return nil
}

// AttendPurge purges all attends
func (c *DB) AttendPurge() error {
	return c.purge(kindAttend)
}
