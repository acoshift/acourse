package assignment

import (
	"github.com/acoshift/acourse/pkg/acourse"
	"github.com/acoshift/acourse/pkg/internal"
	"github.com/acoshift/ds"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// New creates new assignment service server
func New(client *ds.Client, course acourse.CourseServiceClient) acourse.AssignmentServiceServer {
	return &service{client, course}
}

type service struct {
	client *ds.Client
	course acourse.CourseServiceClient
}

func (s *service) isCourseOwner(ctx context.Context, courseID, userID string) error {
	course, err := s.course.GetCourse(ctx, &acourse.CourseIDRequest{CourseId: courseID})
	if err != nil {
		return err
	}
	if course.GetCourse().GetOwner() != userID {
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

	_, err = s.course.FindEnroll(ctx, &acourse.EnrollFindRequest{UserId: userID, CourseId: x.CourseID})
	if grpc.Code(err) == codes.NotFound {
		return nil, grpc.Errorf(codes.PermissionDenied, "can not submit user assignment for this course")
	}
	if err != nil {
		return nil, err
	}

	// create model
	userAssignment := &userAssignment{
		AssignmentID: x.ID(),
		UserID:       userID,
		URL:          req.GetUrl(),
	}

	err = s.client.SaveModel(ctx, kindUserAssignment, userAssignment)
	if err != nil {
		return nil, err
	}

	return toUserAssignment(userAssignment), nil
}

func (s *service) DeleteUserAssignment(ctx context.Context, req *acourse.UserAssignmentIDRequest) (*acourse.Empty, error) {
	userID := internal.GetUserID(ctx)
	if userID == "" {
		return nil, grpc.Errorf(codes.Unauthenticated, "authorization required")
	}

	var x userAssignment
	err := s.client.GetByStringID(ctx, kindUserAssignment, req.GetUserAssignmentId(), &x)
	err = ds.IgnoreFieldMismatch(err)
	if ds.NotFound(err) {
		return nil, ErrUserAssignmentNotFound
	}
	if err != nil {
		return nil, err
	}

	if x.UserID != userID {
		return nil, ErrPermissionDenied
	}

	err = s.client.DeleteByStringID(ctx, kindUserAssignment, req.GetUserAssignmentId())
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
	xs := make([]*userAssignment, 0)
	for _, assignment := range assignments {
		// TODO: need refactor
		var tmp []*userAssignment
		err := s.client.Query(ctx, kindUserAssignment, &tmp,
			ds.Filter("AssignmentID =", assignment.ID()),
			ds.Filter("UserID =", userID),
		)
		err = ds.IgnoreFieldMismatch(err)
		if err != nil {
			return nil, err
		}
		xs = append(xs, tmp...)
	}

	return &acourse.UserAssignmentsResponse{
		UserAssignments: toUserAssignments(xs),
	}, nil
}

func (s *service) ListUserAssignments(ctx context.Context, req *acourse.CourseIDRequest) (*acourse.UserAssignmentsResponse, error) {
	userID := internal.GetUserID(ctx)
	if userID == "" {
		return nil, grpc.Errorf(codes.Unauthenticated, "authorization required")
	}

	course, err := s.course.GetCourse(ctx, &acourse.CourseIDRequest{CourseId: req.GetCourseId()})
	if err != nil {
		return nil, err
	}

	if course.GetCourse().GetOwner() != userID {
		return nil, grpc.Errorf(codes.PermissionDenied, "only instructor can list user assignments")
	}

	var assignments []*assignment
	err = s.client.Query(ctx, kindAssignment, &assignments, ds.Filter("CourseID =", course.GetCourse().GetId()))
	err = ds.IgnoreFieldMismatch(err)
	err = ds.IgnoreNotFound(err)
	if err != nil {
		return nil, err
	}

	xs := make([]*userAssignment, 0)

	for _, assignment := range assignments {
		var tmp []*userAssignment
		err := s.client.Query(ctx, kindUserAssignment, &tmp,
			ds.Filter("AssignmentID =", assignment.ID()),
		)
		err = ds.IgnoreFieldMismatch(err)
		if err != nil {
			return nil, err
		}
		xs = append(xs, tmp...)
	}

	return &acourse.UserAssignmentsResponse{
		UserAssignments: toUserAssignments(xs),
	}, nil
}
