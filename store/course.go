package store

import (
	"time"

	"cloud.google.com/go/datastore"
)

// Course model
type Course struct {
	Base
	Timestampable
	Title            string `datastore:",noindex"`
	ShortDescription string `datastore:",noindex"`
	Description      string `datastore:",noindex"` // Markdown
	Photo            string `datastore:",noindex"` // URL
	Owner            string
	Start            time.Time
	URL              string
	Type             CourseType
	Video            string `datastore:",noindex"` // Cover Video
	Price            float64
	DiscountedPrice  float64
	Options          CourseOption
	Contents         []CourseContent `datastore:",noindex"`
}

// CourseOption type
type CourseOption struct {
	Public     bool
	Enroll     bool `datastore:",noindex"`
	Attend     bool `datastore:",noindex"`
	Assignment bool `datastore:",noindex"`
	Purchase   bool
	Discount   bool
}

// CourseContent type
type CourseContent struct {
	Title       string `datastore:",noindex"`
	Description string `datastore:",noindex"` // Markdown
	Video       string `datastore:",noindex"` // Youtube ID
	DownloadURL string `datastore:",noindex"` // Video download link
}

// CourseType type
type CourseType string

// CourseType
const (
	CourseTypeLive  = string("live")
	CourseTypeVideo = string("video")
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

var cacheCourse = NewCache(time.Second * 10)
var cacheCourseURL = NewCache(time.Minute)

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
func (c *DB) CourseGet(courseID string) (*Course, error) {
	id := idInt(courseID)
	if id == 0 {
		return nil, nil
	}

	if cache := cacheCourse.Get(courseID); cache != nil {
		return cache.(*Course), nil
	}

	ctx, cancel := getContext()
	defer cancel()

	var err error
	var x Course

	key := datastore.IDKey(kindCourse, id, nil)
	err = c.client.Get(ctx, key, &x)
	if notFound(err) {
		return nil, nil
	}
	if datastoreError(err) {
		return nil, err
	}
	x.setKey(key)
	cacheCourse.Set(courseID, &x)
	return &x, nil
}

// CourseSave saves course to database
func (c *DB) CourseSave(x *Course) error {
	ctx, cancel := getContext()
	defer cancel()

	// Check duplicate URL
	if x.URL != "" {
		if t, err := c.CourseFind(x.URL); (t != nil && t.ID != x.ID) || err != nil {
			if err != nil {
				return err
			}
			return ErrConflict("course url already exists")
		}
	}

	var err error
	x.Stamp()
	if x.key == nil {
		x.setKey(datastore.IncompleteKey(kindCourse, nil))
	}

	key, err := c.client.Put(ctx, x.key, x)
	if err != nil {
		return err
	}
	x.setKey(key)

	cacheCourse.Del(x.ID)
	return nil
}

// CourseList retrieves courses from database
func (c *DB) CourseList(opts ...CourseListOption) ([]*Course, error) {
	ctx, cancel := getContext()
	defer cancel()

	var xs []*Course

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

	keys, err := c.getAll(ctx, q, &xs)
	if err != nil {
		return nil, err
	}

	for i := range keys {
		xs[i].setKey(keys[i])
	}

	return xs, nil
}

// CourseDelete delete course from database
func (c *DB) CourseDelete(courseID string) error {
	id := idInt(courseID)
	if id == 0 {
		return nil
	}

	ctx, cancel := getContext()
	defer cancel()

	err := c.client.Delete(ctx, datastore.IDKey(kindCourse, id, nil))
	if err != nil {
		return err
	}
	cacheCourse.Del(courseID)
	return nil
}

// CourseFind retrieves course from URL
func (c *DB) CourseFind(courseURL string) (*Course, error) {
	if cache := cacheCourseURL.Get(courseURL); cache != nil {
		return c.CourseGet(cache.(string))
	}

	ctx, cancel := getContext()
	defer cancel()

	var x Course
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
	cacheCourseURL.Set(courseURL, x.ID)
	return &x, nil
}

// CourseGetAllByIDs retrieves courses by given course ids
func (c *DB) CourseGetAllByIDs(courseIDs []string) ([]*Course, error) {
	if len(courseIDs) == 0 {
		return []*Course{}, nil
	}

	keys := make([]*datastore.Key, len(courseIDs))
	for i, id := range courseIDs {
		d := idInt(id)
		if d == 0 {
			return nil, ErrInvalidID
		}
		keys[i] = datastore.IDKey(kindCourse, d, nil)
	}

	ctx, cancel := getContext()
	defer cancel()

	xs := make([]*Course, len(keys))
	err := c.client.GetMulti(ctx, keys, xs)
	if err != nil {
		return nil, err
	}
	for i, x := range xs {
		x.setKey(keys[i])
	}
	return xs, nil
}

// CoursePurge purges all users
func (c *DB) CoursePurge() error {
	return c.purge(kindCourse)
}
