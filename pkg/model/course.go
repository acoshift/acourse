package model

import (
	"strconv"
	"time"

	"github.com/acoshift/acourse/pkg/internal"
	"github.com/lib/pq"
)

// Course model
type Course struct {
	ID            int64
	Option        CourseOption
	Owner         User
	enrollCount   int
	Title         string
	ShortDesc     string
	Desc          string
	Image         string
	UserID        string
	Start         time.Time
	URL           string // MUST not parsable to int
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
	ID          int64
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

// EnrollCount returns count of enrolled user
func (x *Course) EnrollCount() int {
	return x.enrollCount
}

const (
	selectCourses = `
		SELECT
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
		FROM courses
			LEFT JOIN course_options ON courses.id = course_options.id
`
)

var (
	getCourseStmt, _ = internal.GetDB().Prepare(selectCourses + `
		WHERE courses.id = $1;
	`)

	getCoursesStmt, _ = internal.GetDB().Prepare(selectCourses + `
		WHERE courses.id = ANY($1);
	`)

	getCourseFromURLStmt, _ = internal.GetDB().Prepare(selectCourses + `
		WHERE courses.url = $1;
	`)

	getCourseContentsStmt, _ = internal.GetDB().Prepare(`
		SELECT
			courses.id,
			course_contents.id,
			course_contents.title,
			course_contents.long_desc,
			course_contents.video_id,
			course_contents.video_type,
			course_contents.download_url
		FROM courses
			INNER JOIN course_contents
			ON courses.id = course_contents.course_id,
		WHERE courses.id = ANY($1);
	`)

	listCoursesStmt, _ = internal.GetDB().Prepare(selectCourses + `
		ORDER BY courses.created_at DESC;
	`)

	listCoursesPublicStmt, _ = internal.GetDB().Prepare(selectCourses + `
		WHERE course_options.public = true
		ORDER BY courses.created_at DESC;
	`)

	listCoursesOwnStmt, _ = internal.GetDB().Prepare(selectCourses + `
		WHERE courses.user_id = $1
		ORDER BY courses.created_at DESC;
	`)

	listCoursesEnrolledStmt, _ = internal.GetDB().Prepare(selectCourses + `
		INNER JOIN enrolls ON courses.id = enrolls.course_id
		WHERE enrolls.user_id = $1
		ORDER BY enrolls.created_at DESC;
	`)

	saveCourseStmt, _ = internal.GetDB().Prepare(`
		UPSERT INTO courses
			(id, user_id, title, short_desc, long_desc, image, start, url, type, price, discount, enroll_detail, updated_at)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, now());
	`)

	saveCourseOptionStmt, _ = internal.GetDB().Prepare(`
		UPSERT INTO course_options
			(id, public, enroll, attend, assignment, discount)
		VALUES
			($1, $2, $3, $4, $5, %6);
	`)
)

// Save saves course
func (x *Course) Save() error {
	tx, err := internal.GetDB().Begin()
	if err != nil {
		return err
	}
	var start *time.Time
	if !x.Start.IsZero() {
		start = &x.Start
	}
	var url *string
	if len(x.URL) > 0 && x.URL != strconv.FormatInt(x.ID, 10) {
		url = &x.URL
	}

	tx.Stmt(saveCourseStmt).Exec(x.ID, x.UserID, x.Title, x.ShortDesc, x.Desc, x.Image, start, url, x.Type, x.Price, x.Discount, x.EnrollDetail)
	tx.Stmt(saveCourseOptionStmt).Exec(x.ID, x.Option.Public, x.Option.Enroll, x.Option.Attend, x.Option.Assignment, x.Option.Discount)
	// TODO: save contents
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func scanCourse(scan scanFunc, x *Course) error {
	var start *time.Time
	var u *string
	err := scan(&x.ID,
		&x.Title, &x.ShortDesc, &x.Desc, &x.Image, &start, &u, &x.Type, &x.Price, &x.Discount, &x.EnrollDetail,
		&x.CreatedAt, &x.UpdatedAt,
		&x.Option.Public, &x.Option.Enroll, &x.Option.Attend, &x.Option.Assignment, &x.Option.Discount,
	)
	if err != nil {
		return err
	}
	if start != nil {
		x.Start = *start
	}
	if u != nil {
		x.URL = *u
	}
	if len(x.URL) == 0 {
		x.URL = strconv.FormatInt(x.ID, 10)
	}
	return nil
}

func scanCourseContent(scan scanFunc, courseID *int64, x *CourseContent) error {
	return scan(courseID, &x.ID, &x.Title, &x.Desc, &x.VideoID, &x.VideoType, &x.DownloadURL)
}

// GetCourses gets courses
func GetCourses(courseIDs []int64) ([]*Course, error) {
	xs := make([]*Course, 0, len(courseIDs))
	rows, err := getCoursesStmt.Query(pq.Array(courseIDs))
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var x Course
		err = scanCourse(rows.Scan, &x)
		if err != nil {
			return nil, err
		}
		xs = append(xs, &x)
	}
	return xs, nil
}

// GetCourse gets course
func GetCourse(courseID int64) (*Course, error) {
	var x Course
	err := scanCourse(getCourseStmt.QueryRow(courseID).Scan, &x)
	if err != nil {
		return nil, err
	}
	return &x, nil
}

// GetCourseFromURL gets course from url
func GetCourseFromURL(url string) (*Course, error) {
	var x Course
	err := scanCourse(getCourseFromURLStmt.QueryRow(url).Scan, &x)
	if err != nil {
		return nil, err
	}
	return &x, nil
}

// GetCourFromIDOrURL gets course from id if given v can parse to int,
// otherwise get from url
func GetCourFromIDOrURL(v string) (*Course, error) {
	if id, err := strconv.ParseInt(v, 10, 64); err == nil {
		return GetCourse(id)
	}
	return GetCourseFromURL(v)
}

// ListCourses lists courses
// TODO: pagination
func ListCourses() ([]*Course, error) {
	xs := make([]*Course, 0)
	rows, err := listCoursesStmt.Query()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var x Course
		err = scanCourse(rows.Scan, &x)
		if err != nil {
			return nil, err
		}
		xs = append(xs, &x)
	}
	return xs, nil
}

// ListPublicCourses lists public course sort by created at desc
// TODO: add pagination
func ListPublicCourses() ([]*Course, error) {
	xs := make([]*Course, 0)
	rows, err := listCoursesPublicStmt.Query()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var x Course
		err = scanCourse(rows.Scan, &x)
		if err != nil {
			return nil, err
		}
		xs = append(xs, &x)
	}
	return xs, nil
}

// ListOwnCourses lists courses that owned by given user
// TODO: add pagination
func ListOwnCourses(userID string) ([]*Course, error) {
	xs := make([]*Course, 0)
	rows, err := listCoursesOwnStmt.Query(userID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var x Course
		err = scanCourse(rows.Scan, &x)
		if err != nil {
			return nil, err
		}
		xs = append(xs, &x)
	}
	return xs, nil
}

// ListEnrolledCourses lists courses that enrolled by given user
// TODO: add pagination
func ListEnrolledCourses(userID string) ([]*Course, error) {
	xs := make([]*Course, 0)
	rows, err := listCoursesEnrolledStmt.Query(userID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var x Course
		err = scanCourse(rows.Scan, &x)
		if err != nil {
			return nil, err
		}
		xs = append(xs, &x)
	}
	return xs, nil
}
