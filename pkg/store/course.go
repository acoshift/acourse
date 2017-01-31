package store

import (
	"context"
	"time"

	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/ds"
	"github.com/acoshift/gotcha"
)

var cacheCourse = gotcha.New()
var cacheCourseURL = gotcha.New()

// CourseGet retrieves course from database
func (c *DB) CourseGet(ctx context.Context, courseID string) (*model.Course, error) {
	if cache := cacheCourse.Get(courseID); cache != nil {
		return cache.(*model.Course), nil
	}

	var x model.Course
	err := c.client.GetByStringID(ctx, kindCourse, courseID, &x)
	err = ds.IgnoreFieldMismatch(err)
	if ds.NotFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	cacheCourse.SetTTL(courseID, &x, time.Minute)
	return &x, nil
}

// CourseGetMulti retrieves multiple courses from database
func (c *DB) CourseGetMulti(ctx context.Context, courseIDs []string) (model.Courses, error) {
	if len(courseIDs) == 0 {
		return []*model.Course{}, nil
	}

	courses := make([]*model.Course, 0, len(courseIDs))
	ids := make([]string, 0, len(courseIDs))

	for _, id := range courseIDs {
		if c := cacheCourse.Get(id); c != nil {
			courses = append(courses, c.(*model.Course))
		} else {
			ids = append(ids, id)
		}
	}
	if len(ids) == 0 {
		return courses, nil
	}

	xs := make([]*model.Course, len(ids))
	err := c.client.GetByStringIDs(ctx, kindCourse, ids, xs)
	err = ds.IgnoreFieldMismatch(err)
	if err != nil {
		return nil, err
	}
	for _, x := range xs {
		if x == nil {
			continue
		}
		courses = append(courses, x)
		cacheCourse.SetTTL(x.ID(), x, time.Minute)
	}
	return courses, nil
}
