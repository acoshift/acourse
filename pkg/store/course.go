package store

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/acoshift/acourse/pkg/model"
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

const kindCourse = "Course"

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
	id := idInt(courseID)
	if id == 0 {
		return nil, nil
	}

	if cache := cacheCourse.Get(courseID); cache != nil {
		return cache.(*model.Course), nil
	}

	var err error
	var x model.Course

	err = c.get(ctx, datastore.IDKey(kindCourse, id, nil), &x)
	if notFound(err) {
		return nil, nil
	}
	if datastoreError(err) {
		return nil, err
	}
	cacheCourse.SetTTL(courseID, &x, time.Minute)
	return &x, nil
}

// CourseSave saves course to database
func (c *DB) CourseSave(ctx context.Context, x *model.Course) error {
	// Check duplicate URL
	if x.URL != "" {
		if t, err := c.CourseFind(ctx, x.URL); (t != nil && t.ID != x.ID) || err != nil {
			if err != nil {
				return err
			}
			return ErrConflict("course url already exists")
		}
	}

	x.Stamp()
	err := c.save(ctx, kindCourse, x)
	if err != nil {
		return err
	}

	cacheCourse.Unset(x.ID)
	return nil
}

// CourseList retrieves courses from database
func (c *DB) CourseList(ctx context.Context, opts ...CourseListOption) (model.Courses, error) {
	var xs []*model.Course

	q := datastore.NewQuery(kindCourse)

	opt := &CourseListOptions{}
	for _, setter := range opts {
		setter(opt)
	}

	if opt.Offset != nil {
		q = q.Offset(*opt.Offset)
	}
	if opt.Limit != nil {
		q = q.Limit(*opt.Limit)
	}
	if opt.Public != nil {
		q = q.Filter("Options.Public =", *opt.Public)
	}
	if opt.Owner != nil {
		q = q.Filter("Owner =", *opt.Owner)
	}
	if opt.Start.From != nil {
		q = q.Filter("Start >=", *opt.Start.From)
	}
	if opt.Start.To != nil {
		q = q.Filter("Start <=", *opt.Start.To)
	}

	err := c.getAll(ctx, q, &xs)
	if err != nil {
		return nil, err
	}

	return xs, nil
}

// CourseDelete delete course from database
func (c *DB) CourseDelete(ctx context.Context, courseID string) error {
	err := c.deleteByIDStr(ctx, kindCourse, courseID)
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
	q := datastore.
		NewQuery(kindCourse).
		Filter("URL =", courseURL).
		Limit(1)

	err := c.findFirst(ctx, q, &x)
	if notFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	cacheCourseURL.SetTTL(courseURL, x.ID, time.Minute*5)
	return &x, nil
}

// CourseGetAllByIDs retrieves courses by given course ids
func (c *DB) CourseGetAllByIDs(ctx context.Context, courseIDs []string) (model.Courses, error) {
	if len(courseIDs) == 0 {
		return []*model.Course{}, nil
	}

	keys := make([]*datastore.Key, len(courseIDs))
	for i, id := range courseIDs {
		d := idInt(id)
		if d == 0 {
			return nil, ErrInvalidID
		}
		keys[i] = datastore.IDKey(kindCourse, d, nil)
	}

	xs := make([]*model.Course, len(keys))
	err := c.client.GetMulti(ctx, keys, xs)
	if datastoreError(err) {
		return nil, err
	}
	for i, x := range xs {
		x.SetKey(keys[i])
	}
	return xs, nil
}

// CourseGetMulti retrieves multiple courses from database
func (c *DB) CourseGetMulti(ctx context.Context, courseIDs []string) (model.Courses, error) {
	if len(courseIDs) == 0 {
		return []*model.Course{}, nil
	}

	courses := make([]*model.Course, 0, len(courseIDs))
	keys := make([]*datastore.Key, 0, len(courseIDs))

	for _, id := range courseIDs {
		if c := cacheCourse.Get(id); c != nil {
			courses = append(courses, c.(*model.Course))
		} else {
			tempID := idInt(id)
			if tempID != 0 {
				keys = append(keys, datastore.IDKey(kindCourse, tempID, nil))
			}
		}
	}
	if len(keys) == 0 {
		return courses, nil
	}

	xs := make([]*model.Course, len(keys))
	err := c.client.GetMulti(ctx, keys, xs)
	if multiError(err) {
		return nil, err
	}
	for i, x := range xs {
		if x == nil {
			continue
		}
		x.SetKey(keys[i])
		courses = append(courses, x)
		cacheCourse.SetTTL(x.ID, x, time.Minute)
	}
	return courses, nil
}
