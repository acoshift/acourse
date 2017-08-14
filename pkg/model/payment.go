package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
)

// Payment model
type Payment struct {
	ID            int64
	UserID        string
	CourseID      int64
	Image         string
	Price         float64
	OriginalPrice float64
	Code          string
	Status        int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	At            pq.NullTime

	User   User
	Course Course
}

// PaymentStatus values
const (
	Pending = iota
	Accepted
	Rejected
	Refunded
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
func CreatePayment(ctx context.Context, tx *sql.Tx, x *Payment) error {
	_, err := tx.ExecContext(ctx, `
		insert into payments
			(user_id, course_id, image, price, original_price, code, status)
		values
			($1, $2, $3, $4, $5, $6, $7)
		returning id
	`, x.UserID, x.CourseID, x.Image, x.Price, x.OriginalPrice, x.Code, Pending)
	if err != nil {
		return err
	}
	return nil
}

// Accept accepts a payment and create new enroll
func (x *Payment) Accept(ctx context.Context, tx *sql.Tx) error {
	if x.ID == 0 {
		return fmt.Errorf("payment must be save before accept")
	}

	_, err := tx.Exec(`
		update payments
		set
			status = $2,
			updated_at = now(),
			at = now()
		where id = $1`, x.ID, Accepted)
	if err != nil {
		return err
	}

	err = Enroll(ctx, tx, x.UserID, x.CourseID)
	if err != nil {
		return err
	}

	return nil
}

// Reject rejects a payment
func (x *Payment) Reject(ctx context.Context, db DB) error {
	if x.ID == 0 {
		return fmt.Errorf("payment must be save before accept")
	}
	_, err := db.ExecContext(ctx, `
		update payments
		set
			status = $2,
			updated_at = now(),
			at = now()
		where id = $1
	`, x.ID, Rejected)
	if err != nil {
		return err
	}
	return nil
}

func scanPayment(scan scanFunc, x *Payment) error {
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
func GetPayments(ctx context.Context, db DB, paymentIDs []string) ([]*Payment, error) {
	xs := make([]*Payment, 0, len(paymentIDs))
	rows, err := db.QueryContext(ctx, queryGetPayments, pq.Array(paymentIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var x Payment
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
func GetPayment(ctx context.Context, db DB, paymentID int64) (*Payment, error) {
	var x Payment
	err := scanPayment(db.QueryRowContext(ctx, queryGetPayment, paymentID).Scan, &x)
	if err != nil {
		return nil, err
	}
	return &x, nil
}

// HasPendingPayment returns ture if given user has pending payment for given course
func HasPendingPayment(ctx context.Context, db DB, userID string, courseID int64) (bool, error) {
	var p int
	err := db.QueryRowContext(ctx, `
		select 1 from payments
		where user_id = $1 and course_id = $2 and status = $3`,
		userID, courseID, Pending,
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
func ListHistoryPayments(ctx context.Context, db DB, limit, offset int64) ([]*Payment, error) {
	xs := make([]*Payment, 0)
	rows, err := db.QueryContext(ctx, queryListPaymentsWithStatus, pq.Array([]int{Accepted, Rejected, Refunded}), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var x Payment
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
func ListPendingPayments(ctx context.Context, db DB, limit, offset int64) ([]*Payment, error) {
	xs := make([]*Payment, 0)
	rows, err := db.QueryContext(ctx, queryListPaymentsWithStatus, pq.Array([]int{Pending}), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var x Payment
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
func CountHistoryPayments(ctx context.Context, db DB) (int64, error) {
	var cnt int64
	err := db.QueryRowContext(ctx, queryCountPaymentsWithStatus, pq.Array([]int{Accepted, Rejected})).Scan(&cnt)
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

// CountPendingPayments returns pending payments count
func CountPendingPayments(ctx context.Context, db DB) (int64, error) {
	var cnt int64
	err := db.QueryRowContext(ctx, queryCountPaymentsWithStatus, pq.Array([]int{Pending})).Scan(&cnt)
	if err != nil {
		return 0, err
	}
	return cnt, nil
}
