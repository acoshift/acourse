package course

import (
	"time"

	"github.com/acoshift/ds"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

var (
	errAttendNotFound = grpc.Errorf(codes.NotFound, "course: attend not found")
)

func (s *service) findAttend(ctx context.Context, userID, courseID string) (*attendModel, error) {
	var x attendModel
	err := s.client.QueryFirst(ctx, kindAttend, &x,
		ds.Filter("UserID =", userID),
		ds.Filter("CourseID =", courseID),
		ds.CreateAfter(time.Now().Add(-8*time.Hour), true),
	)
	err = ds.IgnoreFieldMismatch(err)
	if ds.NotFound(err) {
		return nil, errAttendNotFound
	}
	if err != nil {
		return nil, err
	}
	return &x, nil
}

func (s *service) saveAttend(ctx context.Context, x *attendModel) error {
	return s.client.SaveModel(ctx, x)
}
