package model

import "time"

// User model
type User struct {
	ID          string
	role        *UserRole
	oldUsername string
	Username    string
	Name        string
	Email       string
	AboutMe     string
	Image       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// UserRole type
type UserRole struct {
	Admin      bool
	Instructor bool
}

// Role returns user role
func (x *User) Role() *UserRole {
	if x.role == nil {
		x.role = &UserRole{}
	}
	return x.role
}

// Save saves user
// func (x *User) Save(c redis.Conn) error {
// 	if len(x.ID) == 0 {
// 		return fmt.Errorf("invalid id")
// 	}

// 	c.Do("WATCH", key("u", "username"))
// 	// verify is new username duplicate
// 	{
// 		uID, err := redis.String(c.Do("HGET", key("u", "username"), x.Username))
// 		if err != redis.ErrNil && err != nil {
// 			c.Do("UNWATCH")
// 			return err
// 		}
// 		if len(uID) > 0 && x.ID != uID {
// 			c.Do("UNWATCH")
// 			return fmt.Errorf("username already exists")
// 		}
// 	}

// 	c.Send("MULTI")
// 	c.Send("SADD", key("u", "all"), x.id)

// 	x.UpdatedAt = time.Now()
// 	if x.CreatedAt.IsZero() {
// 		x.CreatedAt = x.UpdatedAt
// 		c.Send("ZADD", key("u", "t0"), x.CreatedAt.UnixNano(), x.id)
// 	} else {
// 		c.Send("ZADD", key("u", "t0"), x.CreatedAt.UnixNano(), x.id) // TODO: remove after migrate
// 	}

// 	c.Send("ZADD", key("u", "t1"), x.UpdatedAt.UnixNano(), x.id)
// 	c.Send("HSET", key("u"), x.id, enc(x))

// 	if x.role != nil {
// 		if x.role.Admin {
// 			c.Send("SADD", key("u", "admin"), x.id)
// 		} else {
// 			c.Send("SREM", key("u", "admin"), x.id)
// 		}
// 		if x.role.Instructor {
// 			c.Send("SADD", key("u", "instructor"), x.id)
// 		} else {
// 			c.Send("SREM", key("u", "instructor"), x.id)
// 		}
// 	}
// 	if x.oldUsername != x.Username {
// 		if len(x.oldUsername) > 0 {
// 			c.Send("HDEL", key("u", "username"), x.oldUsername)
// 		}
// 		if len(x.Username) > 0 {
// 			c.Send("HSET", key("u", "username"), x.Username, x.id)
// 		}
// 	}

// 	_, err := c.Do("EXEC")
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// // GetUsers gets users
// func GetUsers(c redis.Conn, userIDs []string) ([]*User, error) {
// 	xs := make([]*User, len(userIDs))
// 	for _, userID := range userIDs {
// 		c.Send("SISMEMBER", key("u", "all"), userID)
// 		c.Send("HGET", key("u"), userID)
// 		c.Send("SISMEMBER", key("u", "admin"), userID)
// 		c.Send("SISMEMBER", key("u", "instructor"), userID)
// 	}
// 	c.Flush()
// 	for i, userID := range userIDs {
// 		exists, _ := redis.Bool(c.Receive())
// 		if !exists {
// 			c.Receive()
// 			c.Receive()
// 			c.Receive()
// 			continue
// 		}
// 		var x User
// 		b, err := redis.Bytes(c.Receive())
// 		if err != nil {
// 			return nil, err
// 		}
// 		err = dec(b, &x)
// 		if err != nil {
// 			return nil, err
// 		}
// 		x.role = &UserRole{}
// 		x.role.Admin, _ = redis.Bool(c.Receive())
// 		x.role.Instructor, _ = redis.Bool(c.Receive())
// 		x.oldUsername = x.Username
// 		x.id = userID
// 		xs[i] = &x
// 	}
// 	return xs, nil
// }

// // GetUser gets user from id
// func GetUser(c redis.Conn, userID string) (*User, error) {
// 	xs, err := GetUsers(c, []string{userID})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return xs[0], nil
// }

// // GetUserFromUsername gets user from username
// func GetUserFromUsername(c redis.Conn, username string) (*User, error) {
// 	userID, err := redis.String(c.Do("HGET", key("u", "username"), username))
// 	if err == redis.ErrNil {
// 		return nil, ErrNotFound
// 	}
// 	if err != nil {
// 		return nil, err
// 	}
// 	return GetUser(c, userID)
// }

// // ListUsers lists users
// // TODO: pagination
// func ListUsers(c redis.Conn) ([]*User, error) {
// 	userIDs, err := redis.Strings(c.Do("ZREVRANGE", key("u", "t0"), 0, -1))
// 	if err != nil {
// 		return nil, err
// 	}
// 	return GetUsers(c, userIDs)
// }
