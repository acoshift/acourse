package admin

import (
	"context"
	"time"

	"github.com/acoshift/pgsql"

	"github.com/acoshift/acourse/internal/pkg/context/sqlctx"
	"github.com/acoshift/acourse/internal/pkg/course"
)

// CourseItem type
type CourseItem struct {
	ID        string
	Title     string
	Image     string
	Type      int
	Price     float64
	Discount  float64
	URL       string
	CreatedAt time.Time
	UpdatedAt time.Time
	Option    course.Option
	Owner     struct {
		ID       string
		Username string
		Image    string
	}
}

func GetCourses(ctx context.Context, limit, offset int64) ([]*CourseItem, error) {
	rows, err := sqlctx.Query(ctx, `
		select
			c.id, c.title, c.image,
			c.url, c.type, c.price, c.discount,
			c.created_at, c.updated_at,
			opt.public, opt.enroll, opt.attend, opt.assignment, opt.discount,
			u.id, u.username, u.image
		from courses as c
			left join course_options as opt on opt.course_id = c.id
			left join users as u on u.id = c.user_id
		order by c.created_at desc
		limit $1 offset $2
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var xs []*CourseItem
	for rows.Next() {
		var x CourseItem
		err = rows.Scan(
			&x.ID, &x.Title, &x.Image,
			pgsql.NullString(&x.URL), &x.Type, &x.Price, &x.Discount,
			&x.CreatedAt, &x.UpdatedAt,
			&x.Option.Public, &x.Option.Enroll, &x.Option.Attend, &x.Option.Assignment, &x.Option.Discount,
			&x.Owner.ID, &x.Owner.Username, &x.Owner.Image,
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

func CountCourses(ctx context.Context) (cnt int64, err error) {
	err = sqlctx.QueryRow(ctx,
		`select count(*) from courses`,
	).Scan(&cnt)
	return
}
