package course

import (
	"github.com/acoshift/acourse/pkg/acourse"
	"github.com/acoshift/ds"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

var (
	errCourseNotFound      = grpc.Errorf(codes.NotFound, "course: not found")
	errCourseURLExists     = grpc.Errorf(codes.AlreadyExists, "course: url already exists")
	errEnrollAlreadyExists = grpc.Errorf(codes.AlreadyExists, "course: enroll already exists")
	errEnrollNotFound      = grpc.Errorf(codes.NotFound, "course: enroll not found")
)

func (s *service) countEnroll(ctx context.Context, courseID string) (int, error) {
	// TODO: get from cache

	cnt, err := s.client.QueryCount(ctx, kindEnroll, ds.Filter("CourseID =", courseID))
	if err != nil {
		return 0, err
	}

	// TODO: save to cache

	return cnt, nil
}

func (s *service) FindEnroll(ctx context.Context, req *acourse.EnrollFindRequest) (*acourse.Enroll, error) {
	var x enrollModel

	err := s.client.QueryFirst(ctx, kindEnroll, &x,
		ds.Filter("UserID =", req.UserId),
		ds.Filter("CourseID =", req.CourseId),
	)

	if ds.NotFound(err) {
		return nil, errEnrollNotFound
	}
	if err != nil {
		return nil, err
	}
	return toEnroll(&x), nil
}

func (s *service) listEnrollByUserID(ctx context.Context, userID string) ([]*enrollModel, error) {
	var xs []*enrollModel

	err := s.client.Query(ctx, kindEnroll, &xs,
		ds.Filter("UserID =", userID),
	)
	if err != nil {
		return nil, err
	}

	return xs, nil
}

func (s *service) saveEnroll(ctx context.Context, x *enrollModel) error {
	// TODO: race condition
	// TODO: use keysonly query
	var t enrollModel
	err := s.client.QueryFirst(ctx, kindEnroll, &t,
		ds.Filter("UserID =", x.UserID),
		ds.Filter("CourseID =", x.CourseID),
	)
	if err == nil {
		return errEnrollAlreadyExists
	}

	err = s.client.AllocateModel(ctx, kindEnroll, x)
	if err != nil {
		return err
	}

	return nil
}

func toEnroll(x *enrollModel) *acourse.Enroll {
	return &acourse.Enroll{
		UserID:   x.UserID,
		CourseID: x.CourseID,
	}
}

func fromEnroll(x *acourse.Enroll) *enrollModel {
	return &enrollModel{
		UserID:   x.UserID,
		CourseID: x.CourseID,
	}
}

func fromEnrolls(xs []*acourse.Enroll) []*enrollModel {
	rs := make([]*enrollModel, len(xs))
	for i, x := range xs {
		rs[i] = fromEnroll(x)
	}
	return rs
}

func (s *service) CreateEnrolls(ctx context.Context, req *acourse.EnrollsRequest) (*acourse.Empty, error) {
	enrolls := fromEnrolls(req.Enrolls)

	var err error
	for _, enroll := range enrolls {
		var t enrollModel
		err = s.client.QueryFirst(ctx, kindEnroll, &t,
			ds.Filter("UserID =", enroll.UserID),
			ds.Filter("CourseID =", enroll.CourseID),
		)
		if t.Key != nil {
			return nil, errEnrollAlreadyExists
		}
	}

	err = s.client.AllocateModels(ctx, kindEnroll, enrolls)
	if err != nil {
		return nil, err
	}

	return new(acourse.Empty), nil
}
