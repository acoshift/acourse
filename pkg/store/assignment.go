package store

import (
	"context"

	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/ds"
)

// AssignmentList retrieves assignments from course id
func (c *DB) AssignmentList(ctx context.Context, courseID string) (model.Assignments, error) {
	var xs model.Assignments
	err := c.client.Query(ctx, kindAssignment, &xs, ds.Filter("CourseID =", courseID))
	if err != nil {
		return nil, err
	}

	return xs, nil
}

// AssignmentGet retrieves assignment from database
func (c *DB) AssignmentGet(ctx context.Context, assignmentID string) (*model.Assignment, error) {
	var x model.Assignment
	err := c.client.GetByStringID(ctx, kindAssignment, assignmentID, &x)
	if ds.NotFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &x, nil
}

// AssignmentGetMulti retrieves assignments from database
func (c *DB) AssignmentGetMulti(ctx context.Context, assignmentIDs []string) (model.Assignments, error) {
	xs := make(model.Assignments, len(assignmentIDs))
	err := c.client.GetByStringIDs(ctx, kindAssignment, assignmentIDs, xs)
	err = ds.IgnoreFieldMismatch(err)
	if err != nil {
		return nil, err
	}
	return xs, nil
}

// AssignmentSave saves assignment to database
func (c *DB) AssignmentSave(ctx context.Context, x *model.Assignment) error {
	return c.client.Save(ctx, kindAssignment, x)
}

// AssignmentDelete deletes assignment from database
func (c *DB) AssignmentDelete(ctx context.Context, assignmentID string) error {
	return c.client.DeleteByStringID(ctx, kindAssignment, assignmentID)
}

// UserAssignmentSave saves user assignment to database
func (c *DB) UserAssignmentSave(ctx context.Context, x *model.UserAssignment) error {
	return c.client.Save(ctx, kindUserAssignment, x)
}

// UserAssignmentGet retrieves an User Assignment from database
func (c *DB) UserAssignmentGet(ctx context.Context, userAssignmentID string) (*model.UserAssignment, error) {
	var x model.UserAssignment
	err := c.client.GetByStringID(ctx, kindUserAssignment, userAssignmentID, &x)
	if ds.NotFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &x, nil
}

// UserAssignmentDelete deletes user assignment from database
func (c *DB) UserAssignmentDelete(ctx context.Context, userAssignmentID string) error {
	return c.client.DeleteByStringID(ctx, kindUserAssignment, userAssignmentID)
}

// UserAssignmentGetMulti retrieves assignments from database
func (c *DB) UserAssignmentGetMulti(ctx context.Context, userAssignmentIDs []string) (model.UserAssignments, error) {
	xs := make(model.UserAssignments, len(userAssignmentIDs))
	err := c.client.GetByStringIDs(ctx, kindUserAssignment, userAssignmentIDs, xs)
	err = ds.IgnoreFieldMismatch(err)
	if err != nil {
		return nil, err
	}
	return xs, nil
}
