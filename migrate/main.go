package main

import (
	"context"
	"log"
	"time"

	"golang.org/x/oauth2/google"
	identitytoolkit "google.golang.org/api/identitytoolkit/v3"

	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/configfile"
	"github.com/acoshift/ds"
	"github.com/garyburd/redigo/redis"
)

var ctx = context.Background()

var client, _ = ds.NewClient(ctx, "acourse-156413")
var conn, _ = redis.Dial("tcp", "localhost:6379", redis.DialDatabase(4))

var config = configfile.NewReader("../config")
var serviceAccount = config.Bytes("service_account")

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

	var users []*userModel
	var roles []*roleModel
	var courses []*courseModel
	var payments []*paymentModel
	// var attends []*attendModel
	var enrolls []*enrollModel
	var assignments []*assignment
	// var userAssignments []*userAssignment

	log.Println("load old database")
	must(client.Query(ctx, "User", &users))
	must(client.Query(ctx, "Role", &roles))
	must(client.Query(ctx, "Course", &courses, ds.Order("CreatedAt")))
	must(client.Query(ctx, "Payment", &payments, ds.Order("CreatedAt")))
	// must(client.Query(ctx, "Attend", &attends))
	must(client.Query(ctx, "Enroll", &enrolls))
	must(client.Query(ctx, "Assignment", &assignments, ds.Order("CreatedAt")))
	// must(client.Query(ctx, "UserAssignment", &userAssignments))

	findRole := func(userID string) *roleModel {
		for _, p := range roles {
			if p.ID() == userID {
				return p
			}
		}
		return &roleModel{}
	}

	findCourse := func(courseID string) *courseModel {
		for _, p := range courses {
			if p.ID() == courseID {
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

	// save users and create mapper
	log.Println("migrate users")
	for _, p := range respUser.Users {
		u := findUser(p.LocalId)
		r := findRole(p.LocalId)
		x := model.User{
			Username:  u.Username,
			Name:      u.Name,
			Image:     u.Photo,
			AboutMe:   u.AboutMe,
			Email:     p.Email,
			CreatedAt: time.Unix(0, p.CreatedAt),
			UpdatedAt: u.UpdatedAt,
		}
		if len(x.Name) == 0 {
			x.Name = p.DisplayName
		}
		x.Role().Admin = r.Admin
		x.Role().Instructor = r.Instructor
		x.SetID(p.LocalId)
		must(x.Save(conn))
	}

	log.Println("migrate assignments")
	for _, p := range assignments {
		x := model.Assignment{
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
			Title:     p.Title,
			Desc:      p.Description,
			Open:      p.Open,
		}
		must(x.Save(conn))
		c := findCourse(p.CourseID)
		c.assignments = append(c.assignments, x.ID())
	}

	// save course and create mapper
	log.Println("migrate courses")
	for _, p := range courses {
		x := model.Course{
			CreatedAt:    p.CreatedAt,
			UpdatedAt:    p.UpdatedAt,
			Title:        p.Title,
			ShortDesc:    p.ShortDescription,
			Desc:         p.Description,
			Image:        p.Photo,
			UserID:       p.Owner,
			Start:        p.Start,
			URL:          p.URL,
			Price:        p.Price,
			Discount:     p.DiscountedPrice,
			EnrollDetail: p.EnrollDetail,
		}
		switch p.Type {
		case CourseTypeLive:
			x.Type = model.Live
		case CourseTypeVideo:
			x.Type = model.Video
		case CourseTypeEbook:
			x.Type = model.EBook
		}
		x.Option().Public = p.Options.Public
		x.Option().Enroll = p.Options.Enroll
		x.Option().Attend = p.Options.Attend
		x.Option().Assignment = p.Options.Assignment
		x.Option().Discount = p.Options.Discount
		for _, c := range p.Contents {
			t := model.CourseContent{
				Title:       c.Title,
				Desc:        c.Description,
				VideoID:     c.Video,
				DownloadURL: c.DownloadURL,
			}
			if len(p.Video) > 0 {
				t.VideoType = model.Youtube
			}
			x.Contents = append(x.Contents, &t)
		}
		must(x.Save(conn))
		p.newID = x.ID()
	}

	// save payments
	log.Println("migrate payments")
	for _, p := range payments {
		c := findCourse(p.CourseID)
		if c == nil {
			log.Println("course not found")
			continue
		}
		x := model.Payment{
			CourseID:      c.newID,
			UserID:        p.UserID,
			CreatedAt:     p.CreatedAt,
			UpdatedAt:     p.UpdatedAt,
			Image:         p.URL,
			Price:         p.Price,
			OriginalPrice: p.OriginalPrice,
			At:            p.At,
			Code:          p.Code,
		}
		switch p.Status {
		case statusWaiting:
			x.Status = model.Pending
		case statusApproved:
			x.Status = model.Accepted
		case statusRejected:
			x.Status = model.Rejected
		}
		must(x.Save(conn))
	}

	// save enrolls
	log.Println("migrate enrolls")
	for _, p := range enrolls {
		c := findCourse(p.CourseID)
		must(model.Enroll(conn, p.UserID, c.newID))
	}
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
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
	newID       string
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
	newID string

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
	newID string

	ds.StringIDModel
	ds.StampModel
	CourseID    string
	Title       string `datastore:",noindex"`
	Description string `datastore:",noindex"`
	Open        bool   `datastore:",noindex"`
}

type userAssignment struct {
	newID string

	ds.StringIDModel
	ds.StampModel
	AssignmentID string
	UserID       string
	URL          string `datastore:",noindex"`
}
