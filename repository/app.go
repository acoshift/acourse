package repository

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/gob"
	"time"

	"github.com/acoshift/pgsql"

	"github.com/acoshift/acourse/context/redisctx"
	"github.com/acoshift/acourse/context/sqlctx"
	"github.com/acoshift/acourse/controller/app"
	"github.com/acoshift/acourse/entity"
	"github.com/acoshift/acourse/model/course"
)

// NewApp creates new app repository
func NewApp() app.Repository {
	return &appRepo{}
}

type appRepo struct {
}

func (appRepo) GetCourse(ctx context.Context, courseID string) (*app.Course, error) {
	q := sqlctx.GetQueryer(ctx)

	var x app.Course
	err := q.QueryRow(`
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
		return nil, entity.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &x, nil
}

func (appRepo) GetCourseIDByURL(ctx context.Context, url string) (courseID string, err error) {
	q := sqlctx.GetQueryer(ctx)

	err = q.QueryRow(`
		select id
		from courses
		where url = $1
	`, url).Scan(&courseID)
	if err == sql.ErrNoRows {
		err = entity.ErrNotFound
	}
	return
}

func (appRepo) HasPendingPayment(ctx context.Context, userID string, courseID string) (exists bool, err error) {
	q := sqlctx.GetQueryer(ctx)

	err = q.QueryRow(`
		select exists (
			select 1
			from payments
			where user_id = $1 and course_id = $2 and status = $3
		)
	`, userID, courseID, entity.Pending).Scan(&exists)
	return
}

func (appRepo) GetCourseContents(ctx context.Context, courseID string) ([]*course.Content, error) {
	q := sqlctx.GetQueryer(ctx)

	rows, err := q.Query(`
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

func (appRepo) GetUser(ctx context.Context, userID string) (*entity.User, error) {
	q := sqlctx.GetQueryer(ctx)

	var x entity.User
	err := q.QueryRow(`
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
		return nil, entity.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &x, nil
}

func (appRepo) FindAssignmentsByCourseID(ctx context.Context, courseID string) ([]*entity.Assignment, error) {
	q := sqlctx.GetQueryer(ctx)

	rows, err := q.Query(`
		select id, title, long_desc, open
		from assignments
		where course_id = $1
		order by i asc
	`, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var xs []*entity.Assignment
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

func (appRepo) ListPublicCourses(ctx context.Context) ([]*app.PublicCourse, error) {
	c := redisctx.GetClient(ctx)
	cachePrefix := redisctx.GetPrefix(ctx)

	// look from cache
	{
		bs, err := c.Get(cachePrefix + "cache:list_public_course").Bytes()
		if err == nil {
			var xs []*app.PublicCourse
			err = gob.NewDecoder(bytes.NewReader(bs)).Decode(&xs)
			if err == nil {
				return xs, nil
			}
		}
	}

	q := sqlctx.GetQueryer(ctx)

	rows, err := q.Query(`
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

	var xs []*app.PublicCourse
	for rows.Next() {
		var x app.PublicCourse
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

func (appRepo) ListOwnCourses(ctx context.Context, userID string) ([]*app.OwnCourse, error) {
	q := sqlctx.GetQueryer(ctx)

	rows, err := q.Query(`
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

	var xs []*app.OwnCourse
	for rows.Next() {
		var x app.OwnCourse
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

func (appRepo) ListEnrolledCourses(ctx context.Context, userID string) ([]*app.EnrolledCourse, error) {
	q := sqlctx.GetQueryer(ctx)

	rows, err := q.Query(`
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

	var xs []*app.EnrolledCourse
	for rows.Next() {
		var x app.EnrolledCourse
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
