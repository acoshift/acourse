package repository

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/gob"
	"time"

	"github.com/acoshift/pgsql"
	"github.com/lib/pq"

	"github.com/acoshift/acourse/context/redisctx"
	"github.com/acoshift/acourse/context/sqlctx"
	"github.com/acoshift/acourse/controller/app"
	"github.com/acoshift/acourse/entity"
)

// NewApp creates new app repository
func NewApp() app.Repository {
	return &appRepo{}
}

type appRepo struct {
}

func (appRepo) GetCourse(ctx context.Context, courseID string) (*entity.Course, error) {
	q := sqlctx.GetQueryer(ctx)

	var x entity.Course
	err := q.QueryRow(`
		select
			id, user_id, title, short_desc, long_desc, image,
			start, url, type, price, courses.discount, enroll_detail,
			opt.public, opt.enroll, opt.attend, opt.assignment, opt.discount
		from courses
			left join course_options as opt on opt.course_id = courses.id
		where id = $1
	`, courseID).Scan(
		&x.ID, &x.UserID, &x.Title, &x.ShortDesc, &x.Desc, &x.Image,
		&x.Start, &x.URL, &x.Type, &x.Price, &x.Discount, &x.EnrollDetail,
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

func (appRepo) IsEnrolled(ctx context.Context, userID string, courseID string) (enrolled bool, err error) {
	q := sqlctx.GetQueryer(ctx)

	err = q.QueryRow(`
		select exists (
			select 1
			from enrolls
			where user_id = $1 and course_id = $2
		)
	`, userID, courseID).Scan(&enrolled)
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

func (appRepo) GetCourseContents(ctx context.Context, courseID string) ([]*entity.CourseContent, error) {
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

	var xs []*entity.CourseContent
	for rows.Next() {
		var x entity.CourseContent
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

func (appRepo) ListPublicCourses(ctx context.Context) ([]*entity.Course, error) {
	// TODO: move cache logic out from repo

	c := redisctx.GetClient(ctx)
	cachePrefix := redisctx.GetPrefix(ctx)

	// look from cache
	{
		bs, err := c.Get(cachePrefix + "cache:list_public_course").Bytes()
		if err == nil {
			var xs []*entity.Course
			err = gob.NewDecoder(bytes.NewReader(bs)).Decode(&xs)
			if err == nil {
				return xs, nil
			}
		}
	}

	xs := make([]*entity.Course, 0)
	m := make(map[string]*entity.Course)
	ids := make([]string, 0)

	q := sqlctx.GetQueryer(ctx)

	{
		rows, err := q.Query(queryListCoursesPublic)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var x entity.Course
			err = scanCourse(rows.Scan, &x)
			if err != nil {
				return nil, err
			}
			xs = append(xs, &x)
			ids = append(ids, x.ID)
			m[x.ID] = &x
		}
		if err = rows.Err(); err != nil {
			return nil, err
		}
		rows.Close()
	}

	rows, err := q.Query(`
		select course_id, count(*)
		from enrolls
		where course_id = any($1)
		group by course_id
	`, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var courseID string
		var cnt int64
		err = rows.Scan(&courseID, &cnt)
		if err != nil {
			return nil, err
		}
		m[courseID].EnrollCount = cnt
	}

	// save to cache
	go func() {
		buf := bytes.Buffer{}
		err := gob.NewEncoder(&buf).Encode(xs)
		if err == nil {
			c.Set(cachePrefix+"cache:list_public_course", buf.Bytes(), 10*time.Second)
		}
	}()

	return xs, nil
}

func (appRepo) ListOwnCourses(ctx context.Context, userID string) ([]*entity.Course, error) {
	q := sqlctx.GetQueryer(ctx)

	rows, err := q.Query(queryListCoursesOwn, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	xs := make([]*entity.Course, 0)
	ids := make([]string, 0)
	m := make(map[string]*entity.Course)
	for rows.Next() {
		var x entity.Course
		err = scanCourse(rows.Scan, &x)
		if err != nil {
			return nil, err
		}
		xs = append(xs, &x)
		ids = append(ids, x.ID)
		m[x.ID] = &x
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	rows.Close()

	rows, err = q.Query(`
		select course_id, count(*)
		from enrolls
		where course_id = any($1)
		group by course_id
	`, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var courseID string
		var cnt int64
		err = rows.Scan(&courseID, &cnt)
		if err != nil {
			return nil, err
		}
		m[courseID].EnrollCount = cnt
	}
	return xs, nil
}

func (appRepo) ListEnrolledCourses(ctx context.Context, userID string) ([]*entity.Course, error) {
	q := sqlctx.GetQueryer(ctx)

	xs := make([]*entity.Course, 0)
	rows, err := q.Query(queryListCoursesEnrolled, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var x entity.Course
		err = scanCourse(rows.Scan, &x)
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
