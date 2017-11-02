package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/acoshift/acourse/pkg/app"
	"github.com/lib/pq"
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
func (repo) CreatePayment(ctx context.Context, x *app.Payment) error {
	tx := app.GetTransaction(ctx)

	_, err := tx.ExecContext(ctx, `
		insert into payments
			(user_id, course_id, image, price, original_price, code, status)
		values
			($1, $2, $3, $4, $5, $6, $7)
		returning id
	`, x.UserID, x.CourseID, x.Image, x.Price, x.OriginalPrice, x.Code, app.Pending)
	if err != nil {
		return err
	}
	return nil
}

// Accept accepts a payment and create new enroll
func (repo *repo) AcceptPayment(ctx context.Context, x *app.Payment) error {
	tx := app.GetTransaction(ctx)

	if len(x.ID) == 0 {
		return fmt.Errorf("payment must be save before accept")
	}

	_, err := tx.ExecContext(ctx, `
		update payments
		set
			status = $2,
			updated_at = now(),
			at = now()
		where id = $1`, x.ID, app.Accepted)
	if err != nil {
		return err
	}

	err = repo.Enroll(ctx, x.UserID, x.CourseID)
	if err != nil {
		return err
	}

	return nil
}

// Reject rejects a payment
func (repo) RejectPayment(ctx context.Context, x *app.Payment) error {
	db := app.GetDatabase(ctx)

	if len(x.ID) == 0 {
		return fmt.Errorf("payment must be save before accept")
	}
	_, err := db.ExecContext(ctx, `
		update payments
		set
			status = $2,
			updated_at = now(),
			at = now()
		where id = $1
	`, x.ID, app.Rejected)
	if err != nil {
		return err
	}
	return nil
}

func scanPayment(scan scanFunc, x *app.Payment) error {
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
func (repo) GetPayments(ctx context.Context, paymentIDs []string) ([]*app.Payment, error) {
	db := app.GetDatabase(ctx)

	xs := make([]*app.Payment, 0, len(paymentIDs))
	rows, err := db.QueryContext(ctx, queryGetPayments, pq.Array(paymentIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var x app.Payment
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
func (repo) GetPayment(ctx context.Context, paymentID string) (*app.Payment, error) {
	db := app.GetDatabase(ctx)

	var x app.Payment
	err := scanPayment(db.QueryRowContext(ctx, queryGetPayment, paymentID).Scan, &x)
	if err != nil {
		return nil, err
	}
	return &x, nil
}

// HasPendingPayment returns ture if given user has pending payment for given course
func (repo) HasPendingPayment(ctx context.Context, userID string, courseID string) (bool, error) {
	db := app.GetDatabase(ctx)

	var p int
	err := db.QueryRowContext(ctx, `
		select 1 from payments
		where user_id = $1 and course_id = $2 and status = $3`,
		userID, courseID, app.Pending,
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
func (repo) ListHistoryPayments(ctx context.Context, limit, offset int64) ([]*app.Payment, error) {
	db := app.GetDatabase(ctx)

	xs := make([]*app.Payment, 0)
	rows, err := db.QueryContext(ctx, queryListPaymentsWithStatus, pq.Array([]int{app.Accepted, app.Rejected, app.Refunded}), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var x app.Payment
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
func (repo) ListPendingPayments(ctx context.Context, limit, offset int64) ([]*app.Payment, error) {
	db := app.GetDatabase(ctx)

	xs := make([]*app.Payment, 0)
	rows, err := db.QueryContext(ctx, queryListPaymentsWithStatus, pq.Array([]int{app.Pending}), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var x app.Payment
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
func (repo) CountHistoryPayments(ctx context.Context) (int64, error) {
	db := app.GetDatabase(ctx)

	var cnt int64
	err := db.QueryRowContext(ctx, queryCountPaymentsWithStatus, pq.Array([]int{app.Accepted, app.Rejected})).Scan(&cnt)
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

// CountPendingPayments returns pending payments count
func (repo) CountPendingPayments(ctx context.Context) (int64, error) {
	db := app.GetDatabase(ctx)

	var cnt int64
	err := db.QueryRowContext(ctx, queryCountPaymentsWithStatus, pq.Array([]int{app.Pending})).Scan(&cnt)
	if err != nil {
		return 0, err
	}
	return cnt, nil
}
