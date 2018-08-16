package repository

import (
	"context"
	"database/sql"

	"github.com/acoshift/pgsql"

	"github.com/lib/pq"

	"github.com/acoshift/acourse/context/sqlctx"
	"github.com/acoshift/acourse/entity"
)

const (
	selectPayment = `
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
	`

	queryGetPayment = selectPayment + `
		where payments.id = $1
	`

	queryGetPayments = selectPayment + `
		where payments.id = any($1)
	`

	queryListPayments = selectPayment + `
		order by payments.created_at desc
	`

	queryListPaymentsWithStatus = selectPayment + `
		where payments.status = any($1)
		order by payments.created_at desc
		limit $2 offset $3
	`
)

// CreatePayment creates new payment
func CreatePayment(ctx context.Context, x *entity.Payment) error {
	q := sqlctx.GetQueryer(ctx)

	_, err := q.Exec(`
		insert into payments
			(user_id, course_id, image, price, original_price, code, status)
		values
			($1, $2, $3, $4, $5, $6, $7)
		returning id
	`, x.UserID, x.CourseID, x.Image, x.Price, x.OriginalPrice, x.Code, entity.Pending)
	if err != nil {
		return err
	}
	return nil
}

// SetPaymentStatus sets payment status
func SetPaymentStatus(ctx context.Context, paymentID string, status int) error {
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

func scanPayment(scan scanFunc, x *entity.Payment) error {
	err := scan(&x.ID,
		&x.Image, &x.Price, &x.OriginalPrice, &x.Code, &x.Status, &x.CreatedAt, &x.UpdatedAt, &x.At,
		&x.User.ID, &x.User.Username, &x.User.Name, pgsql.NullString(&x.User.Email), &x.User.Image,
		&x.Course.ID, &x.Course.Title, &x.Course.Image, &x.Course.URL,
	)
	if err != nil {
		return err
	}
	x.UserID = x.User.ID
	x.CourseID = x.Course.ID
	return nil
}

// GetPayments gets payments
func GetPayments(ctx context.Context, paymentIDs []string) ([]*entity.Payment, error) {
	q := sqlctx.GetQueryer(ctx)

	rows, err := q.Query(queryGetPayments, pq.Array(paymentIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	xs := make([]*entity.Payment, 0, len(paymentIDs))
	for rows.Next() {
		var x entity.Payment
		err = scanPayment(rows.Scan, &x)
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

// GetPayment gets payment from given id
func GetPayment(ctx context.Context, paymentID string) (*entity.Payment, error) {
	q := sqlctx.GetQueryer(ctx)

	var x entity.Payment
	err := scanPayment(q.QueryRow(queryGetPayment, paymentID).Scan, &x)
	if err == sql.ErrNoRows {
		return nil, entity.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &x, nil
}

// HasPendingPayment returns ture if given user has pending payment for given course
func HasPendingPayment(ctx context.Context, userID string, courseID string) (exists bool, err error) {
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

// ListPaymentsByStatus lists history payments by statuses
func ListPaymentsByStatus(ctx context.Context, statuses []int, limit, offset int64) ([]*entity.Payment, error) {
	q := sqlctx.GetQueryer(ctx)

	rows, err := q.Query(queryListPaymentsWithStatus, pq.Array(statuses), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var xs []*entity.Payment
	for rows.Next() {
		var x entity.Payment
		err = scanPayment(rows.Scan, &x)
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

// CountPaymentsByStatuses returns payments count by statuses
func CountPaymentsByStatuses(ctx context.Context, statuses []int) (cnt int64, err error) {
	q := sqlctx.GetQueryer(ctx)

	err = q.QueryRow(`
		select count(*)
		from payments
		where status = any($1)
	`, pq.Array(statuses)).Scan(&cnt)
	return
}
