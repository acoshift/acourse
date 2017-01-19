package store

import (
	"context"

	"cloud.google.com/go/datastore"
	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/ds"
	"github.com/acoshift/gotcha"
)

var cacheEnrollCount = gotcha.New()

// EnrollFind finds enroll for given user id and course id
func (c *DB) EnrollFind(ctx context.Context, userID, courseID string) (*model.Enroll, error) {
	var x model.Enroll

	err := c.client.QueryFirst(ctx, &x,
		ds.Filter("UserID =", userID),
		ds.Filter("CourseID =", courseID),
	)

	if ds.NotFound(err) {
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

	err := c.client.Query(ctx, &model.Enroll{}, &xs,
		ds.Filter("UserID =", userID),
	)
	if err != nil {
		return nil, err
	}

	return xs, nil
}

// EnrollSave saves enroll to database
func (c *DB) EnrollSave(ctx context.Context, x *model.Enroll) error {
	// TODO: race condition
	// TODO: use keysonly query
	var t model.Enroll
	err := c.client.QueryFirst(ctx, &t,
		ds.Filter("UserID =", x.UserID),
		ds.Filter("CourseID =", x.CourseID),
	)
	if err == nil {
		return ErrConflict("enroll already exists")
	}

	err = c.client.Save(ctx, x)
	if err != nil {
		return err
	}

	cacheEnrollCount.Unset(x.CourseID)
	return nil
}

// EnrollCourseCount counts enroll from course id
func (c *DB) EnrollCourseCount(ctx context.Context, courseID string) (int, error) {
	if cache := cacheEnrollCount.Get(courseID); cache != nil {
		return cache.(int), nil
	}

	keys, err := c.client.QueryKeys(ctx, &model.Enroll{}, ds.Filter("CourseID =", courseID))
	if err != nil {
		return 0, err
	}
	r := len(keys)

	cacheEnrollCount.Set(courseID, r)
	return r, nil
}

// EnrollSaveMulti saves multiple enrolls to database
func (c *DB) EnrollSaveMulti(ctx context.Context, enrolls []*model.Enroll) error {
	// TODO: change to ds
	keys := make([]*datastore.Key, 0, len(enrolls))
	kind := (&model.Enroll{}).Kind()
	for _, enroll := range enrolls {
		enroll.Stamp()
		keys = append(keys, datastore.IncompleteKey(kind, nil))
	}

	var pKeys []*datastore.PendingKey

	commit, err := c.client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		var t model.Enroll
		var err error
		for _, enroll := range enrolls {
			err = c.client.QueryFirst(ctx, &t,
				ds.Filter("UserID =", enroll.UserID),
				ds.Filter("CourseID =", enroll.CourseID),
				ds.Transaction(tx),
			)
			if err == nil {
				return ErrConflict("enroll already exists")
			}
		}

		pKeys, err = tx.PutMulti(keys, enrolls)
		return err
	})
	if err != nil {
		return err
	}

	ds.SetCommitKeys(commit, pKeys, enrolls)
	for _, enroll := range enrolls {
		cacheEnrollCount.Unset(enroll.CourseID)
	}

	return nil
}
