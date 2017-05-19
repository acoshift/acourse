package model

import (
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
	At            time.Time

	User   User
	Course Course
}

// PaymentStatus values
const (
	Pending = iota
	Accepted
	Rejected
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
	`

	querySavePayment = `
		insert into payments
			(user_id, course_id, image, price, original_price, code, status, updated_at)
		values
			($1, $2, $3, $4, $5, $6, $7, now())
		returning id
	`

	queryChangePaymentStatus = `
		update payments
		set status = $2
		where id = $1
	`
)

// Save saves payment, allow for create only
func (x *Payment) Save() error {
	if x.ID > 0 {
		return fmt.Errorf("payment already created")
	}
	if len(x.UserID) == 0 {
		return fmt.Errorf("invalid user")
	}
	if x.CourseID <= 0 {
		return fmt.Errorf("invalid course")
	}
	err := db.QueryRow(querySavePayment, x.UserID, x.CourseID, x.Image, x.Price, x.OriginalPrice, x.Code, Pending).Scan(&x.ID)
	if err != nil {
		return err
	}
	return nil
}

// Accept accepts a payment and create new enroll
func (x *Payment) Accept() error {
	if x.ID <= 0 {
		return fmt.Errorf("payment must be save before accept")
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(queryChangePaymentStatus, x.ID, Accepted)
	if err != nil {
		return err
	}
	_, err = tx.Exec(querySaveEnroll, x.UserID, x.CourseID)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// Reject rejects a payment
func (x *Payment) Reject() error {
	if x.ID <= 0 {
		return fmt.Errorf("payment must be save before accept")
	}
	_, err := db.Exec(queryChangePaymentStatus, x.ID, Rejected)
	if err != nil {
		return err
	}
	return nil
}

func scanPayment(scan scanFunc, x *Payment) error {
	var at *time.Time
	var courseURL *string
	err := scan(&x.ID,
		&x.Image, &x.Price, &x.OriginalPrice, &x.Code, &x.Status, &x.CreatedAt, &x.UpdatedAt, &at,
		&x.User.ID, &x.User.Username, &x.User.Image,
		&x.Course.ID, &x.Course.Title, &x.Course.Image, &courseURL,
	)
	if err != nil {
		return err
	}
	if at != nil {
		x.At = *at
	}
	if courseURL != nil {
		x.Course.URL = *courseURL
	}
	x.UserID = x.User.ID
	x.CourseID = x.Course.ID
	return nil
}

// GetPayments gets payments
func GetPayments(paymentIDs []int64) ([]*Payment, error) {
	xs := make([]*Payment, 0, len(paymentIDs))
	rows, err := db.Query(queryGetPayments, pq.Array(paymentIDs))
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var x Payment
		err = scanPayment(rows.Scan, &x)
		if err != nil {
			return nil, err
		}
		xs = append(xs, &x)
	}
	return xs, nil
}

// GetPayment gets payment from given id
func GetPayment(paymentID int64) (*Payment, error) {
	var x Payment
	err := scanPayment(db.QueryRow(queryGetPayment, paymentID).Scan, &x)
	if err != nil {
		return nil, err
	}
	return &x, nil
}

// ListHistoryPayments lists history payments
// TODO: pagination
func ListHistoryPayments() ([]*Payment, error) {
	xs := make([]*Payment, 0)
	rows, err := db.Query(queryListPaymentsWithStatus, pq.Array([]int{Accepted, Rejected}))
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var x Payment
		err = scanPayment(rows.Scan, &x)
		if err != nil {
			return nil, err
		}
		xs = append(xs, &x)
	}
	return xs, nil
}

// ListPendingPayments lists pending payments
// TODO: pagination
func ListPendingPayments() ([]*Payment, error) {
	xs := make([]*Payment, 0)
	rows, err := db.Query(queryListPaymentsWithStatus, pq.Array([]int{Pending}))
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var x Payment
		err = scanPayment(rows.Scan, &x)
		if err != nil {
			return nil, err
		}
		xs = append(xs, &x)
	}
	return xs, nil
}
