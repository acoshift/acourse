package app

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/gob"
	"time"

	"github.com/acoshift/pgsql"

	"github.com/acoshift/acourse/internal/pkg/app"
	"github.com/acoshift/acourse/internal/pkg/context/redisctx"
	"github.com/acoshift/acourse/internal/pkg/context/sqlctx"
	"github.com/acoshift/acourse/internal/pkg/course"
	"github.com/acoshift/acourse/internal/pkg/payment"
	"github.com/acoshift/acourse/internal/pkg/user"
)

func getCourse(ctx context.Context, courseID string) (*Course, error) {
	var x Course
	err := sqlctx.QueryRow(ctx, `
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

func getCourseIDByURL(ctx context.Context, url string) (courseID string, err error) {
	err = sqlctx.QueryRow(ctx, `
		select id
		from courses
		where url = $1
	`, url).Scan(&courseID)
	if err == sql.ErrNoRows {
		err = app.ErrNotFound
	}
	return
}

func hasPendingPayment(ctx context.Context, userID string, courseID string) (exists bool, err error) {
	err = sqlctx.QueryRow(ctx, `
		select exists (
			select 1
			from payments
			where user_id = $1 and course_id = $2 and status = $3
		)
	`, userID, courseID, payment.Pending).Scan(&exists)
	return
}

func getCourseContents(ctx context.Context, courseID string) ([]*course.Content, error) {
	rows, err := sqlctx.Query(ctx, `
		select
			id, course_id, title, long_desc, video_id, video_type, download_url
		from course_contents
		where course_id = $1
		order by i
	`, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var xs []*course.Content
	for rows.Next() {
		var x course.Content
		err = rows.Scan(
			&x.ID, &x.CourseID, &x.Title, &x.Desc, &x.VideoID, &x.VideoType, &x.DownloadURL,
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

func getUser(ctx context.Context, userID string) (*user.User, error) {
	var x user.User
	err := sqlctx.QueryRow(ctx, `
		select
			users.id,
			users.name,
			users.username,
			users.email,
			users.about_me,
			users.image,
			coalesce(roles.admin, false),
			coalesce(roles.instructor, false)
		from users
			left join roles on users.id = roles.user_id
		where users.id = $1
	`, userID).Scan(
		&x.ID, &x.Name, &x.Username, pgsql.NullString(&x.Email), &x.AboutMe, &x.Image,
		&x.Role.Admin, &x.Role.Instructor,
	)
	if err == sql.ErrNoRows {
		return nil, app.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &x, nil
}

func findAssignmentsByCourseID(ctx context.Context, courseID string) ([]*course.Assignment, error) {
	rows, err := sqlctx.Query(ctx, `
		select id, title, long_desc, open
		from assignments
		where course_id = $1
		order by i asc
	`, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var xs []*course.Assignment
	for rows.Next() {
		var x course.Assignment
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

	rows, err := sqlctx.Query(ctx, `
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

func listOwnCourses(ctx context.Context, userID string) ([]*OwnCourse, error) {
	rows, err := sqlctx.Query(ctx, `
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

func listEnrolledCourses(ctx context.Context, userID string) ([]*EnrolledCourse, error) {
	rows, err := sqlctx.Query(ctx, `
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
