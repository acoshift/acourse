package model

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

// User model
type User struct {
	id        string
	Username  string
	Password  string
	Name      string
	Email     string
	AboutMe   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ID returns user id
func (x *User) ID() string {
	return x.id
}

// UserGet gets user from id
func UserGet(c redis.Conn, userID string) (*User, error) {
	var x User
	b, err := redis.Bytes(c.Do("HGET", key("u"), userID))
	if err != nil {
		return nil, err
	}
	err = dec(b, &x)
	if err != nil {
		return nil, err
	}
	return &x, nil
}

// UserGetFromUsername gets user from username
func UserGetFromUsername(c redis.Conn, username string) (*User, error) {
	userID, err := redis.String(c.Do("HGET", key("u", "username"), username))
	if err == redis.ErrNil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return UserGet(c, userID)
}

// UserGetFromEmail gets user from email
func UserGetFromEmail(c redis.Conn, email string) (*User, error) {
	userID, err := redis.String(c.Do("HGET", key("u", "email"), email))
	if err == redis.ErrNil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return UserGet(c, userID)
}

// UserGetFromProvider gets user from provider
func UserGetFromProvider(c redis.Conn, provider string, providerUserID string) (*User, error) {
	userID, err := redis.String(c.Do("HGET", key("u", "provider", provider), providerUserID))
	if err == redis.ErrNil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return UserGet(c, userID)
}

// UserSave saves user
func UserSave(c redis.Conn, x *User) error {
	var err error
	if len(x.id) == 0 {
		x.id, err = redis.String(c.Do("INCR", key("id", "u")))
		if err != nil {
			return err
		}
	}

	c.Send("MULTI")

	x.UpdatedAt = time.Now()
	if x.CreatedAt.IsZero() {
		x.CreatedAt = x.UpdatedAt
		c.Send("ZADD", key("u", "t0"), x.CreatedAt.UnixNano(), x.id)
	}

	c.Send("HSET", key("u"), x.id, enc(x))
	c.Send("ZADD", key("u", "t1"), x.UpdatedAt.UnixNano(), x.id)
	_, err = c.Do("EXEC")
	if err != nil {
		return err
	}
	return nil
}
