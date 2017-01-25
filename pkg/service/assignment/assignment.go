package assignment

import (
	"context"

	"github.com/acoshift/acourse/pkg/acourse"
	"github.com/acoshift/acourse/pkg/internal"
	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/ds"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// New creates new assignment service server
func New(store Store, client *ds.Client) acourse.AssignmentServiceServer {
	return &service{store, client}
}

// Store is the store interface for assignment service
type Store interface {
	CourseGet(context.Context, string) (*model.Course, error)
	EnrollFind(context.Context, string, string) (*model.Enroll, error)
	UserAssignmentSave(context.Context, *model.UserAssignment) error
	UserAssignmentGet(context.Context, string) (*model.UserAssignment, error)
	UserAssignmentDelete(context.Context, string) error
	UserAssignmentList(context.Context, string, string) (model.UserAssignments, error)
}

type service struct {
	store  Store
	client *ds.Client
}

func (s *service) isCourseOwner(ctx context.Context, courseID, userID string) error {
	course, err := s.store.CourseGet(ctx, courseID)
	if err != nil {
		return err
	}
	if course == nil {
		return grpc.Errorf(codes.NotFound, "course not found")
	}
	if course.Owner != userID {
		return grpc.Errorf(codes.PermissionDenied, "only course owner can create assignment")
	}
	return nil
}

func (s *service) CreateAssignment(ctx context.Context, req *acourse.Assignment) (*acourse.Assignment, error) {
	userID := internal.GetUserID(ctx)
	if userID == "" {
		return nil, grpc.Errorf(codes.Unauthenticated, "authorization required")
	}
	if err := s.isCourseOwner(ctx, req.GetCourseId(), userID); err != nil {
		return nil, err
	}

	x := assignment{
		CourseID:    req.GetCourseId(),
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		Open:        req.GetOpen(),
	}

	// save model
	err := s.client.SaveModel(ctx, kindAssignment, &x)
	if err != nil {
		return nil, err
	}
	return toAssignment(&x), nil
}

func (s *service) UpdateAssignment(ctx context.Context, req *acourse.Assignment) (*acourse.Empty, error) {
	userID := internal.GetUserID(ctx)
	if userID == "" {
		return nil, grpc.Errorf(codes.Unauthenticated, "authorization required")
	}

	var x assignment
	err := s.client.GetByStringID(ctx, kindAssignment, req.GetId(), &x)
	err = ds.IgnoreFieldMismatch(err)
	if ds.NotFound(err) {
		return nil, ErrAssignmentNotFound
	}
	if err != nil {
		return nil, err
	}

	err = s.isCourseOwner(ctx, x.CourseID, userID)
	if err != nil {
		return nil, err
	}

	x.Title = req.GetTitle()
	x.Description = req.GetDescription()

	err = s.client.SaveModel(ctx, "", &x)
	if err != nil {
		return nil, err
	}
	return new(acourse.Empty), nil
}

func (s *service) changeOpenAssignment(ctx context.Context, req *acourse.AssignmentIDRequest, open bool) (*acourse.Empty, error) {
	userID := internal.GetUserID(ctx)
	if userID == "" {
		return nil, grpc.Errorf(codes.Unauthenticated, "authorization required")
	}

	var x assignment
	err := s.client.GetByStringID(ctx, kindAssignment, req.GetAssignmentId(), &x)
	err = ds.IgnoreFieldMismatch(err)
	if ds.NotFound(err) {
		return nil, ErrAssignmentNotFound
	}
	if err != nil {
		return nil, err
	}

	err = s.isCourseOwner(ctx, x.CourseID, userID)
	if err != nil {
		return nil, err
	}

	x.Open = open

	err = s.client.SaveModel(ctx, "", &x)
	if err != nil {
		return nil, err
	}
	return new(acourse.Empty), nil
}

func (s *service) OpenAssignment(ctx context.Context, req *acourse.AssignmentIDRequest) (*acourse.Empty, error) {
	return s.changeOpenAssignment(ctx, req, true)
}

func (s *service) CloseAssignment(ctx context.Context, req *acourse.AssignmentIDRequest) (*acourse.Empty, error) {
	return s.changeOpenAssignment(ctx, req, false)
}

func (s *service) ListAssignments(ctx context.Context, req *acourse.CourseIDRequest) (*acourse.AssignmentsResponse, error) {
	// TODO: check is owner or enrolled

	var xs []*assignment
	err := s.client.Query(ctx, kindAssignment, &xs, ds.Filter("CourseID =", req.GetCourseId()))
	err = ds.IgnoreFieldMismatch(err)
	err = ds.IgnoreNotFound(err)
	if err != nil {
		return nil, err
	}
	return &acourse.AssignmentsResponse{
		Assignments: toAssignments(xs),
	}, nil
}

func (s *service) DeleteAssignment(ctx context.Context, req *acourse.AssignmentIDRequest) (*acourse.Empty, error) {
	userID := internal.GetUserID(ctx)
	if userID == "" {
		return nil, grpc.Errorf(codes.Unauthenticated, "authorization required")
	}

	var x assignment
	err := s.client.GetByStringID(ctx, kindAssignment, req.GetAssignmentId(), &x)
	err = ds.IgnoreFieldMismatch(err)
	if ds.NotFound(err) {
		return nil, ErrAssignmentNotFound
	}
	if err != nil {
		return nil, err
	}

	err = s.isCourseOwner(ctx, x.CourseID, userID)
	if err != nil {
		return nil, err
	}

	err = s.client.DeleteByStringID(ctx, kindAssignment, req.GetAssignmentId())
	if err != nil {
		return nil, err
	}
	return new(acourse.Empty), nil
}

func (s *service) SubmitUserAssignment(ctx context.Context, req *acourse.UserAssignment) (*acourse.UserAssignment, error) {
	userID := internal.GetUserID(ctx)
	if userID == "" {
		return nil, grpc.Errorf(codes.Unauthenticated, "authorization required")
	}

	var x assignment
	err := s.client.GetByStringID(ctx, kindAssignment, req.GetAssignmentId(), &x)
	err = ds.IgnoreFieldMismatch(err)
	if ds.NotFound(err) {
		return nil, ErrAssignmentNotFound
	}
	if err != nil {
		return nil, err
	}
	if !x.Open {
		return nil, ErrAssignmentNotOpen
	}

	enroll, err := s.store.EnrollFind(ctx, userID, x.CourseID)
	if err != nil {
		return nil, err
	}
	if enroll == nil {
		return nil, grpc.Errorf(codes.PermissionDenied, "can not submit user assignment for this course")
	}

	// create model
	userAssignment := &model.UserAssignment{
		AssignmentID: x.ID(),
		UserID:       userID,
		URL:          req.GetUrl(),
	}

	err = s.store.UserAssignmentSave(ctx, userAssignment)
	if err != nil {
		return nil, err
	}

	return acourse.ToUserAssignment(userAssignment), nil
}

func (s *service) DeleteUserAssignment(ctx context.Context, req *acourse.UserAssignmentIDRequest) (*acourse.Empty, error) {
	userID := internal.GetUserID(ctx)
	if userID == "" {
		return nil, grpc.Errorf(codes.Unauthenticated, "authorization required")
	}

	userAssignment, err := s.store.UserAssignmentGet(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	if userAssignment == nil {
		return nil, grpc.Errorf(codes.NotFound, "user assignment not found")
	}

	if userAssignment.UserID != userID {
		return nil, grpc.Errorf(codes.PermissionDenied, "can not delete this user assignment")
	}

	err = s.store.UserAssignmentDelete(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return new(acourse.Empty), nil
}

func (s *service) GetUserAssignments(ctx context.Context, req *acourse.AssignmentIDsRequest) (*acourse.UserAssignmentsResponse, error) {
	userID := internal.GetUserID(ctx)
	if userID == "" {
		return nil, grpc.Errorf(codes.Unauthenticated, "authorization required")
	}

	var assignments []*assignment
	err := s.client.GetByStringIDs(ctx, kindAssignment, req.GetAssignmentIds(), &assignments)
	err = ds.IgnoreFieldMismatch(err)
	if err != nil {
		return nil, err
	}
	userAssignments := make(model.UserAssignments, 0)
	for _, assignment := range assignments {
		// TODO: need refactor
		tmp, err := s.store.UserAssignmentList(ctx, assignment.ID(), userID)
		if err != nil {
			return nil, err
		}
		userAssignments = append(userAssignments, tmp...)
	}

	return &acourse.UserAssignmentsResponse{
		UserAssignments: acourse.ToUserAssignments(userAssignments),
	}, nil
}

func (s *service) ListUserAssignments(ctx context.Context, req *acourse.CourseIDRequest) (*acourse.UserAssignmentsResponse, error) {
	userID := internal.GetUserID(ctx)
	if userID == "" {
		return nil, grpc.Errorf(codes.Unauthenticated, "authorization required")
	}

	course, err := s.store.CourseGet(ctx, req.GetCourseId())
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, grpc.Errorf(codes.NotFound, "course not found")
	}

	if course.Owner != userID {
		return nil, grpc.Errorf(codes.PermissionDenied, "only instructor can list user assignments")
	}

	var xs []*assignment
	err = s.client.Query(ctx, kindAssignment, ds.Filter("CourseID =", course.ID()))
	err = ds.IgnoreFieldMismatch(err)
	err = ds.IgnoreNotFound(err)
	if err != nil {
		return nil, err
	}

	userAssignments := model.UserAssignments{}

	for _, x := range xs {
		temp, err := s.store.UserAssignmentList(ctx, x.ID(), "")
		if err != nil {
			return nil, err
		}
		userAssignments = append(userAssignments, temp...)
	}

	return &acourse.UserAssignmentsResponse{
		UserAssignments: acourse.ToUserAssignments(userAssignments),
	}, nil
}
