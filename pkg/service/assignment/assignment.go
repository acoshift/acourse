package assignment

import (
	"github.com/acoshift/acourse/pkg/acourse"
	_context "golang.org/x/net/context"
)

// New creates new assignment service server
func New(store Store) acourse.AssignmentServiceServer {
	return &service{store}
}

// Store is the store interface for assignment service
type Store interface {
}

type service struct {
	store Store
}

func (s *service) CreateAssignment(ctx _context.Context, req *acourse.Assignment) (*acourse.Assignment, error) {
	return nil, nil
}

func (s *service) UpdateAssignment(ctx _context.Context, req *acourse.Assignment) (*acourse.Empty, error) {
	return nil, nil
}

func (s *service) OpenAssignment(ctx _context.Context, req *acourse.AssignmentIDRequest) (*acourse.Empty, error) {
	return nil, nil
}

func (s *service) CloseAssignment(ctx _context.Context, req *acourse.AssignmentIDRequest) (*acourse.Empty, error) {
	return nil, nil
}

func (s *service) DeleteAssignment(ctx _context.Context, req *acourse.AssignmentIDRequest) (*acourse.Empty, error) {
	return nil, nil
}

func (s *service) SubmitUserAssignment(ctx _context.Context, req *acourse.UserAssignment) (*acourse.UserAssignment, error) {
	return nil, nil
}

func (s *service) DeleteUserAssignment(ctx _context.Context, req *acourse.UserAssignmentIDRequest) (*acourse.Empty, error) {
	return nil, nil
}

func (s *service) GetUserAssignments(ctx _context.Context, req *acourse.AssignmentIDsRequest) (*acourse.UserAssignmentsResponse, error) {
	return nil, nil
}
