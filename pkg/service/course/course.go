package course

import (
	"context"

	"github.com/acoshift/acourse/pkg/app"
	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/acourse/pkg/store"
	"github.com/acoshift/httperror"
)

// New creates new course service
func New(store Store) app.CourseService {
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

func (s *service) ListCourses(ctx context.Context, req *app.CourseListRequest) (*app.CoursesReply, error) {
	currentUserID, _ := ctx.Value(app.KeyCurrentUserID).(string)

	var courses model.Courses
	var err error

	opt := make([]store.CourseListOption, 0, 3)

	if req.Owner != "" {
		opt = append(opt, store.CourseListOptionOwner(req.Owner))
		if req.Owner != currentUserID {
			opt = append(opt, store.CourseListOptionPublic(true))
		}
	} else if req.Public == nil || *req.Public == true {
		opt = append(opt, store.CourseListOptionPublic(true))
	} else {
		if currentUserID == "" {
			return nil, httperror.Unauthorized
		}
		// check is admin
		var role *model.Role
		role, err = s.store.RoleGet(currentUserID)
		if err != nil {
			return nil, err
		}
		if !role.Admin {
			return nil, httperror.Forbidden
		}
	}

	courses, err = s.store.CourseList(opt...)
	if err != nil {
		return nil, err
	}

	courses.SetView(model.CourseViewTiny)

	reply := new(app.CoursesReply)
	reply.Courses = courses

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
	users.SetView(model.UserViewTiny)
	reply.Users = users

	if req.EnrollCount {
		enrollCount := map[string]int{}
		for _, course := range courses {
			enrollCount[course.ID], err = s.store.EnrollCourseCount(course.ID)
			if err != nil {
				return nil, err
			}
		}
		reply.EnrollCount = enrollCount
	}

	return reply, nil
}

func (s *service) ListEnrolledCourses(ctx context.Context) (*app.CoursesReply, error) {
	currentUserID, ok := ctx.Value(app.KeyCurrentUserID).(string)
	if !ok {
		return nil, httperror.Unauthorized
	}

	enrolls, err := s.store.EnrollListByUserID(currentUserID)
	if err != nil {
		return nil, err
	}
	ids := make([]string, len(enrolls))
	for i, e := range enrolls {
		ids[i] = e.CourseID
	}
	courses, err := s.store.CourseGetAllByIDs(ids)
	courses.SetView(model.CourseViewTiny)

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
	users.SetView(model.UserViewTiny)

	enrollCount := map[string]int{}
	for _, course := range courses {
		enrollCount[course.ID], err = s.store.EnrollCourseCount(course.ID)
		if err != nil {
			return nil, err
		}
	}
	return &app.CoursesReply{
		Courses:     courses,
		Users:       users,
		EnrollCount: enrollCount,
	}, nil
}
