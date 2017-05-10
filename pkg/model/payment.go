package model

import (
	"fmt"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
)

// Payment model
type Payment struct {
	id        string
	UserID    string
	CourseID  string
	Image     string
	Status    PaymentStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

// PaymentStatus type
type PaymentStatus int

// PaymentStatus values
const (
	Pending PaymentStatus = iota
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
