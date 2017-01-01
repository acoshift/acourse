package store

import (
	"time"

	"cloud.google.com/go/datastore"
	"github.com/acoshift/acourse/model"
)

const kindEnroll = "Enroll"

var cacheEnrollCount = NewCache(time.Hour)

// EnrollFind finds enroll for given user id and course id
func (c *DB) EnrollFind(userID, courseID string) (*model.Enroll, error) {
	ctx, cancel := getContext()
	defer cancel()

	var x model.Enroll
	q := datastore.
		NewQuery(kindEnroll).
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

// EnrollListByUserID list all enroll by given user id
func (c *DB) EnrollListByUserID(userID string) ([]*model.Enroll, error) {
	ctx, cancel := getContext()
	defer cancel()

	var xs []*model.Enroll
	q := datastore.
		NewQuery(kindEnroll).
		Filter("UserID =", userID)

	keys, err := c.getAll(ctx, q, &xs)
	if err != nil {
		return nil, err
	}
	for i, x := range xs {
		x.SetKey(keys[i])
	}
	return xs, nil
}

// EnrollSave saves enroll to database
func (c *DB) EnrollSave(x *model.Enroll) error {
	ctx, cancel := getContext()
	defer cancel()

	var pKey *datastore.PendingKey
	x.Stamp()

	commit, err := c.client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		var t model.Enroll

		q := datastore.
			NewQuery(kindEnroll).
			Filter("UserID =", x.UserID).
			Filter("CourseID =", x.CourseID).
			Limit(1).
			Transaction(tx)

		err := c.findFirst(ctx, q, &t)
		if err == nil {
			return ErrConflict("enroll already exists")
		}

		pKey, err = tx.Put(datastore.IncompleteKey(kindEnroll, nil), x)
		return err
	})
	if err != nil {
		return err
	}
	x.SetKey(commit.Key(pKey))
	cacheEnrollCount.Del(x.CourseID)
	return nil
}

// EnrollCreateAll creates all enrolls
func (c *DB) EnrollCreateAll(xs []*model.Enroll) error {
	ctx, cancel := getContext()
	defer cancel()

	keys := make([]*datastore.Key, len(xs))
	for i, x := range xs {
		x.Stamp()
		keys[i] = datastore.IncompleteKey(kindEnroll, nil)
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

// EnrollPurge purges all users
func (c *DB) EnrollPurge() error {
	return c.purge(kindEnroll)
}

// EnrollCourseCount counts enroll from course id
func (c *DB) EnrollCourseCount(courseID string) (int, error) {
	if cache := cacheEnrollCount.Get(courseID); cache != nil {
		return cache.(int), nil
	}

	ctx, cancel := getContext()
	defer cancel()

	q := datastore.
		NewQuery(kindEnroll).
		Filter("CourseID =", courseID).
		KeysOnly()

	keys, err := c.client.GetAll(ctx, q, nil)
	if err != nil {
		return 0, err
	}
	r := len(keys)

	cacheEnrollCount.Set(courseID, r)
	return r, nil
}
