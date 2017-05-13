package model

import (
	"fmt"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
)

// Payment model
type Payment struct {
	id            string
	UserID        string
	CourseID      string
	Image         string
	Price         float64
	OriginalPrice float64
	Code          string
	Status        int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	At            time.Time
}

// PaymentStatus values
const (
	Pending = iota
	Accepted
	Rejected
)

// ID returns payment id
func (x *Payment) ID() string {
	return x.id
}

func (x *Payment) save(c redis.Conn) {
	x.UpdatedAt = time.Now()
	if x.CreatedAt.IsZero() {
		x.CreatedAt = x.UpdatedAt
		n := x.CreatedAt.UnixNano()
		c.Send("ZADD", key("p", "t0"), n, x.id)
		c.Send("ZADD", key("u", x.UserID, "p"), n, x.id)
		c.Send("SADD", key("c", x.CourseID, "p"), n, x.id)
	} else {
		// TODO: remove after migrate
		n := x.CreatedAt.UnixNano()
		c.Send("ZADD", key("p", "t0"), n, x.id)
		c.Send("ZADD", key("u", x.UserID, "p"), n, x.id)
		c.Send("SADD", key("c", x.CourseID, "p"), n, x.id)
	}

	c.Send("ZADD", key("p", "t1"), x.UpdatedAt.UnixNano(), x.id)
	c.Send("HSET", key("p"), x.id, enc(x))
}

// Save saves payment
func (x *Payment) Save(c redis.Conn) error {
	if len(x.UserID) == 0 {
		return fmt.Errorf("invalid user")
	}
	if len(x.CourseID) == 0 {
		return fmt.Errorf("invalid course")
	}

	var err error
	if len(x.id) == 0 {
		id, err := redis.Int64(c.Do("INCR", key("id", "p")))
		if err != nil {
			return err
		}
		x.id = strconv.FormatInt(id, 10)
	}

	c.Send("MULTI")
	c.Send("SADD", key("p", "all"), x.id)
	if x.Status == Pending {
		c.Send("SADD", key("p", "pending"), x.id)
	} else {
		c.Send("SREM", key("p", "pending"), x.id)
	}
	x.save(c)
	_, err = c.Do("EXEC")
	if err != nil {
		return err
	}
	return nil
}

// Accept accepts payment and save
func (x *Payment) Accept(c redis.Conn) error {
	// TODO: accept and reject should run in transaction in `watch`
	if len(x.id) == 0 {
		return fmt.Errorf("invalid payment")
	}
	if x.Status != Pending {
		return fmt.Errorf("invalid payment status")
	}
	x.Status = Accepted
	c.Send("MULTI")
	enroll(c, x.UserID, x.CourseID)
	x.save(c)
	_, err := c.Do("EXEC")
	if err != nil {
		return err
	}
	return nil
}

// Reject rejects a payment
func (x *Payment) Reject(c redis.Conn) error {
	if len(x.id) == 0 {
		return fmt.Errorf("invalid payment")
	}
	if x.Status != Pending {
		return fmt.Errorf("invalid payment status")
	}
	x.Status = Rejected
	c.Send("MULTI")
	x.save(c)
	_, err := c.Do("EXEC")
	if err != nil {
		return err
	}
	return nil
}

// GetPayments gets payments
func GetPayments(c redis.Conn, paymentIDs []string) ([]*Payment, error) {
	xs := make([]*Payment, len(paymentIDs))
	for _, paymentID := range paymentIDs {
		c.Send("SISMEMBER", key("p", "all"), paymentID)
		c.Send("HGET", key("p"), paymentID)
	}
	c.Flush()
	for i, paymentID := range paymentIDs {
		exists, _ := redis.Bool(c.Receive())
		if !exists {
			c.Receive()
			continue
		}
		var x Payment
		b, err := redis.Bytes(c.Receive())
		if err != nil {
			return nil, err
		}
		err = dec(b, &x)
		if err != nil {
			return nil, err
		}
		x.id = paymentID
		xs[i] = &x
	}
	return xs, nil
}

// GetPayment gets payment from given id
func GetPayment(c redis.Conn, paymentID string) (*Payment, error) {
	xs, err := GetPayments(c, []string{paymentID})
	if err != nil {
		return nil, err
	}
	return xs[0], nil
}

// ListPayments lists payments
// TODO: pagination
func ListPayments(c redis.Conn) ([]*Payment, error) {
	paymentIDs, err := redis.Strings(c.Do("ZREVRANGE", key("p", "t0"), 0, -1))
	if err != nil {
		return nil, err
	}
	return GetPayments(c, paymentIDs)
}

// ListPendingPayments lists pending payments
// TODO: pagination
func ListPendingPayments(c redis.Conn) ([]*Payment, error) {
	c.Send("MULTI")
	c.Send("ZINTERSTORE", key("result"), 2, key("p", "t0"), key("p", "pending"), "WEIGHTS", 1, 0)
	c.Send("ZREVRANGE", key("result"), 0, -1)
	reply, err := redis.Values(c.Do("EXEC"))
	if err != nil {
		return nil, err
	}
	paymentIDs, _ := redis.Strings(reply[1], nil)
	return GetPayments(c, paymentIDs)
}
