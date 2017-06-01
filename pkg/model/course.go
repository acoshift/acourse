package model

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/lib/pq"
)

// Course model
type Course struct {
	ID            int64
	Option        CourseOption
	Owner         User
	EnrollCount   int64
	Title         string
	ShortDesc     string
	Desc          string
	Image         string
	UserID        string
	Start         pq.NullTime
	URL           sql.NullString // MUST not parsable to int
	Type          int
	Price         float64
	Discount      float64
	Contents      []*CourseContent
	EnrollDetail  string
	AssignmentIDs []string
	CreatedAt     time.Time
	UpdatedAt     time.Time
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

// Link returns id if url is invalid
func (x *Course) Link() string {
	if !x.URL.Valid || len(x.URL.String) == 0 {
		return strconv.FormatInt(x.ID, 10)
	}
	return x.URL.String
}

const (
	selectCourses = `
		select
			courses.id,
			courses.title,
			courses.short_desc,
			courses.long_desc,
			courses.image,
			courses.start,
			courses.url,
			courses.type,
			courses.price,
			courses.discount,
			courses.enroll_detail,
			courses.created_at,
			courses.updated_at,
			course_options.public,
			course_options.enroll,
			course_options.attend,
			course_options.assignment,
			course_options.discount
		from courses
			left join course_options on courses.id = course_options.course_id
	`

	queryGetCourses = selectCourses + `
		where courses.id = any($1)
	`

	queryListCoursesPublic = selectCourses + `
		where course_options.public = true
		order by courses.created_at desc
	`

	queryListCoursesOwn = selectCourses + `
		where courses.user_id = $1
		order by courses.created_at desc
	`

	queryListCoursesEnrolled = selectCourses + `
		inner join enrolls on courses.id = enrolls.course_id
		where enrolls.user_id = $1
		order by enrolls.created_at desc
	`
)

// Save saves course
func (x *Course) Save() error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	strID := strconv.FormatInt(x.ID, 10)
	if len(x.URL.String) > 0 && x.URL.String != strID {
		x.URL.Valid = true
	} else {
		x.URL.String = strID
		x.URL.Valid = false
	}

	_, err = tx.Exec(`
		upsert into courses
			(id, user_id, title, short_desc, long_desc, image, start, url, type, price, discount, enroll_detail, updated_at)
		values
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, now())
	`, x.ID, x.UserID, x.Title, x.ShortDesc, x.Desc, x.Image, x.Start, x.URL, x.Type, x.Price, x.Discount, x.EnrollDetail)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
		upsert into course_options
			(course_id, public, enroll, attend, assignment, discount)
		values
			($1, $2, $3, $4, $5, $6)
	`, x.ID, x.Option.Public, x.Option.Enroll, x.Option.Attend, x.Option.Assignment, x.Option.Discount)
	if err != nil {
		return err
	}
	// TODO: save contents
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func scanCourse(scan scanFunc, x *Course) error {
	err := scan(&x.ID,
		&x.Title, &x.ShortDesc, &x.Desc, &x.Image, &x.Start, &x.URL, &x.Type, &x.Price, &x.Discount, &x.EnrollDetail,
		&x.CreatedAt, &x.UpdatedAt,
		&x.Option.Public, &x.Option.Enroll, &x.Option.Attend, &x.Option.Assignment, &x.Option.Discount,
	)
	if err != nil {
		return err
	}
	if len(x.URL.String) == 0 {
		x.URL.String = strconv.FormatInt(x.ID, 10)
	}
	return nil
}

// GetCourses gets courses
func GetCourses(courseIDs []int64) ([]*Course, error) {
	xs := make([]*Course, 0, len(courseIDs))
	rows, err := db.Query(queryGetCourses, pq.Array(courseIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var x Course
		err = scanCourse(rows.Scan, &x)
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

// GetCourse gets course
func GetCourse(courseID int64) (*Course, error) {
	var x Course
	err := db.QueryRow(`
		select
			id, user_id, title, short_desc, long_desc, image, start, url, type, price, courses.discount, enroll_detail,
			opt.public, opt.enroll, opt.attend, opt.assignment, opt.discount
		from courses left join course_options as opt on courses.id = opt.course_id
		where id = $1
	`, courseID).Scan(
		&x.ID, &x.UserID, &x.Title, &x.ShortDesc, &x.Desc, &x.Image, &x.Start, &x.URL, &x.Type, &x.Price, &x.Discount, &x.EnrollDetail,
		&x.Option.Public, &x.Option.Enroll, &x.Option.Attend, &x.Option.Assignment, &x.Option.Discount,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &x, nil
}

// GetCourseContents gets course contents for given course id
func GetCourseContents(courseID int64) ([]*CourseContent, error) {
	rows, err := db.Query(`
		select
			title,
			long_desc,
			video_id,
			video_type,
			download_url
		from course_contents
		where course_id = $1
		order by i asc
	`, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	xs := make([]*CourseContent, 0)
	for rows.Next() {
		var x CourseContent
		err = rows.Scan(&x.Title, &x.Desc, &x.VideoID, &x.VideoType, &x.DownloadURL)
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

// GetCourseIDFromURL gets course id from url
func GetCourseIDFromURL(url string) (int64, error) {
	var id int64
	err := db.QueryRow(`
		select id
		from courses
		where url = $1
	`, url).Scan(&id)
	if err == sql.ErrNoRows {
		return 0, ErrNotFound
	}
	if err != nil {
		return 0, err
	}
	return id, nil
}

// ListCourses lists all courses
// TODO: pagination
func ListCourses() ([]*Course, error) {
	xs := make([]*Course, 0)
	rows, err := db.Query(`
		select
			courses.id,
			courses.title,
			courses.short_desc,
			courses.long_desc,
			courses.image,
			courses.start,
			courses.url,
			courses.type,
			courses.price,
			courses.discount,
			courses.enroll_detail,
			courses.created_at,
			courses.updated_at,
			course_options.public,
			course_options.enroll,
			course_options.attend,
			course_options.assignment,
			course_options.discount,
			users.id,
			users.username,
			users.image
		from courses
			left join course_options on courses.id = course_options.course_id
			left join users on courses.user_id = users.id
			order by courses.created_at desc
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var x Course
		err := rows.Scan(&x.ID,
			&x.Title, &x.ShortDesc, &x.Desc, &x.Image, &x.Start, &x.URL, &x.Type, &x.Price, &x.Discount, &x.EnrollDetail,
			&x.CreatedAt, &x.UpdatedAt,
			&x.Option.Public, &x.Option.Enroll, &x.Option.Attend, &x.Option.Assignment, &x.Option.Discount,
			&x.Owner.ID, &x.Owner.Username, &x.Owner.Image,
		)
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

// ListPublicCourses lists public course sort by created at desc
// TODO: add pagination
func ListPublicCourses() ([]*Course, error) {
	rows, err := db.Query(queryListCoursesPublic)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	xs := make([]*Course, 0)
	ids := make([]int64, 0)
	m := make(map[int64]*Course)
	for rows.Next() {
		var x Course
		err = scanCourse(rows.Scan, &x)
		if err != nil {
			return nil, err
		}
		xs = append(xs, &x)
		ids = append(ids, x.ID)
		m[x.ID] = &x
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	rows.Close()

	rows, err = db.Query(`select course_id, count(*) from enrolls where course_id = any($1) group by course_id`, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var courseID, cnt int64
		err = rows.Scan(&courseID, &cnt)
		if err != nil {
			return nil, err
		}
		m[courseID].EnrollCount = cnt
	}
	return xs, nil
}

// ListOwnCourses lists courses that owned by given user
// TODO: add pagination
func ListOwnCourses(userID string) ([]*Course, error) {
	rows, err := db.Query(queryListCoursesOwn, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	xs := make([]*Course, 0)
	ids := make([]int64, 0)
	m := make(map[int64]*Course)
	for rows.Next() {
		var x Course
		err = scanCourse(rows.Scan, &x)
		if err != nil {
			return nil, err
		}
		xs = append(xs, &x)
		ids = append(ids, x.ID)
		m[x.ID] = &x
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	rows.Close()

	rows, err = db.Query(`select course_id, count(*) from enrolls where course_id = any($1) group by course_id`, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var courseID, cnt int64
		err = rows.Scan(&courseID, &cnt)
		if err != nil {
			return nil, err
		}
		m[courseID].EnrollCount = cnt
	}
	return xs, nil
}

// ListEnrolledCourses lists courses that enrolled by given user
// TODO: add pagination
func ListEnrolledCourses(userID string) ([]*Course, error) {
	xs := make([]*Course, 0)
	rows, err := db.Query(queryListCoursesEnrolled, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var x Course
		err = scanCourse(rows.Scan, &x)
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
