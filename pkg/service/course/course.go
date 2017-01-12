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
}

type service struct {
	store Store
}

func (s *service) ListCourses(ctx context.Context, req *app.CourseListRequest) (*app.CoursesReply, error) {
	var courses model.Courses
	var err error

	if req.Public == nil || *req.Public == true {
		// return public courses
		courses, err = s.store.CourseList(store.CourseListOptionPublic(true))
	} else {
		// check is admin
		currentUserID, ok := ctx.Value(app.KeyCurrentUserID).(string)
		if !ok {
			return nil, httperror.Unauthorized
		}
		var role *model.Role
		role, err = s.store.RoleGet(currentUserID)
		if err != nil {
			return nil, err
		}
		if !role.Admin {
			return nil, httperror.Forbidden
		}
		courses, err = s.store.CourseList()
	}

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

	if req.Student {
		student := map[string]int{}
		for _, course := range courses {
			student[course.ID], err = s.store.EnrollCourseCount(course.ID)
			if err != nil {
				return nil, err
			}
		}
		reply.Students = student
	}

	return reply, nil
}
