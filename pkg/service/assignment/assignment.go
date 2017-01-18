package assignment

import (
	"context"

	"github.com/acoshift/acourse/pkg/acourse"
	"github.com/acoshift/acourse/pkg/model"
	_context "golang.org/x/net/context"
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
	EnrollFind(context.Context, string, string) (*model.Enroll, error)
	UserAssignmentSave(context.Context, *model.UserAssignment) error
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

func (s *service) CreateAssignment(ctx _context.Context, req *acourse.Assignment) (*acourse.Assignment, error) {
	userID, ok := ctx.Value(acourse.KeyUserID).(string)
	if !ok || userID == "" {
		return nil, grpc.Errorf(codes.Unauthenticated, "authorization required")
	}
	if err := s.isCourseOwner(ctx, req.GetCourseId(), userID); err != nil {
		return nil, err
	}

	assignment := &model.Assignment{
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

func (s *service) UpdateAssignment(ctx _context.Context, req *acourse.Assignment) (*acourse.Empty, error) {
	userID, ok := ctx.Value(acourse.KeyUserID).(string)
	if !ok || userID == "" {
		return nil, grpc.Errorf(codes.Unauthenticated, "authorization required")
	}
	if err := s.isCourseOwner(ctx, req.GetCourseId(), userID); err != nil {
		return nil, err
	}

	assignment, err := s.store.AssignmentGet(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	if assignment == nil {
		return nil, grpc.Errorf(codes.NotFound, "assignment not found")
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

func (s *service) OpenAssignment(ctx _context.Context, req *acourse.AssignmentIDRequest) (*acourse.Empty, error) {
	return s.changeOpenAssignment(ctx, req, true)
}

func (s *service) CloseAssignment(ctx _context.Context, req *acourse.AssignmentIDRequest) (*acourse.Empty, error) {
	return s.changeOpenAssignment(ctx, req, false)
}

func (s *service) DeleteAssignment(ctx _context.Context, req *acourse.AssignmentIDRequest) (*acourse.Empty, error) {
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

func (s *service) SubmitUserAssignment(ctx _context.Context, req *acourse.UserAssignment) (*acourse.UserAssignment, error) {
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

	enroll, err := s.store.EnrollFind(ctx, userID, assignment.CourseID)
	if err != nil {
		return nil, err
	}
	if enroll == nil {
		return nil, grpc.Errorf(codes.PermissionDenied, "can not submit user assignment for this course")
	}

	// create model
	userAssignment := &model.UserAssignment{
		AssignmentID: assignment.ID,
		UserID:       userID,
		URL:          req.GetUrl(),
	}

	err = s.store.UserAssignmentSave(ctx, userAssignment)
	if err != nil {
		return nil, err
	}

	return acourse.ToUserAssignment(userAssignment), nil
}

func (s *service) DeleteUserAssignment(ctx _context.Context, req *acourse.UserAssignmentIDRequest) (*acourse.Empty, error) {
	// check is owner

	// delete model
	return nil, nil
}

func (s *service) GetUserAssignments(ctx _context.Context, req *acourse.AssignmentIDsRequest) (*acourse.UserAssignmentsResponse, error) {
	// check is owner or course owner

	// get models
	return nil, nil
}
