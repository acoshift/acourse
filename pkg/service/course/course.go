package course

import (
	"context"

	"github.com/acoshift/acourse/pkg/app"
	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/acourse/pkg/store"
)

// New creates new course service
func New(store Store) app.CourseService {
	return &service{store}
}

// Store is the store interface for course service
type Store interface {
	CourseList(opts ...store.CourseListOption) (model.Courses, error)
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
	}

	if err != nil {
		return nil, err
	}

	courses.SetView(model.CourseViewTiny)

	return &app.CoursesReply{
		Courses: courses,
	}, nil
}
