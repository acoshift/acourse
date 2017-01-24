package assignment

import (
	"context"

	"github.com/acoshift/acourse/pkg/acourse"
	"github.com/acoshift/acourse/pkg/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// New creates new assignment service server
func New(store Store) acourse.AssignmentServiceServer {
	return &service{store}
}

// Store is the store interface for assignment service
type Store interface {
	CourseGet(context.Context, string) (*model.Course, error)
	AssignmentGet(context.Context, string) (*model.Assignment, error)
	AssignmentSave(context.Context, *model.Assignment) error
	AssignmentList(context.Context, string) (model.Assignments, error)
	AssignmentDelete(context.Context, string) error
	AssignmentGetMulti(context.Context, []string) (model.Assignments, error)
	EnrollFind(context.Context, string, string) (*model.Enroll, error)
	UserAssignmentSave(context.Context, *model.UserAssignment) error
	UserAssignmentGet(context.Context, string) (*model.UserAssignment, error)
	UserAssignmentDelete(context.Context, string) error
	UserAssignmentList(context.Context, string, string) (model.UserAssignments, error)
}

type service struct {
	store Store
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
	userID, ok := ctx.Value(acourse.KeyUserID).(string)
	if !ok || userID == "" {
		return nil, grpc.Errorf(codes.Unauthenticated, "authorization required")
	}
	if err := s.isCourseOwner(ctx, req.GetCourseId(), userID); err != nil {
		return nil, err
	}

	assignment := &model.Assignment{
		CourseID:    req.GetCourseId(),
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		Open:        req.GetOpen(),
	}

	// save model
	err := s.store.AssignmentSave(ctx, assignment)
	if err != nil {
		return nil, err
	}
	return acourse.ToAssignment(assignment), nil
}

func (s *service) UpdateAssignment(ctx context.Context, req *acourse.Assignment) (*acourse.Empty, error) {
	userID, ok := ctx.Value(acourse.KeyUserID).(string)
	if !ok || userID == "" {
		return nil, grpc.Errorf(codes.Unauthenticated, "authorization required")
	}
	assignment, err := s.store.AssignmentGet(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	if assignment == nil {
		return nil, grpc.Errorf(codes.NotFound, "assignment not found")
	}

	if err := s.isCourseOwner(ctx, assignment.CourseID, userID); err != nil {
		return nil, err
	}

	assignment.Title = req.GetTitle()
	assignment.Description = req.GetDescription()

	err = s.store.AssignmentSave(ctx, assignment)
	if err != nil {
		return nil, err
	}
	return new(acourse.Empty), nil
}

func (s *service) changeOpenAssignment(ctx context.Context, req *acourse.AssignmentIDRequest, open bool) (*acourse.Empty, error) {
	userID, ok := ctx.Value(acourse.KeyUserID).(string)
	if !ok || userID == "" {
		return nil, grpc.Errorf(codes.Unauthenticated, "authorization required")
	}

	assignment, err := s.store.AssignmentGet(ctx, req.GetAssignmentId())
	if err != nil {
		return nil, err
	}
	if assignment == nil {
		return nil, grpc.Errorf(codes.NotFound, "assignment not found")
	}

	if err := s.isCourseOwner(ctx, assignment.CourseID, userID); err != nil {
		return nil, err
	}

	assignment.Open = open

	err = s.store.AssignmentSave(ctx, assignment)
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

	assignments, err := s.store.AssignmentList(ctx, req.GetCourseId())
	if err != nil {
		return nil, err
	}
	return &acourse.AssignmentsResponse{
		Assignments: acourse.ToAssignments(assignments),
	}, nil
}

func (s *service) DeleteAssignment(ctx context.Context, req *acourse.AssignmentIDRequest) (*acourse.Empty, error) {
	userID, ok := ctx.Value(acourse.KeyUserID).(string)
	if !ok || userID == "" {
		return nil, grpc.Errorf(codes.Unauthenticated, "authorization required")
	}

	assignment, err := s.store.AssignmentGet(ctx, req.GetAssignmentId())
	if err != nil {
		return nil, err
	}
	if assignment == nil {
		return nil, grpc.Errorf(codes.NotFound, "assignment not found")
	}

	if err := s.isCourseOwner(ctx, assignment.CourseID, userID); err != nil {
		return nil, err
	}

	err = s.store.AssignmentDelete(ctx, req.GetAssignmentId())
	if err != nil {
		return nil, err
	}
	return new(acourse.Empty), nil
}

func (s *service) SubmitUserAssignment(ctx context.Context, req *acourse.UserAssignment) (*acourse.UserAssignment, error) {
	userID, ok := ctx.Value(acourse.KeyUserID).(string)
	if !ok || userID == "" {
		return nil, grpc.Errorf(codes.Unauthenticated, "authorization required")
	}

	assignment, err := s.store.AssignmentGet(ctx, req.GetAssignmentId())
	if err != nil {
		return nil, err
	}
	if assignment == nil {
		return nil, grpc.Errorf(codes.NotFound, "assignment not found")
	}
	if !assignment.Open {
		return nil, grpc.Errorf(codes.PermissionDenied, "assignment not open")
	}

	enroll, err := s.store.EnrollFind(ctx, userID, assignment.CourseID)
	if err != nil {
		return nil, err
	}
	if enroll == nil {
		return nil, grpc.Errorf(codes.PermissionDenied, "can not submit user assignment for this course")
	}

	// create model
	userAssignment := &model.UserAssignment{
		AssignmentID: assignment.ID(),
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
	userID, ok := ctx.Value(acourse.KeyUserID).(string)
	if !ok || userID == "" {
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
	userID, ok := ctx.Value(acourse.KeyUserID).(string)
	if !ok || userID == "" {
		return nil, grpc.Errorf(codes.Unauthenticated, "authorization required")
	}

	assignments, err := s.store.AssignmentGetMulti(ctx, req.GetAssignmentIds())
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
	userID, ok := ctx.Value(acourse.KeyUserID).(string)
	if !ok || userID == "" {
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

	assignments, err := s.store.AssignmentList(ctx, course.ID())
	if err != nil {
		return nil, err
	}

	userAssignments := model.UserAssignments{}

	for _, assignment := range assignments {
		temp, err := s.store.UserAssignmentList(ctx, assignment.ID(), "")
		if err != nil {
			return nil, err
		}
		userAssignments = append(userAssignments, temp...)
	}

	return &acourse.UserAssignmentsResponse{
		UserAssignments: acourse.ToUserAssignments(userAssignments),
	}, nil
}
