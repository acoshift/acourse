package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/configfile"
	"github.com/acoshift/ds"
	"github.com/garyburd/redigo/redis"
	_ "github.com/lib/pq"
	"golang.org/x/oauth2/google"
	identitytoolkit "google.golang.org/api/identitytoolkit/v3"
)

var ctx = context.Background()

var client, _ = ds.NewClient(ctx, "acourse-156413")
var conn, _ = redis.Dial("tcp", "localhost:6379", redis.DialDatabase(4))

var config = configfile.NewReader("../config")
var serviceAccount = config.Bytes("service_account")
var sqlURL = config.String("sql_url")

// migrate all data from datastore and firebase to redis
func main() {
	time.Local = time.UTC
	ctx := context.Background()

	gconf, err := google.JWTConfigFromJSON(serviceAccount, identitytoolkit.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}
	gitService, err := identitytoolkit.New(gconf.Client(ctx))
	if err != nil {
		log.Fatal(err)
	}
	gitClient := gitService.Relyingparty

	db, err := sql.Open("postgres", sqlURL)
	must(err)

	var users []*userModel
	var roles []*roleModel
	var courses []*courseModel
	var payments []*paymentModel
	var attends []*attendModel
	var enrolls []*enrollModel
	var assignments []*assignment
	var userAssignments []*userAssignment

	log.Println("load old database")
	must(client.Query(ctx, "User", &users))
	must(client.Query(ctx, "Role", &roles))
	must(client.Query(ctx, "Course", &courses))
	must(client.Query(ctx, "Payment", &payments))
	must(client.Query(ctx, "Attend", &attends))
	must(client.Query(ctx, "Enroll", &enrolls))
	must(client.Query(ctx, "Assignment", &assignments))
	must(client.Query(ctx, "UserAssignment", &userAssignments))

	findCourse := func(courseID string) *courseModel {
		for _, p := range courses {
			if p.ID() == courseID {
				return p
			}
		}
		return nil
	}

	findAssignment := func(assignmentID string) *assignment {
		for _, p := range assignments {
			if p.ID() == assignmentID {
				return p
			}
		}
		return nil
	}

	respUser, err := gitClient.DownloadAccount(&identitytoolkit.IdentitytoolkitRelyingpartyDownloadAccountRequest{
		MaxResults: 5000,
	}).Do()
	must(err)

	findUser := func(userID string) *userModel {
		for _, p := range users {
			if p.ID() == userID {
				return p
			}
		}
		return &userModel{}
	}

	db.Exec("DELETE FROM payments;")
	db.Exec("DELETE FROM user_assignments;")
	db.Exec("DELETE FROM course_assignments;")
	db.Exec("DELETE FROM course_contents;")
	db.Exec("DELETE FROM course_options;")
	db.Exec("DELETE FROM courses;")
	db.Exec("DELETE FROM roles;")
	db.Exec("DELETE FROM users;")

	log.Println("migrate users")
	stmt, err := db.Prepare(`
		INSERT INTO users
			(id, username, name, image, about_me, email, created_at, updated_at)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8);
	`)
	must(err)
	for _, p := range respUser.Users {
		u := findUser(p.LocalId)
		x := model.User{
			ID:        p.LocalId,
			Username:  u.Username,
			Name:      u.Name,
			Image:     u.Photo,
			AboutMe:   u.AboutMe,
			Email:     p.Email,
			CreatedAt: time.Unix(0, p.CreatedAt*1000000),
			UpdatedAt: u.UpdatedAt,
		}
		id := p.LocalId
		username := u.Username
		if len(username) == 0 {
			username = id
		}
		name := u.Name
		if len(name) == 0 {
			name = p.DisplayName
		}
		createdAt := time.Unix(0, p.CreatedAt*1000000)
		updatedAt := u.UpdatedAt
		if updatedAt.IsZero() {
			updatedAt = time.Now()
		}
		if createdAt.IsZero() {
			createdAt = updatedAt
		}
		var email *string
		if len(x.Email) > 0 {
			email = &x.Email
		}
		_, err = stmt.Exec(id, username, name, u.Photo, u.AboutMe, email, createdAt, updatedAt)
		must(err)
	}

	log.Println("migrate role")
	stmt, err = db.Prepare(`
		INSERT INTO roles
			(id, admin, instructor, created_at, updated_at)
		VALUES
			($1, $2, $3, $4, $5);
	`)
	must(err)
	for _, p := range roles {
		_, err = stmt.Exec(p.ID(), p.Admin, p.Instructor, p.CreatedAt, p.UpdatedAt)
		must(err)
	}

	log.Println("migrate courses")
	stmt, err = db.Prepare(`
		INSERT INTO courses
			(user_id, title, short_desc, long_desc, image, start, url, type, price, discount, enroll_detail, created_at, updated_at)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id;
	`)
	must(err)
	stmt2, err := db.Prepare(`
		INSERT INTO course_options
			(id, public, enroll, attend, assignment, discount)
		VALUES
			($1, $2, $3, $4, $5, $6);
	`)
	must(err)
	stmt3, err := db.Prepare(`
		INSERT INTO course_contents
			(course_id, i, title, long_desc, video_id, video_type, download_url)
		VALUES
			($1, $2, $3, $4, $5, $6, $7);
	`)
	must(err)
	for _, p := range courses {
		var tp int
		switch p.Type {
		case CourseTypeLive:
			tp = model.Live
		case CourseTypeVideo:
			tp = model.Video
		case CourseTypeEbook:
			tp = model.EBook
		}
		var st *time.Time
		if !p.Start.IsZero() {
			st = &p.Start
		}
		var url *string
		if len(p.URL) > 0 {
			url = &p.URL
		}
		var id int64
		err = stmt.QueryRow(p.Owner, p.Title, p.ShortDescription, p.Description, p.Photo, st, url, tp, p.Price, p.DiscountedPrice, p.EnrollDetail, p.CreatedAt, p.UpdatedAt).Scan(&id)
		must(err)
		p.newID = id

		_, err = stmt2.Exec(id, p.Options.Public, p.Options.Enroll, p.Options.Attend, p.Options.Assignment, p.Options.Discount)
		must(err)
		for i, c := range p.Contents {
			var vt int
			if len(p.Video) > 0 {
				vt = model.Youtube
			}
			_, err = stmt3.Exec(id, i, c.Title, c.Description, c.Video, vt, c.DownloadURL)
			must(err)
		}
	}

	log.Println("migrate payments")
	stmt, err = db.Prepare(`
		INSERT INTO payments
			(user_id, course_id, image, price, original_price, code, status, created_at, updated_at, at)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);
	`)
	must(err)
	for _, p := range payments {
		c := findCourse(p.CourseID)
		if c == nil {
			log.Println("course not found")
			continue
		}
		var status int
		switch p.Status {
		case statusWaiting:
			status = model.Pending
		case statusApproved:
			status = model.Accepted
		case statusRejected:
			status = model.Rejected
		}
		_, err = stmt.Exec(p.UserID, c.newID, p.URL, p.Price, p.OriginalPrice, p.Code, status, p.CreatedAt, p.UpdatedAt, p.At)
		must(err)
	}

	log.Println("migrate assignments")
	stmt, err = db.Prepare(`
		INSERT INTO assignments
			(course_id, title, long_desc, open, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id;
	`)
	must(err)
	for _, p := range assignments {
		c := findCourse(p.CourseID)
		var id int64
		err = stmt.QueryRow(c.newID, p.Title, p.Description, p.Open, p.CreatedAt, p.UpdatedAt).Scan(&id)
		must(err)
		p.newID = id
	}

	log.Println("migrate enroll")
	stmt, err = db.Prepare(`
		INSERT INTO enrolls
			(user_id, course_id, created_at)
		VALUES
			($1, $2, $3);
	`)
	must(err)
	for _, p := range enrolls {
		c := findCourse(p.CourseID)
		_, err = stmt.Exec(p.UserID, c.newID, p.CreatedAt)
		must(err)
	}

	log.Println("migrate attend")
	stmt, err = db.Prepare(`
		INSERT INTO attends
			(user_id, course_id, created_at)
		VALUES
			($1, $2, $3);
	`)
	must(err)
	for _, p := range attends {
		c := findCourse(p.CourseID)
		_, err = stmt.Exec(p.UserID, c.newID, p.CreatedAt)
		must(err)
	}

	log.Println("migrate user assignments")
	stmt, err = db.Prepare(`
		INSERT INTO user_assignments
			(user_id, assignment_id, download_url, created_at)
		VALUES
			($1, $2, $3, $4);
	`)
	must(err)
	for _, p := range userAssignments {
		c := findAssignment(p.AssignmentID)
		_, err = stmt.Exec(p.UserID, c.newID, p.URL, p.CreatedAt)
		must(err)
	}
}

