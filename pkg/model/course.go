package model

import (
	"time"

	"strconv"

	"github.com/garyburd/redigo/redis"
)

// Course model
type Course struct {
	id           string
	option       *CourseOption
	studentCount int
	oldURL       string
	Title        string
	ShortDesc    string
	Desc         string
	Image        string
	UserID       string
	Start        time.Time
	URL          string
	Type         CourseType
	Price        float64
	Discount     float64
	Contents     []*CourseContent
	EnrollDetail string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// CourseType type
type CourseType int

// CourseType values
const (
	_ CourseType = iota
	Live
	Video
	EBook
)

// CourseContent type
type CourseContent struct {
	Title       string
	Desc        string
	YoutubeID   string
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
	return x.option
}

// StudentCount returns student count
func (x *Course) StudentCount() int {
	return x.studentCount
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
	}

	c.Send("ZADD", key("c", "t1"), x.UpdatedAt.UnixNano(), x.id)
	c.Send("HSET", key("c"), x.id, enc(x))

	if x.oldURL != x.URL {
		// url updated
		if len(x.oldURL) > 0 {
			c.Send("HDEL", key("c", "url"), x.id)
		}
		if len(x.URL) > 0 {
			c.Send("HSET", key("c", "url"), x.id, x.URL)
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
		c.Send("ZCARD", key("c", courseID, "u"), courseID)
	}
	c.Flush()
	for i := range courseIDs {
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
		x.studentCount, _ = redis.Int(c.Receive())
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
