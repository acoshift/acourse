package model

import (
	"context"

	"github.com/acoshift/acourse/pkg/app"
)

// GetAssignments gets assignments
func GetAssignments(ctx context.Context, courseID string) ([]*app.Assignment, error) {
	db := app.GetDatabase(ctx)

	rows, err := db.QueryContext(ctx, `
		select
			id, title, long_desc, open
		from assignments
		where course_id = $1
		order by i asc
	`, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	xs := make([]*app.Assignment, 0)
	for rows.Next() {
		var x app.Assignment
		err = rows.Scan(&x.ID, &x.Title, &x.Desc, &x.Open)
		if err != nil {
			return nil, err
		}
		xs = append(xs, &x)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return xs, nil
}
