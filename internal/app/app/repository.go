package app

import (
	"bytes"
	"context"
	"encoding/gob"
	"time"

	"github.com/acoshift/pgsql"
	"github.com/acoshift/pgsql/pgctx"

	"github.com/acoshift/acourse/internal/pkg/context/redisctx"
)

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

	// language=SQL
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
