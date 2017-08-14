package model

import "context"

// Assignment model
type Assignment struct {
	ID    string
	Title string
	Desc  string
	Open  bool
}

// GetAssignments gets assignments
func GetAssignments(ctx context.Context, db DB, courseID int64) ([]*Assignment, error) {
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
	xs := make([]*Assignment, 0)
	for rows.Next() {
		var x Assignment
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
