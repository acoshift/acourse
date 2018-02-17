package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"

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

	queryCountPaymentsWithStatus = `
		select count(*)
		from payments
		where status = any($1)
	`
)

// CreatePayment creates new payment
func CreatePayment(ctx context.Context, q Queryer, x *entity.Payment) error {
	_, err := q.ExecContext(ctx, `
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

// AcceptPayment accepts a payment and create new enroll
func AcceptPayment(ctx context.Context, q Queryer, x *entity.Payment) error {
	if len(x.ID) == 0 {
		return fmt.Errorf("payment must be save before accept")
	}

	_, err := q.ExecContext(ctx, `
		update payments
		set
			status = $2,
			updated_at = now(),
			at = now()
		where id = $1`, x.ID, entity.Accepted)
	if err != nil {
		return err
	}

	err = Enroll(ctx, q, x.UserID, x.CourseID)
	if err != nil {
		return err
	}

	return nil
}

// RejectPayment rejects a payment
func RejectPayment(ctx context.Context, q Queryer, x *entity.Payment) error {
	if len(x.ID) == 0 {
		return fmt.Errorf("payment must be save before accept")
	}
	_, err := q.ExecContext(ctx, `
		update payments
		set
			status = $2,
			updated_at = now(),
			at = now()
		where id = $1
	`, x.ID, entity.Rejected)
	if err != nil {
		return err
	}
	return nil
}

func scanPayment(scan scanFunc, x *entity.Payment) error {
	err := scan(&x.ID,
		&x.Image, &x.Price, &x.OriginalPrice, &x.Code, &x.Status, &x.CreatedAt, &x.UpdatedAt, &x.At,
		&x.User.ID, &x.User.Username, &x.User.Name, &x.User.Email, &x.User.Image,
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
func GetPayments(ctx context.Context, q Queryer, paymentIDs []string) ([]*entity.Payment, error) {
	xs := make([]*entity.Payment, 0, len(paymentIDs))
	rows, err := q.QueryContext(ctx, queryGetPayments, pq.Array(paymentIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
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
func GetPayment(ctx context.Context, q Queryer, paymentID string) (*entity.Payment, error) {
	var x entity.Payment
	err := scanPayment(q.QueryRowContext(ctx, queryGetPayment, paymentID).Scan, &x)
	if err != nil {
		return nil, err
	}
	return &x, nil
}

// HasPendingPayment returns ture if given user has pending payment for given course
func HasPendingPayment(ctx context.Context, q Queryer, userID string, courseID string) (bool, error) {
	var p int
	err := q.QueryRowContext(ctx, `
		select 1 from payments
		where user_id = $1 and course_id = $2 and status = $3`,
		userID, courseID, entity.Pending,
	).Scan(&p)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// ListHistoryPayments lists history payments
func ListHistoryPayments(ctx context.Context, q Queryer, limit, offset int64) ([]*entity.Payment, error) {
	xs := make([]*entity.Payment, 0)
	rows, err := q.QueryContext(ctx, queryListPaymentsWithStatus, pq.Array([]int{entity.Accepted, entity.Rejected, entity.Refunded}), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
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

// ListPendingPayments lists pending payments
func ListPendingPayments(ctx context.Context, q Queryer, limit, offset int64) ([]*entity.Payment, error) {
	xs := make([]*entity.Payment, 0)
	rows, err := q.QueryContext(ctx, queryListPaymentsWithStatus, pq.Array([]int{entity.Pending}), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
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

// CountHistoryPayments returns history payments count
func CountHistoryPayments(ctx context.Context, q Queryer) (int64, error) {
	var cnt int64
	err := q.QueryRowContext(ctx, queryCountPaymentsWithStatus, pq.Array([]int{entity.Accepted, entity.Rejected})).Scan(&cnt)
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

// CountPendingPayments returns pending payments count
func CountPendingPayments(ctx context.Context, q Queryer) (int64, error) {
	var cnt int64
	err := q.QueryRowContext(ctx, queryCountPaymentsWithStatus, pq.Array([]int{entity.Pending})).Scan(&cnt)
	if err != nil {
		return 0, err
	}
	return cnt, nil
}
