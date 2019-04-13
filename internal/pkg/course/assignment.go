package course

import (
	"context"

	"github.com/acoshift/pgsql/pgctx"
)

// Assignment model
type Assignment struct {
	ID    string
	Title string
	Desc  string
	Open  bool
}

func GetAssignments(ctx context.Context, courseID string) ([]*Assignment, error) {
	// language=SQL
	rows, err := pgctx.Query(ctx, `
		select id, title, long_desc, open
		from assignments
		where course_id = $1
		order by i asc
	`, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var xs []*Assignment
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
