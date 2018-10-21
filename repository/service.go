package repository

import (
	"context"
	"database/sql"

	"github.com/acoshift/pgsql"

	"github.com/acoshift/acourse/context/sqlctx"
	"github.com/acoshift/acourse/entity"
	"github.com/acoshift/acourse/model/course"
	"github.com/acoshift/acourse/service"
)

// NewService creates new service repository
func NewService() service.Repository {
	return &svcRepo{}
}

type svcRepo struct {
}

func (svcRepo) GetUserByEmail(ctx context.Context, email string) (*service.User, error) {
	q := sqlctx.GetQueryer(ctx)

	var x service.User
	err := q.QueryRow(`
		select
			id, name, email
		from users
		where email = $1
	`, email).Scan(
		&x.ID, &x.Name, pgsql.NullString(&x.Email),
	)
	if err == sql.ErrNoRows {
		return nil, entity.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &x, nil
}

func (svcRepo) RegisterCourse(ctx context.Context, x *service.RegisterCourse) (courseID string, err error) {
	q := sqlctx.GetQueryer(ctx)

	err = q.QueryRow(`
		insert into courses
			(user_id, title, short_desc, long_desc, image, start)
		values
			($1, $2, $3, $4, $5, $6)
		returning id
	`, x.UserID, x.Title, x.ShortDesc, x.LongDesc, x.Image, pgsql.NullTime(&x.Start)).Scan(&courseID)
	return
}

func (svcRepo) GetCourse(ctx context.Context, courseID string) (*entity.Course, error) {
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
		&x.Start, pgsql.NullString(&x.URL), &x.Type, &x.Price, &x.Discount, &x.EnrollDetail,
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

func (svcRepo) UpdateCourse(ctx context.Context, x *service.UpdateCourseModel) error {
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
	`, x.ID, x.Title, x.ShortDesc, x.LongDesc, pgsql.NullTime(&x.Start))
	return err
}

func (svcRepo) SetCourseImage(ctx context.Context, courseID string, image string) error {
	q := sqlctx.GetQueryer(ctx)

	_, err := q.Exec(`update courses set image = $2 where id = $1`, courseID, image)
	return err
}

func (svcRepo) SetCourseOption(ctx context.Context, courseID string, x *entity.CourseOption) error {
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

func (svcRepo) RegisterCourseContent(ctx context.Context, x *entity.RegisterCourseContent) (contentID string, err error) {
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

func (svcRepo) GetCourseContent(ctx context.Context, contentID string) (*course.Content, error) {
	q := sqlctx.GetQueryer(ctx)

	var x course.Content
	err := q.QueryRow(`
		select
			id, course_id, title, long_desc, video_id, video_type, download_url
		from course_contents
		where id = $1
	`, contentID).Scan(
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

func (svcRepo) ListCourseContents(ctx context.Context, courseID string) ([]*course.Content, error) {
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

func (svcRepo) UpdateCourseContent(ctx context.Context, contentID, title, desc, videoID string) error {
	q := sqlctx.GetQueryer(ctx)

	_, err := q.Exec(`
		update course_contents
		set
			title = $2,
			long_desc = $3,
			video_id = $4,
			updated_at = now()
		where id = $1
	`, contentID, title, desc, videoID)
	return err
}

func (svcRepo) DeleteCourseContent(ctx context.Context, contentID string) error {
	q := sqlctx.GetQueryer(ctx)

	_, err := q.Exec(`delete from course_contents where id = $1`, contentID)
	return err
}

func (svcRepo) RegisterPayment(ctx context.Context, x *service.RegisterPayment) error {
	q := sqlctx.GetQueryer(ctx)

	_, err := q.Exec(`
		insert into payments
			(user_id, course_id, image, price, original_price, code, status)
		values
			($1, $2, $3, $4, $5, $6, $7)
		returning id
	`, x.UserID, x.CourseID, x.Image, x.Price, x.OriginalPrice, x.Code, x.Status)
	if err != nil {
		return err
	}
	return nil
}

func (svcRepo) GetPayment(ctx context.Context, paymentID string) (*service.Payment, error) {
	q := sqlctx.GetQueryer(ctx)

	var x service.Payment
	err := q.QueryRow(`
		select
			p.id,
			p.image, p.price, p.original_price, p.code,
			p.status, p.created_at, p.at,
			u.id, u.username, u.name, u.email, u.image,
			c.id, c.title, c.image, c.url
		from payments as p
			left join users as u on p.user_id = u.id
			left join courses as c on p.course_id = c.id
		where p.id = $1
	`, paymentID).Scan(
		&x.ID,
		&x.Image, &x.Price, &x.OriginalPrice, &x.Code,
		&x.Status, &x.CreatedAt, &x.At,
		&x.User.ID, &x.User.Username, &x.User.Name, pgsql.NullString(&x.User.Email), &x.User.Image,
		&x.Course.ID, &x.Course.Title, &x.Course.Image, pgsql.NullString(&x.Course.URL),
	)
	if err == sql.ErrNoRows {
		return nil, entity.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &x, nil
}

func (svcRepo) SetPaymentStatus(ctx context.Context, paymentID string, status int) error {
	q := sqlctx.GetQueryer(ctx)

	_, err := q.Exec(`
		update payments
		set
			status = $2,
			updated_at = now(),
			at = now()
		where id = $1
	`, paymentID, status)
	return err
}

func (svcRepo) HasPendingPayment(ctx context.Context, userID string, courseID string) (exists bool, err error) {
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

func (svcRepo) RegisterEnroll(ctx context.Context, userID string, courseID string) error {
	q := sqlctx.GetQueryer(ctx)

	_, err := q.Exec(`
		insert into enrolls
			(user_id, course_id)
		values
			($1, $2)
	`, userID, courseID)
	return err
}

func (svcRepo) IsEnrolled(ctx context.Context, userID string, courseID string) (enrolled bool, err error) {
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