func must(err error) {
	if err != nil {
		// log.Fatal(err)
		log.Println(err)
	}
}

type userModel struct {
	ds.StringIDModel
	ds.StampModel
	Username string
	Name     string `datastore:",noindex"`
	Photo    string `datastore:",noindex"`
	AboutMe  string `datastore:",noindex"`
}

type roleModel struct {
	ds.StringIDModel
	ds.StampModel

	// roles
	Admin      bool
	Instructor bool
}

type courseModel struct {
	newID       int64
	assignments []string

	ds.StringIDModel
	ds.StampModel
	Title            string `datastore:",noindex"`
	ShortDescription string `datastore:",noindex"`
	Description      string `datastore:",noindex"` // Markdown
	Photo            string `datastore:",noindex"` // URL
	Owner            string
	Start            time.Time
	URL              string
	Type             courseType
	Video            string `datastore:",noindex"` // Cover Video
	Price            float64
	DiscountedPrice  float64
	Options          courseOption
	Contents         courseContents `datastore:",noindex"`
	EnrollDetail     string         `datastore:",noindex"`
}

type courseModels []*courseModel

type courseOption struct {
	Public     bool
	Enroll     bool `datastore:",noindex"`
	Attend     bool `datastore:",noindex"`
	Assignment bool `datastore:",noindex"`
	Discount   bool
}

