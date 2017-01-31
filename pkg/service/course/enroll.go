package course

import (
	"context"

	"cloud.google.com/go/datastore"
	"github.com/acoshift/acourse/pkg/acourse"
	"github.com/acoshift/ds"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

var (
	errCourseNotFound      = grpc.Errorf(codes.NotFound, "course: not found")
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

	err = s.client.SaveModel(ctx, kindEnroll, x)
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

	// TODO: change to ds
	keys := make([]*datastore.Key, 0, len(enrolls))
	for _, enroll := range enrolls {
		enroll.Stamp()
		keys = append(keys, datastore.IncompleteKey(kindEnroll, nil))
	}

	var pKeys []*datastore.PendingKey

	commit, err := s.client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		var t enrollModel
		var err error
		for _, enroll := range enrolls {
			err = s.client.QueryFirst(ctx, kindEnroll, &t,
				ds.Filter("UserID =", enroll.UserID),
				ds.Filter("CourseID =", enroll.CourseID),
				ds.Transaction(tx),
			)
			if err == nil {
				return errEnrollAlreadyExists
			}
		}

		pKeys, err = tx.PutMulti(keys, enrolls)
		return err
	})
	if err != nil {
		return nil, err
	}

	ds.SetCommitKeys(commit, pKeys, enrolls)
	return new(acourse.Empty), nil
}
