package store

import (
	"context"
	"time"

	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/ds"
)

// AttendFind finds attends for given user id and course id
func (c *DB) AttendFind(ctx context.Context, userID, courseID string) (*model.Attend, error) {
	var x model.Attend
	err := c.client.QueryFirst(ctx, kindAttend, &x,
		ds.Filter("UserID =", userID),
		ds.Filter("CourseID =", courseID),
		ds.CreateAfter(time.Now().Add(-6*time.Hour), true),
	)
	if ds.NotFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &x, nil
}

// AttendSave saves attend to database
func (c *DB) AttendSave(ctx context.Context, x *model.Attend) error {
	return c.client.Save(ctx, kindAttend, x)
}
