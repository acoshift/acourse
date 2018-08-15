package repository

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/gob"
	"time"

	"github.com/acoshift/pgsql"

	"github.com/go-redis/redis"
	"github.com/lib/pq"

	"github.com/acoshift/acourse/context/sqlctx"
	"github.com/acoshift/acourse/entity"
)

const (
	selectCourses = `
		select
			courses.id,
			courses.title,
			courses.short_desc,
			courses.long_desc,
			courses.image,
			courses.start,
			courses.url,
			courses.type,
			courses.price,
			courses.discount,
			courses.enroll_detail,
			courses.created_at,
			courses.updated_at,
			course_options.public,
			course_options.enroll,
			course_options.attend,
			course_options.assignment,
			course_options.discount
		from courses
			left join course_options on courses.id = course_options.course_id
	`

	queryGetCourses = selectCourses + `
		where courses.id = any($1)
	`

	queryListCoursesPublic = selectCourses + `
		where course_options.public = true
		order by
			case when courses.type = 1
				then 1
				else null
			end,
			courses.created_at desc
	`

	queryListCoursesOwn = selectCourses + `
		where courses.user_id = $1
		order by courses.created_at desc
	`

	queryListCoursesEnrolled = selectCourses + `
		inner join enrolls on courses.id = enrolls.course_id
		where enrolls.user_id = $1
		order by enrolls.created_at desc
	`
)

func scanCourse(scan scanFunc, x *entity.Course) error {
	err := scan(&x.ID,
		&x.Title, &x.ShortDesc, &x.Desc, &x.Image, &x.Start, &x.URL, &x.Type, &x.Price, &x.Discount, &x.EnrollDetail,
		&x.CreatedAt, &x.UpdatedAt,
		&x.Option.Public, &x.Option.Enroll, &x.Option.Attend, &x.Option.Assignment, &x.Option.Discount,
	)
	if err != nil {
		return err
	}
	if len(x.URL.String) == 0 {
		x.URL.String = x.ID
	}
	return nil
}

