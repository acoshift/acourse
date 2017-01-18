package store

import (
	"context"

	"cloud.google.com/go/datastore"
	"github.com/acoshift/acourse/pkg/model"
)

const kindAssignment = "Assignment"
const kindUserAssignment = "UserAssignment"

// AssignmentList retrieves assignments from course id
func (c *DB) AssignmentList(ctx context.Context, courseID string) (model.Assignments, error) {
	q := datastore.
		NewQuery(kindAssignment).
		Filter("CourseID =", courseID)

	var xs model.Assignments
	err := c.getAll(ctx, q, &xs)
	if err != nil {
		return nil, err
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
	return c.save(ctx, kindAssignment, x)
}

// AssignmentDelete deletes assignment from database
func (c *DB) AssignmentDelete(ctx context.Context, assignmentID string) error {
	return c.deleteByIDStr(ctx, kindAssignment, assignmentID)
}

// UserAssignmentSave saves user assignment to database
func (c *DB) UserAssignmentSave(ctx context.Context, x *model.UserAssignment) error {
	x.Stamp()
	return c.save(ctx, kindUserAssignment, x)
}

// UserAssignmentGet retrieves an User Assignment from database
func (c *DB) UserAssignmentGet(ctx context.Context, userAssignmentID string) (*model.UserAssignment, error) {
	id := idInt(userAssignmentID)
	if id == 0 {
		return nil, nil
	}

	var x model.UserAssignment
	err := c.get(ctx, datastore.IDKey(kindUserAssignment, id, nil), &x)
	if notFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &x, nil
}

// UserAssignmentDelete deletes user assignment from database
func (c *DB) UserAssignmentDelete(ctx context.Context, userAssignmentID string) error {
	return c.deleteByIDStr(ctx, kindUserAssignment, userAssignmentID)
}
