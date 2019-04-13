package course

import (
	"context"
	"time"

	"github.com/acoshift/pgsql"
	"github.com/acoshift/pgsql/pgctx"
)

// IsEnroll checks is user enrolled a course
func IsEnroll(ctx context.Context, userID, courseID string) (bool, error) {
	var b bool
	err := pgctx.QueryRow(ctx, `
		select exists (
			select 1
			from enrolls
			where user_id = $1 and course_id = $2
		)
	`, userID, courseID).Scan(&b)
	return b, err
}

// EnrolledCourse type
type EnrolledCourse struct {
	ID    string
	Title string
	Desc  string
	Image string
	Start time.Time
	URL   string
	Type  int
}

// Link returns course link
func (x *EnrolledCourse) Link() string {
	if x.URL != "" {
		return x.URL
	}
	return x.ID
}

// ShowStart returns true if course should show start date
func (x *EnrolledCourse) ShowStart() bool {
	return x.Type == Live && !x.Start.IsZero()
}

func GetEnrolledCourses(ctx context.Context, userID string) ([]*EnrolledCourse, error) {
	// language=SQL
	rows, err := pgctx.Query(ctx, `
		select
			c.id,
			c.title, c.short_desc, c.image,
			c.start, c.url, c.type
		from courses as c
			inner join enrolls as e on c.id = e.course_id
		where e.user_id = $1
		order by e.created_at desc
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var xs []*EnrolledCourse
	for rows.Next() {
		var x EnrolledCourse
		err = rows.Scan(
			&x.ID,
			&x.Title, &x.Desc, &x.Image,
			pgsql.NullTime(&x.Start), pgsql.NullString(&x.URL), &x.Type,
		)
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
