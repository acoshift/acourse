package model

import "github.com/garyburd/redigo/redis"

// User model
type User struct {
	Username string
	Password string
	Name     string
	Email    string
	AboutMe  string
	// CreatedAt time.Time
	// UpdatedAt time.Time
}

// UserGet gets user from id
func UserGet(c redis.Conn, userID int) (*User, error) {
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
	userID, err := redis.Int(c.Do("HGET", key("u", "username"), username))
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
	userID, err := redis.Int(c.Do("HGET", key("u", "email"), email))
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
	userID, err := redis.Int(c.Do("HGET", key("u", "provider", provider), providerUserID))
	if err == redis.ErrNil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return UserGet(c, userID)
}
