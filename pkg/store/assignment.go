package store

import (
	"context"

	"cloud.google.com/go/datastore"

	"github.com/acoshift/acourse/pkg/model"
)

const kindAssignment = "Assignment"

// AssignmentList retrieves assignments from course id
func (c *DB) AssignmentList(ctx context.Context, courseID string) (model.Assignments, error) {
	q := datastore.
		NewQuery(kindAssignment).
		Filter("CourseID =", courseID)

	var xs model.Assignments
	keys, err := c.getAll(ctx, q, &xs)
	if err != nil {
		return nil, err
	}

	for i, x := range xs {
		x.SetKey(keys[i])
	}

	return xs, nil
}

// AssignmentGet retrieves assignment from database
func (c *DB) AssignmentGet(ctx context.Context, assignmentID string) (*model.Assignment, error) {
	id := idInt(assignmentID)
	if id == 0 {
		return nil, nil
	}

	var x model.Assignment
	err := c.get(ctx, datastore.IDKey(kindAssignment, id, nil), &x)
	if notFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &x, nil
}

// AssignmentSave saves assignment to database
func (c *DB) AssignmentSave(ctx context.Context, x *model.Assignment) error {
	x.Stamp()
	if x.Key() == nil {
		x.NewKey(kindAssignment)
	}

	err := c.put(ctx, x)
	if err != nil {
		return err
	}
	return nil
}

// AssignmentDelete deletes assignment from database
func (c *DB) AssignmentDelete(ctx context.Context, assignmentID string) error {
	id := idInt(assignmentID)
	if id == 0 {
		return nil
	}

	err := c.client.Delete(ctx, datastore.IDKey(kindAssignment, id, nil))
	if err != nil {
		return err
	}
	return nil
}
