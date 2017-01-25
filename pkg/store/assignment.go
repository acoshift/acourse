package store

import (
	"context"

	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/ds"
)

// UserAssignmentSave saves user assignment to database
func (c *DB) UserAssignmentSave(ctx context.Context, x *model.UserAssignment) error {
	return c.client.SaveModel(ctx, kindUserAssignment, x)
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

// UserAssignmentList retrieves user assignments
func (c *DB) UserAssignmentList(ctx context.Context, assignmentID, userID string) (model.UserAssignments, error) {
	var xs model.UserAssignments

	qs := []ds.Query{ds.Filter("AssignmentID =", assignmentID)}

	if userID != "" {
		qs = append(qs, ds.Filter("UserID =", userID))
	}

	err := c.client.Query(ctx, kindUserAssignment, &xs, qs...)
	err = ds.IgnoreFieldMismatch(err)
	if err != nil {
		return nil, err
	}
	return xs, nil
}
