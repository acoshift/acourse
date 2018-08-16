package repository

import (
	"context"
	"database/sql"

	"github.com/acoshift/pgsql"
	"github.com/lib/pq"

	"github.com/acoshift/acourse/context/sqlctx"
	"github.com/acoshift/acourse/controller/admin"
	"github.com/acoshift/acourse/entity"
)

// NewAdmin creates admin repository
func NewAdmin() admin.Repository {
	return &adminRepo{}
}

type adminRepo struct {
}

func (adminRepo) ListUsers(ctx context.Context, limit, offset int64) ([]*entity.UserItem, error) {
	q := sqlctx.GetQueryer(ctx)

	rows, err := q.Query(`
		select
			id, name, username, email,
			image, created_at
		from users
		order by created_at desc
		limit $1 offset $2
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var xs []*entity.UserItem
	for rows.Next() {
		var x entity.UserItem
		err = rows.Scan(
			&x.ID, &x.Name, &x.Username, pgsql.NullString(&x.Email),
			&x.Image, &x.CreatedAt,
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

func (adminRepo) CountUsers(ctx context.Context) (cnt int64, err error) {
	q := sqlctx.GetQueryer(ctx)

	err = q.QueryRow(`select count(*) from users`).Scan(&cnt)
	return
}

func (adminRepo) ListCourses(ctx context.Context, limit, offset int64) ([]*entity.Course, error) {
	q := sqlctx.GetQueryer(ctx)

	rows, err := q.Query(`
		  select courses.id, courses.title, courses.short_desc, courses.long_desc, courses.image,
		         courses.start, courses.url, courses.type, courses.price, courses.discount,
		         courses.enroll_detail, courses.created_at, courses.updated_at,
		         opt.public, opt.enroll, opt.attend, opt.assignment, opt.discount,
		         users.id, users.username, users.image
		    from courses
		         left join course_options as opt
		         on opt.course_id = courses.id
		         left join users
		         on users.id = courses.user_id
		order by courses.created_at desc
		   limit $1
		  offset $2
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

func (adminRepo) CountCourses(ctx context.Context) (cnt int64, err error) {
	q := sqlctx.GetQueryer(ctx)

	err = q.QueryRow(`select count(*) from courses`).Scan(&cnt)
	return
}

func (adminRepo) GetPayment(ctx context.Context, paymentID string) (*entity.Payment, error) {
	q := sqlctx.GetQueryer(ctx)

	var x entity.Payment
	err := q.QueryRow(`
		select
			payments.id,
			payments.image,
			payments.price,
			payments.original_price,
			payments.code,
			payments.status,
			payments.created_at,
			payments.updated_at,
			payments.at,
			users.id,
			users.username,
			users.name,
			users.email,
			users.image,
			courses.id,
			courses.title,
			courses.image,
			courses.url
		from payments
			left join users on payments.user_id = users.id
			left join courses on payments.course_id = courses.id
		where payments.id = $1
	`, paymentID).Scan(
		&x.ID,
		&x.Image, &x.Price, &x.OriginalPrice, &x.Code, &x.Status, &x.CreatedAt, &x.UpdatedAt, &x.At,
		&x.User.ID, &x.User.Username, &x.User.Name, pgsql.NullString(&x.User.Email), &x.User.Image,
		&x.Course.ID, &x.Course.Title, &x.Course.Image, &x.Course.URL,
	)
	if err == sql.ErrNoRows {
		return nil, entity.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	x.UserID = x.User.ID
	x.CourseID = x.Course.ID
	return &x, nil
}

func (adminRepo) ListPaymentsByStatus(ctx context.Context, statuses []int, limit, offset int64) ([]*entity.Payment, error) {
	q := sqlctx.GetQueryer(ctx)

	rows, err := q.Query(`
		select
			payments.id,
			payments.image,
			payments.price,
			payments.original_price,
			payments.code,
			payments.status,
			payments.created_at,
			payments.updated_at,
			payments.at,
			users.id,
			users.username,
			users.name,
			users.email,
			users.image,
			courses.id,
			courses.title,
			courses.image,
			courses.url
		from payments
			left join users on payments.user_id = users.id
			left join courses on payments.course_id = courses.id
		where payments.status = any($1)
		order by payments.created_at desc
		limit $2 offset $3
	`, pq.Array(statuses), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var xs []*entity.Payment
	for rows.Next() {
		var x entity.Payment
		err = rows.Scan(
			&x.ID,
			&x.Image, &x.Price, &x.OriginalPrice, &x.Code, &x.Status, &x.CreatedAt, &x.UpdatedAt, &x.At,
			&x.User.ID, &x.User.Username, &x.User.Name, pgsql.NullString(&x.User.Email), &x.User.Image,
			&x.Course.ID, &x.Course.Title, &x.Course.Image, &x.Course.URL,
		)
		if err != nil {
			return nil, err
		}
		x.UserID = x.User.ID
		x.CourseID = x.Course.ID
		xs = append(xs, &x)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return xs, nil
}

func (adminRepo) CountPaymentsByStatuses(ctx context.Context, statuses []int) (cnt int64, err error) {
	q := sqlctx.GetQueryer(ctx)

	err = q.QueryRow(`
		select count(*)
		from payments
		where status = any($1)
	`, pq.Array(statuses)).Scan(&cnt)
	return
}
