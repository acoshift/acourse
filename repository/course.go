package repository

import (
	"bytes"
	"database/sql"
	"encoding/gob"

	"github.com/garyburd/redigo/redis"
	"github.com/lib/pq"

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
func GetCourses(q Queryer, courseIDs []string) ([]*entity.Course, error) {
	xs := make([]*entity.Course, 0, len(courseIDs))
	rows, err := q.Query(queryGetCourses, pq.Array(courseIDs))
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

// GetCourse gets course
func GetCourse(q Queryer, courseID string) (*entity.Course, error) {
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
func GetCourseContents(q Queryer, courseID string) ([]*entity.CourseContent, error) {
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
func GetCourseContent(q Queryer, courseContentID string) (*entity.CourseContent, error) {
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
func GetCourseIDFromURL(q Queryer, url string) (string, error) {
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
func ListCourses(q Queryer, limit, offset int64) ([]*entity.Course, error) {
	xs := make([]*entity.Course, 0)
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
func ListPublicCourses(q Queryer, redisPool *redis.Pool, cachePrefix string) ([]*entity.Course, error) {
	// TODO: move cache logic out from repo

	// look from cache
	{
		c := redisPool.Get()
		bs, err := redis.Bytes(c.Do("GET", cachePrefix+"cache:list_public_course"))
		c.Close()
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
			c := redisPool.Get()
			c.Do("SETEX", cachePrefix+"cache:list_public_course", 5, buf.Bytes())
			c.Close()
		}
	}()

	return xs, nil
}

// ListOwnCourses lists courses that owned by given user
// TODO: add pagination
func ListOwnCourses(q Queryer, userID string) ([]*entity.Course, error) {
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
func ListEnrolledCourses(q Queryer, userID string) ([]*entity.Course, error) {
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
func CountCourses(q Queryer) (cnt int64, err error) {
	err = q.QueryRow(`
		SELECT count(*)
		  FROM courses;
	`).Scan(&cnt)
	return
}
