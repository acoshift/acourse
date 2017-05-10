package main

import (
	"context"
	"log"
	"time"

	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/ds"
	"github.com/garyburd/redigo/redis"
)

var ctx = context.Background()

var client, _ = ds.NewClient(ctx, "acourse-156413")
var conn, _ = redis.Dial("tcp", "localhost:6379", redis.DialDatabase(4))

// migrate all data from datastore and firebase to redis
func main() {
	time.Local = time.UTC

	var users []*userModel
	var roles []*roleModel
	var courses []*courseModel
	var payments []*paymentModel
	// var attends []*attendModel
	// var enrolls []*enrollModel
	// var assignments []*assignment
	// var userAssignments []*userAssignment

	must(client.Query(ctx, "User", &users))
	must(client.Query(ctx, "Role", &roles))
	must(client.Query(ctx, "Course", &courses))
	must(client.Query(ctx, "Payment", &payments))
	// must(client.Query(ctx, "Attend", &attends))
	// must(client.Query(ctx, "Enroll", &enrolls))
	// must(client.Query(ctx, "Assignment", &assignments))
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

	// save users and create mapper
	for _, p := range users {
		r := findRole(p.ID())
		x := model.User{
			Username:  p.Username,
			Image:     p.Photo,
			AboutMe:   p.AboutMe,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		}
		x.Role().Admin = r.Admin
		x.Role().Instructor = r.Instructor

		x.SetID(p.ID())
		must(x.Save(conn))
	}

	// save course and create mapper
	for _, p := range courses {
		x := model.Course{
			CreatedAt:    p.CreatedAt,
			UpdatedAt:    p.UpdatedAt,
			Title:        p.Title,
			ShortDesc:    p.ShortDescription,
			Desc:         p.Description,
			Image:        p.Photo,
			UserID:       p.Owner, // TODO: must map from user
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
			x.Contents = append(x.Contents, &model.CourseContent{
				Title:       c.Title,
				Desc:        c.Description,
				YoutubeID:   c.Video,
				DownloadURL: c.DownloadURL,
			})
		}
		must(x.Save(conn))
		p.newID = x.ID()
	}

	// save payments
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
	newID string

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
