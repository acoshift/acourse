package model

import (
	"fmt"
	"time"

	"github.com/acoshift/acourse/pkg/internal"
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
		SELECT
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
		FROM payments
			LEFT JOIN users ON payments.user_id = users.id
			LEFT JOIN courses ON payments.course_id = courses.id
	`
)

var (
	getPaymentStmt, _ = internal.GetDB().Prepare(selectPayment + `
		WHERE payments.id = $1;
	`)

	getPaymentsStmt, _ = internal.GetDB().Prepare(selectPayment + `
		WHERE payments.id IN ANY($1);
	`)

	listPaymentsStmt, _ = internal.GetDB().Prepare(selectPayment + `
		ORDER BY payments.created_at DESC;
	`)

	listPaymentsWithStatusStmt, _ = internal.GetDB().Prepare(selectPayment + `
		WHERE payments.status = ANY($1)
		ORDER BY payments.created_at DESC;
	`)
)

func (x *Payment) save() {
	// x.UpdatedAt = time.Now()
	// if x.CreatedAt.IsZero() {
	// 	x.CreatedAt = x.UpdatedAt
	// 	n := x.CreatedAt.UnixNano()
	// 	c.Send("ZADD", key("p", "t0"), n, x.id)
	// 	c.Send("ZADD", key("u", x.UserID, "p"), n, x.id)
	// 	c.Send("SADD", key("c", x.CourseID, "p"), n, x.id)
	// } else {
	// 	// TODO: remove after migrate
	// 	n := x.CreatedAt.UnixNano()
	// 	c.Send("ZADD", key("p", "t0"), n, x.id)
	// 	c.Send("ZADD", key("u", x.UserID, "p"), n, x.id)
	// 	c.Send("SADD", key("c", x.CourseID, "p"), n, x.id)
	// }

	// c.Send("ZADD", key("p", "t1"), x.UpdatedAt.UnixNano(), x.id)
	// c.Send("HSET", key("p"), x.id, enc(x))
}

// Save saves payment
func (x *Payment) Save() error {
	if len(x.UserID) == 0 {
		return fmt.Errorf("invalid user")
	}
	if x.CourseID <= 0 {
		return fmt.Errorf("invalid course")
	}

	// var err error
	// if len(x.id) == 0 {
	// 	id, err := redis.Int64(c.Do("INCR", key("id", "p")))
	// 	if err != nil {
	// 		return err
	// 	}
	// 	x.id = strconv.FormatInt(id, 10)
	// }

	// c.Send("MULTI")
	// c.Send("SADD", key("p", "all"), x.id)
	// if x.Status == Pending {
	// 	c.Send("SADD", key("p", "pending"), x.id)
	// } else {
	// 	c.Send("SREM", key("p", "pending"), x.id)
	// }
	// x.save(c)
	// _, err = c.Do("EXEC")
	// if err != nil {
	// 	return err
	// }
	return nil
}

// Accept accepts payment and save
func (x *Payment) Accept() error {
	// TODO: accept and reject should run in transaction in `watch`
	// if len(x.id) == 0 {
	// 	return fmt.Errorf("invalid payment")
	// }
	// if x.Status != Pending {
	// 	return fmt.Errorf("invalid payment status")
	// }
	// x.Status = Accepted
	// c.Send("MULTI")
	// // enroll(c, x.UserID, x.CourseID)
	// x.save(c)
	// _, err := c.Do("EXEC")
	// if err != nil {
	// 	return err
	// }
	return nil
}

// Reject rejects a payment
func (x *Payment) Reject() error {
	// if len(x.id) == 0 {
	// 	return fmt.Errorf("invalid payment")
	// }
	// if x.Status != Pending {
	// 	return fmt.Errorf("invalid payment status")
	// }
	// x.Status = Rejected
	// c.Send("MULTI")
	// x.save(c)
	// _, err := c.Do("EXEC")
	// if err != nil {
	// 	return err
	// }
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
	return nil
}

// GetPayments gets payments
func GetPayments(paymentIDs []int64) ([]*Payment, error) {
	xs := make([]*Payment, 0, len(paymentIDs))
	rows, err := getPaymentsStmt.Query(pq.Array(paymentIDs))
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
	err := scanPayment(getPaymentStmt.QueryRow(paymentID).Scan, &x)
	if err != nil {
		return nil, err
	}
	return &x, nil
}

// ListHistoryPayments lists history payments
// TODO: pagination
func ListHistoryPayments() ([]*Payment, error) {
	xs := make([]*Payment, 0)
	rows, err := listPaymentsWithStatusStmt.Query(pq.Array([]int{Accepted, Rejected}))
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
	rows, err := listPaymentsWithStatusStmt.Query(pq.Array([]int{Pending}))
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
