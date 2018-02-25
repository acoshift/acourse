package repository

import (
	"github.com/acoshift/acourse/entity"
)

// GetAssignments gets assignments
func GetAssignments(q Queryer, courseID string) ([]*entity.Assignment, error) {
	rows, err := q.Query(`
		  SELECT id, title, long_desc, open
		    FROM assignments
		   WHERE course_id = $1
		ORDER BY i ASC;
	`, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	xs := make([]*entity.Assignment, 0)
	for rows.Next() {
		var x entity.Assignment
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
