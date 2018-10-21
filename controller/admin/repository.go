package admin

import (
	"context"
	"database/sql"

	"github.com/acoshift/pgsql"
	"github.com/lib/pq"

	"github.com/acoshift/acourse/context/sqlctx"
	"github.com/acoshift/acourse/entity"
)

func listUsers(ctx context.Context, limit, offset int64) ([]*UserItem, error) {
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

	var xs []*UserItem
	for rows.Next() {
		var x UserItem
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

func countUsers(ctx context.Context) (cnt int64, err error) {
	q := sqlctx.GetQueryer(ctx)

	err = q.QueryRow(`select count(*) from users`).Scan(&cnt)
	return
}

func listCourses(ctx context.Context, limit, offset int64) ([]*CourseItem, error) {
	q := sqlctx.GetQueryer(ctx)

	rows, err := q.Query(`
		select
			c.id, c.title, c.image,
			c.url, c.type, c.price, c.discount,
			c.created_at, c.updated_at,
			opt.public, opt.enroll, opt.attend, opt.assignment, opt.discount,
			u.id, u.username, u.image
		from courses as c
			left join course_options as opt on opt.course_id = c.id
			left join users as u on u.id = c.user_id
		order by c.created_at desc
		limit $1 offset $2
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var xs []*CourseItem
	for rows.Next() {
		var x CourseItem
		err := rows.Scan(
			&x.ID, &x.Title, &x.Image,
			pgsql.NullString(&x.URL), &x.Type, &x.Price, &x.Discount,
			&x.CreatedAt, &x.UpdatedAt,
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

func countCourses(ctx context.Context) (cnt int64, err error) {
	q := sqlctx.GetQueryer(ctx)

	err = q.QueryRow(`select count(*) from courses`).Scan(&cnt)
	return
}

func getPayment(ctx context.Context, paymentID string) (*Payment, error) {
	q := sqlctx.GetQueryer(ctx)

	var x Payment
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
		&x.Status, &x.CreatedAt, pgsql.NullTime(&x.At),
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

func listPaymentsByStatus(ctx context.Context, status []int, limit, offset int64) ([]*Payment, error) {
	q := sqlctx.GetQueryer(ctx)

	rows, err := q.Query(`
		select
			p.id,
			p.image, p.price, p.original_price, p.code,
			p.status, p.created_at, p.at,
			u.id, u.username, u.name, u.email, u.image,
			c.id, c.title, c.image, c.url
		from payments as p
			left join users as u on p.user_id = u.id
			left join courses as c on p.course_id = c.id
		where p.status = any($1)
		order by p.created_at desc
		limit $2 offset $3
	`, pq.Array(status), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var xs []*Payment
	for rows.Next() {
		var x Payment
		err = rows.Scan(
			&x.ID,
			&x.Image, &x.Price, &x.OriginalPrice, &x.Code,
			&x.Status, &x.CreatedAt, pgsql.NullTime(&x.At),
			&x.User.ID, &x.User.Username, &x.User.Name, pgsql.NullString(&x.User.Email), &x.User.Image,
			&x.Course.ID, &x.Course.Title, &x.Course.Image, pgsql.NullString(&x.Course.URL),
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

func countPaymentsByStatus(ctx context.Context, status []int) (cnt int64, err error) {
	q := sqlctx.GetQueryer(ctx)

	err = q.QueryRow(`
		select count(*)
		from payments
		where status = any($1)
	`, pq.Array(status)).Scan(&cnt)
	return
}
