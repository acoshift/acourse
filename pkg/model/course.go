package model

import (
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
)

// Course model
type Course struct {
	id           string
	option       *CourseOption
	owner        *User
	enrollCount  int
	oldURL       string
	Title        string
	ShortDesc    string
	Desc         string
	Image        string
	UserID       string
	Start        time.Time
	URL          string // MUST not parsable to int
	Type         int
	Price        float64
	Discount     float64
	Contents     []*CourseContent
	EnrollDetail string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Course type values
const (
	_ = iota
	Live
	Video
	EBook
)

// Video type values
const (
	_ = iota
	Youtube
)

// CourseContent type
type CourseContent struct {
	Title       string
	Desc        string
	VideoID     string
	VideoType   int
	DownloadURL string
}

// CourseOption type
type CourseOption struct {
	Public     bool
	Enroll     bool
	Attend     bool
	Assignment bool
	Discount   bool
}

// ID returns course id
func (x *Course) ID() string {
	return x.id
}

// Option returns course option
func (x *Course) Option() *CourseOption {
	if x.option == nil {
		x.option = &CourseOption{}
	}
	return x.option
}

// Owner returns course user if fetched
func (x *Course) Owner() *User {
	return x.owner
}

// EnrollCount returns count of enrolled user
func (x *Course) EnrollCount() int {
	return x.enrollCount
}

// Save saves course
func (x *Course) Save(c redis.Conn) error {
	if len(x.id) == 0 {
		id, err := redis.Int64(c.Do("INCR", key("id", "c")))
		if err != nil {
			return err
		}
		x.id = strconv.FormatInt(id, 10)
	}

	c.Send("MULTI")
	c.Send("SADD", key("c", "all"), x.id)

	x.UpdatedAt = time.Now()
	if x.CreatedAt.IsZero() {
		x.CreatedAt = x.UpdatedAt
		c.Send("ZADD", key("c", "t0"), x.CreatedAt.UnixNano(), x.id)
	} else {
		c.Send("ZADD", key("c", "t0"), x.CreatedAt.UnixNano(), x.id) // TODO: remove after migrate
	}

	c.Send("ZADD", key("c", "t1"), x.UpdatedAt.UnixNano(), x.id)
	c.Send("HSET", key("c"), x.id, enc(x))

	if x.oldURL != x.URL {
		// url updated
		if len(x.oldURL) > 0 {
			c.Send("HDEL", key("c", "url"), x.id)
		}
		if len(x.URL) > 0 {
			c.Send("HSET", key("c", "url"), x.URL, x.id)
		}
	}

	if x.option != nil {
		if x.option.Public {
			c.Send("SADD", key("c", "public"), x.id)
		} else {
			c.Send("SREM", key("c", "public"), x.id)
		}
		if x.option.Enroll {
			c.Send("SADD", key("c", "enroll"), x.id)
		} else {
			c.Send("SREM", key("c", "enroll"), x.id)
		}
		if x.option.Attend {
			c.Send("SADD", key("c", "attend"), x.id)
		} else {
			c.Send("SREM", key("c", "attend"), x.id)
		}
		if x.option.Assignment {
			c.Send("SADD", key("c", "assignment"), x.id)
		} else {
			c.Send("SREM", key("c", "assignment"), x.id)
		}
		if x.option.Discount {
			c.Send("SADD", key("c", "discount"), x.id)
		} else {
			c.Send("SREM", key("c", "discount"), x.id)
		}
	}
	_, err := c.Do("EXEC")
	if err != nil {
		return err
	}

	return nil
}

// GetCourses gets courses
func GetCourses(c redis.Conn, courseIDs []string) ([]*Course, error) {
	xs := make([]*Course, len(courseIDs))
	for _, courseID := range courseIDs {
		c.Send("SISMEMBER", key("c", "all"), courseID)
		c.Send("HGET", key("c"), courseID)
		c.Send("SISMEMBER", key("c", "public"), courseID)
		c.Send("SISMEMBER", key("c", "enroll"), courseID)
		c.Send("SISMEMBER", key("c", "attend"), courseID)
		c.Send("SISMEMBER", key("c", "assignment"), courseID)
		c.Send("SISMEMBER", key("c", "discount"), courseID)
		c.Send("ZCARD", key("c", courseID, "u"))
	}
	c.Flush()
	for i, courseID := range courseIDs {
		exists, _ := redis.Bool(c.Receive())
		if !exists {
			c.Receive()
			c.Receive()
			c.Receive()
			c.Receive()
			c.Receive()
			c.Receive()
			continue
		}
		var x Course
		b, err := redis.Bytes(c.Receive())
		if err != nil {
			return nil, err
		}
		err = dec(b, &x)
		if err != nil {
			return nil, err
		}
		x.oldURL = x.URL
		x.option = &CourseOption{}
		x.option.Public, _ = redis.Bool(c.Receive())
		x.option.Enroll, _ = redis.Bool(c.Receive())
		x.option.Attend, _ = redis.Bool(c.Receive())
		x.option.Assignment, _ = redis.Bool(c.Receive())
		x.option.Discount, _ = redis.Bool(c.Receive())
		x.enrollCount, _ = redis.Int(c.Receive())
		x.id = courseID
		xs[i] = &x
	}
	return xs, nil
}

// GetCourse gets course
func GetCourse(c redis.Conn, courseID string) (*Course, error) {
	xs, err := GetCourses(c, []string{courseID})
	if err != nil {
		return nil, err
	}
	return xs[0], nil
}

// GetCourseFromURL gets course from url
func GetCourseFromURL(c redis.Conn, url string) (*Course, error) {
	userID, err := redis.String(c.Do("HGET", key("c", "url"), url))
	if err == redis.ErrNil {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return GetCourse(c, userID)
}

// GetCourFromIDOrURL gets course from id if given v can parse to int,
// otherwise get from url
func GetCourFromIDOrURL(c redis.Conn, v string) (*Course, error) {
	if _, err := strconv.Atoi(v); err == nil {
		return GetCourse(c, v)
	}
	return GetCourseFromURL(c, v)
}

// ListCourses lists courses
// TODO: pagination
func ListCourses(c redis.Conn) ([]*Course, error) {
	courseIDs, err := redis.Strings(c.Do("ZREVRANGE", key("c", "t0"), 0, -1))
	if err != nil {
		return nil, err
	}
	return GetCourses(c, courseIDs)
}

// ListPublicCourses lists public course sort by created at desc
// TODO: add pagination
func ListPublicCourses(c redis.Conn) ([]*Course, error) {
	c.Send("MULTI")
	c.Send("ZINTERSTORE", key("result"), 2, key("c", "t0"), key("c", "public"), "WEIGHTS", 1, 0)
	c.Send("ZREVRANGE", key("result"), 0, -1)
	reply, err := redis.Values(c.Do("EXEC"))
	if err != nil {
		return nil, err
	}
	courseIDs, _ := redis.Strings(reply[1], nil)
	return GetCourses(c, courseIDs)
}

// ListOwnCourses lists courses that owned by given user
// TODO: add pagination
func ListOwnCourses(c redis.Conn, userID string) ([]*Course, error) {
	c.Send("MULTI")
	c.Send("ZINTERSTORE", key("result"), 2, key("c", "t0"), key("u", userID, "c"), "WEIGHTS", 1, 0)
	c.Send("ZREVRANGE", key("result"), 0, -1)
	reply, err := redis.Values(c.Do("EXEC"))
	if err != nil {
		return nil, err
	}
	courseIDs, _ := redis.Strings(reply[1], nil)
	return GetCourses(c, courseIDs)
}

// ListEnrolledCourses lists courses that enrolled by given user
// TODO: add pagination
func ListEnrolledCourses(c redis.Conn, userID string) ([]*Course, error) {
	courseIDs, err := redis.Strings(c.Do("ZREVRANGE", key("u", userID, "e"), 0, -1))
	if err != nil {
		return nil, err
	}
	return GetCourses(c, courseIDs)
}