type courseContent struct {
	Title       string `datastore:",noindex"`
	Description string `datastore:",noindex"` // Markdown
	Video       string `datastore:",noindex"` // Youtube ID
	DownloadURL string `datastore:",noindex"` // Video download link
}

type courseContents []courseContent

type courseType string

// CourseType
const (
	CourseTypeLive  courseType = "live"
	CourseTypeVideo courseType = "video"
	CourseTypeEbook courseType = "ebook"
)

type paymentModel struct {
	ds.StringIDModel
	ds.StampModel
	UserID        string
	CourseID      string
	OriginalPrice float64 `datastore:",noindex"`
	Price         float64 `datastore:",noindex"`
	Code          string
	URL           string `datastore:",noindex"`
	Status        status
	At            time.Time
}

type status string

const (
	statusWaiting  status = "waiting"
	statusApproved status = "approved"
	statusRejected status = "rejected"
)

type attendModel struct {
	ds.StringIDModel
	ds.StampModel
	UserID   string
	CourseID string
}

type enrollModel struct {
	ds.StringIDModel
	ds.StampModel
	UserID   string
	CourseID string
}

type assignment struct {
	newID int64

	ds.StringIDModel
	ds.StampModel
	CourseID    string
	Title       string `datastore:",noindex"`
	Description string `datastore:",noindex"`
	Open        bool   `datastore:",noindex"`
}

type userAssignment struct {
	ds.StringIDModel
	ds.StampModel
	AssignmentID string
	UserID       string
	URL          string `datastore:",noindex"`
}
