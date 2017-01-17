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
	AssignmentSave(context.Context, *model.Assignment) error
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
	// check is course owner

	// get model

	// update model

	// save model
	return nil, nil
}

func (s *service) OpenAssignment(ctx _context.Context, req *acourse.AssignmentIDRequest) (*acourse.Empty, error) {
	// check is course owner

	// get model

	// update model

	// save model
	return nil, nil
}

func (s *service) CloseAssignment(ctx _context.Context, req *acourse.AssignmentIDRequest) (*acourse.Empty, error) {
	// check is course owner

	// get model

	// update model

	// save model
	return nil, nil
}

func (s *service) DeleteAssignment(ctx _context.Context, req *acourse.AssignmentIDRequest) (*acourse.Empty, error) {
	// check is course owner

	// delete model
	return nil, nil
}

func (s *service) SubmitUserAssignment(ctx _context.Context, req *acourse.UserAssignment) (*acourse.UserAssignment, error) {
	// check is enrolled

	// create model

	// save model
	return nil, nil
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
