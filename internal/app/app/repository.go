package app

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/gob"
	"time"

	"github.com/acoshift/pgsql"
	"github.com/acoshift/pgsql/pgctx"

	"github.com/acoshift/acourse/internal/pkg/app"
	"github.com/acoshift/acourse/internal/pkg/context/redisctx"
)

func getCourse(ctx context.Context, courseID string) (*Course, error) {
	var x Course
	err := pgctx.QueryRow(ctx, `
		select
			c.id, c.title, c.short_desc, c.long_desc, c.image,
			c.start, c.url, c.type, c.price, c.discount, c.enroll_detail,
			u.id, u.name, u.image,
			opt.public, opt.enroll, opt.attend, opt.assignment, opt.discount
		from courses as c
			left join course_options as opt on opt.course_id = c.id
			inner join users as u on u.id = c.user_id
		where c.id = $1
	`, courseID).Scan(
		&x.ID, &x.Title, &x.ShortDesc, &x.Desc, &x.Image,
		pgsql.NullTime(&x.Start), pgsql.NullString(&x.URL), &x.Type, &x.Price, &x.Discount, &x.EnrollDetail,
		&x.Owner.ID, &x.Owner.Name, &x.Owner.Image,
		&x.Option.Public, &x.Option.Enroll, &x.Option.Attend, &x.Option.Assignment, &x.Option.Discount,
	)
	if err == sql.ErrNoRows {
		return nil, app.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &x, nil
}

func listPublicCourses(ctx context.Context) ([]*PublicCourse, error) {
	c := redisctx.GetClient(ctx)
	cachePrefix := redisctx.GetPrefix(ctx)

	// look from cache
	{
		bs, err := c.Get(cachePrefix + "cache:list_public_course").Bytes()
		if err == nil {
			var xs []*PublicCourse
			err = gob.NewDecoder(bytes.NewReader(bs)).Decode(&xs)
			if err == nil {
				return xs, nil
			}
		}
	}

	rows, err := pgctx.Query(ctx, `
			select
				c.id,
				c.title, c.short_desc, c.image, c.start, c.url,
				c.type, c.price, c.discount,
				opt.public, opt.enroll, opt.attend, opt.assignment, opt.discount
			from courses as c
				left join course_options as opt on c.id = opt.course_id
			where opt.public = true
			order by
				case
					when c.type = 1 then 1
					else null
				end,
				c.created_at desc
		`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var xs []*PublicCourse
	for rows.Next() {
		var x PublicCourse
		err = rows.Scan(
			&x.ID,
			&x.Title, &x.Desc, &x.Image, pgsql.NullTime(&x.Start), pgsql.NullString(&x.URL),
			&x.Type, &x.Price, &x.Discount,
			&x.Option.Public, &x.Option.Enroll, &x.Option.Attend, &x.Option.Assignment, &x.Option.Discount,
		)
		if err != nil {
			return nil, err
		}
		xs = append(xs, &x)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// save to cache
	go func() {
		buf := bytes.Buffer{}
		err := gob.NewEncoder(&buf).Encode(xs)
		if err == nil {
			c.Set(cachePrefix+"cache:list_public_course", buf.Bytes(), time.Minute)
		}
	}()

	return xs, nil
}
