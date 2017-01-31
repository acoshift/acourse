package store

import (
	"context"

	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/ds"
)

// CourseGetMulti retrieves multiple courses from database
func (c *DB) CourseGetMulti(ctx context.Context, courseIDs []string) (model.Courses, error) {
	if len(courseIDs) == 0 {
		return []*model.Course{}, nil
	}

	var courses []*model.Course
	err := c.client.GetByStringIDs(ctx, kindCourse, courseIDs, &courses)
	err = ds.IgnoreFieldMismatch(err)
	if err != nil {
		return nil, err
	}
	return courses, nil
}
