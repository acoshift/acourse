package course

import (
	"context"
	"time"

	"github.com/acoshift/pgsql"
	"github.com/acoshift/pgsql/pgctx"
)

// OwnCourse type
type OwnCourse struct {
	ID          string
	Title       string
	Desc        string
	Image       string
	Start       time.Time
	URL         string
	Type        int
	EnrollCount int
}

// Link returns course link
func (x *OwnCourse) Link() string {
	if x.URL != "" {
		return x.URL
	}
	return x.ID
}

// ShowStart returns true if course should show start date
func (x *OwnCourse) ShowStart() bool {
	return x.Type == Live && !x.Start.IsZero()
}

func GetOwnCourses(ctx context.Context, userID string) ([]*OwnCourse, error) {
	// language=SQL
	rows, err := pgctx.Query(ctx, `
		select
			c.id,
			c.title, c.short_desc, c.image,
			c.start, c.url, c.type,
			count(e.user_id)
		from courses as c
			left join enrolls as e on e.course_id = c.id
		where c.user_id = $1
		group by c.id
		order by c.created_at desc
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var xs []*OwnCourse
	for rows.Next() {
		var x OwnCourse
		err = rows.Scan(
			&x.ID,
			&x.Title, &x.Desc, &x.Image,
			pgsql.NullTime(&x.Start), pgsql.NullString(&x.URL), &x.Type,
			&x.EnrollCount,
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
