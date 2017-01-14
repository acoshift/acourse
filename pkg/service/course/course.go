package course

import (
	"context"

	"github.com/acoshift/acourse/pkg/acourse"
	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/acourse/pkg/store"
	_context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// New creates new course service
func New(store Store) acourse.CourseServiceServer {
	return &service{store}
}

// Store is the store interface for course service
type Store interface {
	CourseList(opts ...store.CourseListOption) (model.Courses, error)
	UserGetMulti(context.Context, []string) (model.Users, error)
	EnrollCourseCount(string) (int, error)
	RoleGet(string) (*model.Role, error)
	EnrollListByUserID(string) (model.Enrolls, error)
	CourseGetAllByIDs([]string) (model.Courses, error)
}

type service struct {
	store Store
}

func (s *service) listCourses(ctx _context.Context, opts ...store.CourseListOption) (*acourse.CoursesResponse, error) {
	courses, err := s.store.CourseList(opts...)
	if err != nil {
		return nil, err
	}
	// get owners
	userIDMap := map[string]bool{}
	for _, course := range courses {
		userIDMap[course.Owner] = true
	}
	userIDs := make([]string, 0, len(userIDMap))
	for id := range userIDMap {
		userIDs = append(userIDs, id)
	}
	users, err := s.store.UserGetMulti(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	enrollCounts := make([]*acourse.EnrollCount, len(courses))
	for i, course := range courses {
		c, err := s.store.EnrollCourseCount(course.ID)
		if err != nil {
			return nil, err
		}
		enrollCounts[i] = &acourse.EnrollCount{
			CourseId: course.ID,
			Count:    int32(c),
		}
	}
	return &acourse.CoursesResponse{
		Courses:      acourse.ToCoursesSmall(courses),
		Users:        acourse.ToUsersTiny(users),
		EnrollCounts: enrollCounts,
	}, nil
}

func (s *service) ListPublicCourses(ctx _context.Context, req *acourse.ListRequest) (*acourse.CoursesResponse, error) {
	return s.listCourses(ctx, store.CourseListOptionPublic(true))
}

func (s *service) ListCourses(ctx _context.Context, req *acourse.ListRequest) (*acourse.CoursesResponse, error) {
	userID, ok := ctx.Value(acourse.KeyUserID).(string)
	if !ok || userID == "" {
		return nil, grpc.Errorf(codes.Unauthenticated, "authorization required")
	}
	role, err := s.store.RoleGet(userID)
	if err != nil {
		return nil, err
	}
	if !role.Admin {
		return nil, grpc.Errorf(codes.PermissionDenied, "permission denied")
	}
	return s.listCourses(ctx)
}

func (s *service) ListOwnCourses(ctx _context.Context, req *acourse.UserIDRequest) (*acourse.CoursesResponse, error) {
	userID, _ := ctx.Value(acourse.KeyUserID).(string)

	if req.GetUserId() == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "invalid user id")
	}

	opt := make([]store.CourseListOption, 0, 3)
	opt = append(opt, store.CourseListOptionOwner(req.GetUserId()))

	// if not sign in, get only public courses
	if userID == "" {
		opt = append(opt, store.CourseListOptionPublic(true))
	}

	return s.listCourses(ctx, opt...)
}

func (s *service) ListEnrolledCourses(ctx _context.Context, req *acourse.UserIDRequest) (*acourse.CoursesResponse, error) {
	userID, ok := ctx.Value(acourse.KeyUserID).(string)
	if !ok || userID == "" {
		return nil, grpc.Errorf(codes.Unauthenticated, "authorization required")
	}

	if req.GetUserId() == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "invalid user id")
	}

	// only admin allow for get other user enrolled courses
	if req.GetUserId() != userID {
		role, err := s.store.RoleGet(userID)
		if err != nil {
			return nil, err
		}
		if !role.Admin {
			return nil, grpc.Errorf(codes.PermissionDenied, "permission denied")
		}
	}

	enrolls, err := s.store.EnrollListByUserID(req.GetUserId())
	if err != nil {
		return nil, err
	}
	ids := make([]string, len(enrolls))
	for i, e := range enrolls {
		ids[i] = e.CourseID
	}
	courses, err := s.store.CourseGetAllByIDs(ids)

	// get owners
	userIDMap := map[string]bool{}
	for _, course := range courses {
		userIDMap[course.Owner] = true
	}
	userIDs := make([]string, 0, len(userIDMap))
	for id := range userIDMap {
		userIDs = append(userIDs, id)
	}
	users, err := s.store.UserGetMulti(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	enrollCounts := make([]*acourse.EnrollCount, len(courses))
	for i, course := range courses {
		c, err := s.store.EnrollCourseCount(course.ID)
		if err != nil {
			return nil, err
		}
		enrollCounts[i] = &acourse.EnrollCount{
			CourseId: course.ID,
			Count:    int32(c),
		}
	}
	return &acourse.CoursesResponse{
		Courses:      acourse.ToCoursesSmall(courses),
		Users:        acourse.ToUsersTiny(users),
		EnrollCounts: enrollCounts,
	}, nil
}
