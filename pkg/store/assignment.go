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
	var x model.Assignment
	err := c.getByIDStr(ctx, kindAssignment, assignmentID, &x)
	if err == ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &x, nil
}

// AssignmentGetMulti retrieves assignments from database
func (c *DB) AssignmentGetMulti(ctx context.Context, assignmentIDs []string) (model.Assignments, error) {
	keys := make([]*datastore.Key, len(assignmentIDs))
	for i, id := range assignmentIDs {
		keys[i] = datastore.IDKey(kindAssignment, idInt(id), nil)
	}
	xs := make(model.Assignments, len(assignmentIDs))
	err := c.client.GetMulti(ctx, keys, xs)
	if multiError(err) {
		return nil, err
	}
	for i, x := range xs {
		x.SetKey(keys[i])
	}
	return xs, nil
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
	var x model.UserAssignment
	err := c.getByIDStr(ctx, kindUserAssignment, userAssignmentID, &x)
	if err == ErrNotFound {
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

// UserAssignmentGetMulti retrieves assignments from database
func (c *DB) UserAssignmentGetMulti(ctx context.Context, userAssignmentIDs []string) (model.UserAssignments, error) {
	keys := make([]*datastore.Key, len(userAssignmentIDs))
	for i, id := range userAssignmentIDs {
		keys[i] = datastore.IDKey(kindUserAssignment, idInt(id), nil)
	}
	xs := make(model.UserAssignments, len(userAssignmentIDs))
	err := c.client.GetMulti(ctx, keys, xs)
	if multiError(err) {
		return nil, err
	}
	for i, x := range xs {
		x.SetKey(keys[i])
	}
	return xs, nil
}
