package repository

import (
	"context"
	"database/sql"

	"github.com/acoshift/pgsql"

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

// GetCourseContents gets course contents for given course id
func GetCourseContents(ctx context.Context, courseID string) ([]*entity.CourseContent, error) {
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

// GetCourseContent gets course content from id
func GetCourseContent(ctx context.Context, courseContentID string) (*entity.CourseContent, error) {
	q := sqlctx.GetQueryer(ctx)

	var x entity.CourseContent
	err := q.QueryRow(`
		select
			id, course_id, title, long_desc, video_id, video_type, download_url
		from course_contents
		where id = $1
	`, courseContentID).Scan(
		&x.ID, &x.CourseID, &x.Title, &x.Desc, &x.VideoID, &x.VideoType, &x.DownloadURL,
	)
	if err == sql.ErrNoRows {
		return nil, entity.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &x, nil
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
	if err == sql.ErrNoRows {
		err = entity.ErrNotFound
	}
	return
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

// GetCourse get course by id
func GetCourse(ctx context.Context, courseID string) (*entity.Course, error) {
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
