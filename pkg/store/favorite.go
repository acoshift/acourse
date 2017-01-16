package store

import (
	"context"

	"cloud.google.com/go/datastore"
	"github.com/acoshift/acourse/pkg/model"
)

const kindFavorite = "Favorite"

// FavoriteFind finds favorite
func (c *DB) FavoriteFind(ctx context.Context, userID, courseID string) (*model.Favorite, error) {
	var x model.Favorite
	q := datastore.
		NewQuery(kindFavorite).
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

// FavoriteCountCourse counts favorite for a course
func (c *DB) FavoriteCountCourse(ctx context.Context, courseID string) (int, error) {
	q := datastore.
		NewQuery(kindFavorite).
		Filter("CourseID =", courseID).
		KeysOnly()

	keys, err := c.client.GetAll(ctx, q, nil)
	if err != nil {
		return 0, err
	}
	return len(keys), nil
}

// FavoriteAdd adds course favorite to a user
func (c *DB) FavoriteAdd(ctx context.Context, userID, courseID string) error {
	// find duplicate favorite
	x, err := c.FavoriteFind(ctx, userID, courseID)
	if err != nil {
		return err
	}
	if x != nil {
		return nil
	}

	x = &model.Favorite{
		UserID:   userID,
		CourseID: courseID,
	}
	x.Stamp()
	x.SetKey(datastore.IncompleteKey(kindFavorite, nil))
	key, err := c.client.Put(ctx, x.Key(), x)
	if err != nil {
		return err
	}
	x.SetKey(key)
	return nil
}

// FavoriteRemove removes course favorite to a user
func (c *DB) FavoriteRemove(ctx context.Context, userID, courseID string) error {
	// find favorite key
	x, err := c.FavoriteFind(ctx, userID, courseID)
	if err != nil {
		return err
	}
	if x == nil {
		return nil
	}

	return c.client.Delete(ctx, x.Key())
}

// FavoriteCreateAll creates all favorites
func (c *DB) FavoriteCreateAll(ctx context.Context, xs []*model.Favorite) error {
	keys := make([]*datastore.Key, len(xs))
	for i, x := range xs {
		x.Stamp()
		keys[i] = datastore.IncompleteKey(kindFavorite, nil)
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
