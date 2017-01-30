package course

import (
	"context"
	"time"

	"github.com/acoshift/ds"
)

func (s *service) findAttend(ctx context.Context, userID, courseID string) (*attendModel, error) {
	var x attendModel
	err := s.client.QueryFirst(ctx, kindAttend, &x,
		ds.Filter("UserID =", userID),
		ds.Filter("CourseID =", courseID),
		ds.CreateAfter(time.Now().Add(-6*time.Hour), true),
	)
	err = ds.IgnoreFieldMismatch(err)
	if ds.NotFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &x, nil
}

func (s *service) saveAttend(ctx context.Context, x *attendModel) error {
	return s.client.SaveModel(ctx, kindAttend, x)
}
