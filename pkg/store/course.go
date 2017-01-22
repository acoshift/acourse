package store

import (
	"context"
	"time"

	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/ds"
	"github.com/acoshift/gotcha"
)

// CourseListOptions type
type CourseListOptions struct {
	Offset *int
	Limit  *int
	Public *bool
	Owner  *string
	Start  struct {
		From *time.Time
		To   *time.Time
	}
}

// CourseListOption type
type CourseListOption func(*CourseListOptions)

var cacheCourse = gotcha.New()
var cacheCourseURL = gotcha.New()

/* CourseListOptions */

// CourseListOptionOffset sets offset to options
func CourseListOptionOffset(offset int) CourseListOption {
	return func(args *CourseListOptions) {
		args.Offset = &offset
	}
}

// CourseListOptionLimit sets limit to options
func CourseListOptionLimit(limit int) CourseListOption {
	return func(args *CourseListOptions) {
		args.Limit = &limit
	}
}

// CourseListOptionPublic sets open to options
func CourseListOptionPublic(public bool) CourseListOption {
	return func(args *CourseListOptions) {
		args.Public = &public
	}
}

// CourseListOptionOwner sets owner to options
func CourseListOptionOwner(owner string) CourseListOption {
	return func(args *CourseListOptions) {
		args.Owner = &owner
	}
}

// CourseListOptionStartFrom sets start from to options
func CourseListOptionStartFrom(from time.Time) CourseListOption {
	return func(args *CourseListOptions) {
		args.Start.From = &from
	}
}

// CourseListOptionStartTo sets start to to options
func CourseListOptionStartTo(to time.Time) CourseListOption {
	return func(args *CourseListOptions) {
		args.Start.To = &to
	}
}

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

// CourseSave saves course to database
func (c *DB) CourseSave(ctx context.Context, x *model.Course) error {
	// Check duplicate URL
	if x.URL != "" {
		if t, err := c.CourseFind(ctx, x.URL); (t != nil && t.ID() != x.ID()) || err != nil {
			if err != nil {
				return err
			}
			return ErrConflict("course url already exists")
		}
	}

	err := c.client.Save(ctx, kindCourse, x)
	if err != nil {
		return err
	}

	cacheCourse.Unset(x.ID())
	return nil
}

// CourseList retrieves courses from database
func (c *DB) CourseList(ctx context.Context, opts ...CourseListOption) (model.Courses, error) {
	var xs []*model.Course

	qs := []ds.Query{}

	opt := &CourseListOptions{}
	for _, setter := range opts {
		setter(opt)
	}

	if opt.Offset != nil {
		qs = append(qs, ds.Offset(*opt.Offset))
	}
	if opt.Limit != nil {
		qs = append(qs, ds.Limit(*opt.Limit))
	}
	if opt.Public != nil {
		qs = append(qs, ds.Filter("Options.Public =", *opt.Public))
	}
	if opt.Owner != nil {
		qs = append(qs, ds.Filter("Owner =", *opt.Owner))
	}
	if opt.Start.From != nil {
		qs = append(qs, ds.Filter("Start >=", *opt.Start.From))
	}
	if opt.Start.To != nil {
		qs = append(qs, ds.Filter("Start <=", *opt.Start.To))
	}

	err := c.client.Query(ctx, kindCourse, &xs, qs...)
	err = ds.IgnoreFieldMismatch(err)
	if err != nil {
		return nil, err
	}

	return xs, nil
}

// CourseDelete delete course from database
func (c *DB) CourseDelete(ctx context.Context, courseID string) error {
	err := c.client.DeleteByStringID(ctx, kindCourse, courseID)
	if err != nil {
		return err
	}
	cacheCourse.Unset(courseID)
	return nil
}

// CourseFind retrieves course from URL
func (c *DB) CourseFind(ctx context.Context, courseURL string) (*model.Course, error) {
	if cache := cacheCourseURL.Get(courseURL); cache != nil {
		return c.CourseGet(ctx, cache.(string))
	}

	var x model.Course
	err := c.client.QueryFirst(ctx, kindCourse, &x, ds.Filter("URL =", courseURL))
	if ds.NotFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	cacheCourseURL.SetTTL(courseURL, x.ID(), time.Minute*5)
	return &x, nil
}

// CourseGetAllByIDs retrieves courses by given course ids
func (c *DB) CourseGetAllByIDs(ctx context.Context, courseIDs []string) (model.Courses, error) {
	if len(courseIDs) == 0 {
		return []*model.Course{}, nil
	}

	xs := make([]*model.Course, len(courseIDs))
	err := c.client.GetByStringIDs(ctx, kindCourse, courseIDs, xs)
	if ds.NotFound(err) {
		return nil, err
	}
	return xs, nil
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