// GetCourses gets courses
func GetCourses(ctx context.Context, courseIDs []string) ([]*entity.Course, error) {
	q := sqlctx.GetQueryer(ctx)

	rows, err := q.Query(queryGetCourses, pq.Array(courseIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	xs := make([]*entity.Course, 0, len(courseIDs))
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

// GetCourse gets course
func GetCourse(ctx context.Context, courseID string) (*entity.Course, error) {
	q := sqlctx.GetQueryer(ctx)

	var x entity.Course
	err := q.QueryRow(`
		SELECT id, user_id, title, short_desc, long_desc, image,
		       start, url, type, price, courses.discount, enroll_detail,
		       opt.public, opt.enroll, opt.attend, opt.assignment, opt.discount
		  FROM courses
		       LEFT JOIN course_options AS opt
		       ON opt.course_id = courses.id
		 WHERE id = $1;
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

// GetCourseContents gets course contents for given course id
func GetCourseContents(ctx context.Context, courseID string) ([]*entity.CourseContent, error) {
	q := sqlctx.GetQueryer(ctx)

	rows, err := q.Query(`
		  SELECT id, course_id, title, long_desc, video_id, video_type, download_url
		    FROM course_contents
		   WHERE course_id = $1
		ORDER BY i ASC;
	`, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	xs := make([]*entity.CourseContent, 0)
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

// GetCourseContent gets course content from id
func GetCourseContent(ctx context.Context, courseContentID string) (*entity.CourseContent, error) {
	q := sqlctx.GetQueryer(ctx)

	var x entity.CourseContent
	err := q.QueryRow(`
		SELECT id, course_id, title, long_desc, video_id, video_type, download_url
		  FROM course_contents
		 WHERE id = $1;
	`, courseContentID).Scan(
		&x.ID, &x.CourseID, &x.Title, &x.Desc, &x.VideoID, &x.VideoType, &x.DownloadURL,
	)
	if err != nil {
		return nil, err
	}
	return &x, nil
}

// GetCourseIDFromURL gets course id from url
func GetCourseIDFromURL(ctx context.Context, url string) (string, error) {
	q := sqlctx.GetQueryer(ctx)

	var id string
	err := q.QueryRow(`
		SELECT id
		  FROM courses
		 WHERE url = $1;
	`, url).Scan(&id)
	if err == sql.ErrNoRows {
		return "", entity.ErrNotFound
	}
	if err != nil {
		return "", err
	}
	return id, nil
}

// ListCourses lists all courses
func ListCourses(ctx context.Context, limit, offset int64) ([]*entity.Course, error) {
	q := sqlctx.GetQueryer(ctx)

	rows, err := q.Query(`
		  SELECT courses.id, courses.title, courses.short_desc, courses.long_desc, courses.image,
		         courses.start, courses.url, courses.type, courses.price, courses.discount,
		         courses.enroll_detail, courses.created_at, courses.updated_at,
		         opt.public, opt.enroll, opt.attend, opt.assignment, opt.discount,
		         users.id, users.username, users.image
		    FROM courses
		         LEFT JOIN course_options AS opt
		         ON opt.course_id = courses.id
		         LEFT JOIN users
		         ON users.id = courses.user_id
		ORDER BY courses.created_at DESC
		   LIMIT $1
		  OFFSET $2;
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	xs := make([]*entity.Course, 0)
	for rows.Next() {
		var x entity.Course
		x.Owner = &entity.User{}
		err := rows.Scan(
			&x.ID, &x.Title, &x.ShortDesc, &x.Desc, &x.Image,
			&x.Start, &x.URL, &x.Type, &x.Price, &x.Discount,
			&x.EnrollDetail, &x.CreatedAt, &x.UpdatedAt,
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

// ListPublicCourses lists public course sort by created at desc
// TODO: add pagination
func ListPublicCourses(ctx context.Context, c *redis.Client, cachePrefix string) ([]*entity.Course, error) {
	// TODO: move cache logic out from repo

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
		  SELECT course_id, count(*)
		    FROM enrolls
		   WHERE course_id = any($1)
		GROUP BY course_id;
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

// ListOwnCourses lists courses that owned by given user
// TODO: add pagination
func ListOwnCourses(ctx context.Context, userID string) ([]*entity.Course, error) {
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
		  SELECT course_id, count(*)
		    FROM enrolls
		   WHERE course_id = any($1)
		GROUP BY course_id;
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

// ListEnrolledCourses lists courses that enrolled by given user
// TODO: add pagination
func ListEnrolledCourses(ctx context.Context, userID string) ([]*entity.Course, error) {
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

// CountCourses counts courses
func CountCourses(ctx context.Context) (cnt int64, err error) {
	q := sqlctx.GetQueryer(ctx)

	err = q.QueryRow(`
		SELECT count(*)
		  FROM courses;
	`).Scan(&cnt)
	return
}

// DeleteCourseContent deletes course content
func DeleteCourseContent(ctx context.Context, courseID string, contentID string) error {
	q := sqlctx.GetQueryer(ctx)

	_, err := q.Exec(`
		delete from course_contents
		where id = $1 and course_id = $2
	`, contentID, courseID)
	return err
}

// UpdateCourseContent updates course content
func UpdateCourseContent(ctx context.Context, courseID, contentID, title, desc, videoID string) error {
	q := sqlctx.GetQueryer(ctx)

	_, err := q.Exec(`
		update course_contents
		set
			title = $3,
			long_desc = $4,
			video_id = $5,
			updated_at = now()
		where id = $1 and course_id = $2
	`, contentID, courseID, title, desc, videoID)
	return err
}

// RegisterCourse registers new course
func RegisterCourse(ctx context.Context, x *entity.RegisterCourse) (courseID string, err error) {
	q := sqlctx.GetQueryer(ctx)

	err = q.QueryRow(`
		insert into courses
			(user_id, title, short_desc, long_desc, image, start)
		values
			($1, $2, $3, $4, $5, $6)
		returning id
	`, x.UserID, x.Title, x.ShortDesc, x.LongDesc, x.Image, x.Start).Scan(&courseID)
	return
}

// SetCourseOption sets course option
func SetCourseOption(ctx context.Context, courseID string, x *entity.CourseOption) error {
	q := sqlctx.GetQueryer(ctx)

	_, err := q.Exec(`
		insert into course_options
			(course_id, public, enroll, attend, assignment, discount)
		values
			($1, $2, $3, $4, $5, $6)
		on conflict (course_id) do update set
			public = excluded.public,
			enroll = excluded.enroll,
			attend = excluded.attend,
			assignment = excluded.assignment,
			discount = excluded.discount
	`, courseID, x.Public, x.Enroll, x.Attend, x.Assignment, x.Discount)
	return err
}

// GetCourseURL gets course url
func GetCourseURL(ctx context.Context, courseID string) (url string, err error) {
	q := sqlctx.GetQueryer(ctx)

	err = q.QueryRow(`select url from courses where id = $1`, courseID).Scan(pgsql.NullString(&url))
	return
}

// GetCourseUserID gets course user id
func GetCourseUserID(ctx context.Context, courseID string) (userID string, err error) {
	q := sqlctx.GetQueryer(ctx)

	err = q.QueryRow(`select user_id from courses where id = $1`, courseID).Scan(&userID)
	return
}

// SetCourseImage sets course image
func SetCourseImage(ctx context.Context, courseID string, image string) error {
	q := sqlctx.GetQueryer(ctx)

	_, err := q.Exec(`update courses set image = $2 where id = $1`, courseID, image)
	return err
}

// UpdateCourse updates course
func UpdateCourse(ctx context.Context, x *entity.UpdateCourse) error {
	q := sqlctx.GetQueryer(ctx)

	_, err := q.Exec(`
		update courses
		set
			title = $2,
			short_desc = $3,
			long_desc = $4,
			start = $5,
			updated_at = now()
		where id = $1
	`, x.ID, x.Title, x.ShortDesc, x.LongDesc, x.Start)
	return err
}

// RegisterCourseContent registers course content
func RegisterCourseContent(ctx context.Context, x *entity.RegisterCourseContent) (contentID string, err error) {
	q := sqlctx.GetQueryer(ctx)

	err = q.QueryRow(`
		insert into course_contents
			(
				course_id,
				i,
				title, long_desc, video_id, video_type
			)
		values
			(
				$1,
				(select coalesce(max(i)+1, 0) from course_contents where course_id = $1),
				$2, $3, $4, $5
			)
		returning id
	`,
		x.CourseID,
		x.Title, x.LongDesc, x.VideoID, x.VideoType,
	).Scan(&contentID)
	return
}
