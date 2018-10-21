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

func (svcRepo) GetCourse(ctx context.Context, courseID string) (*course.Course, error) {
	q := sqlctx.GetQueryer(ctx)

	var x course.Course
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
