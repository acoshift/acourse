package store

import (
	"context"

	"cloud.google.com/go/datastore"
	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/gotcha"
)

const kindEnroll = "Enroll"

var cacheEnrollCount = gotcha.New()

// EnrollFind finds enroll for given user id and course id
func (c *DB) EnrollFind(ctx context.Context, userID, courseID string) (*model.Enroll, error) {
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
func (c *DB) EnrollListByUserID(ctx context.Context, userID string) (model.Enrolls, error) {
	var xs []*model.Enroll
	q := datastore.
		NewQuery(kindEnroll).
		Filter("UserID =", userID)

	err := c.getAll(ctx, q, &xs)
	if err != nil {
		return nil, err
	}

	return xs, nil
}

// EnrollSave saves enroll to database
func (c *DB) EnrollSave(ctx context.Context, x *model.Enroll) error {
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
	cacheEnrollCount.Unset(x.CourseID)
	return nil
}

// EnrollCreateAll creates all enrolls
func (c *DB) EnrollCreateAll(ctx context.Context, xs []*model.Enroll) error {
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

// EnrollCourseCount counts enroll from course id
func (c *DB) EnrollCourseCount(ctx context.Context, courseID string) (int, error) {
	if cache := cacheEnrollCount.Get(courseID); cache != nil {
		return cache.(int), nil
	}

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

// EnrollSaveMulti saves multiple enrolls to database
func (c *DB) EnrollSaveMulti(ctx context.Context, enrolls []*model.Enroll) error {
	keys := make([]*datastore.Key, 0, len(enrolls))

	for _, enroll := range enrolls {
		enroll.Stamp()
		keys = append(keys, datastore.IncompleteKey(kindEnroll, nil))
	}

	var pKey []*datastore.PendingKey

	commit, err := c.client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		var t model.Enroll
		var err error
		for _, enroll := range enrolls {
			q := datastore.
				NewQuery(kindEnroll).
				Filter("UserID =", enroll.UserID).
				Filter("CourseID =", enroll.CourseID).
				Limit(1).
				Transaction(tx)

			err = c.findFirst(ctx, q, &t)
			if err == nil {
				return ErrConflict("enroll already exists")
			}
		}

		pKey, err = tx.PutMulti(keys, enrolls)
		return err
	})
	if err != nil {
		return err
	}

	for i, enroll := range enrolls {
		enroll.SetKey(commit.Key(pKey[i]))
		cacheEnrollCount.Unset(enroll.CourseID)
	}

	return nil
}
