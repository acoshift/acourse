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

	var courses []*Course

	must(client.Query(ctx, "Course", &courses))

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
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Course model
type Course struct {
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
	Type             CourseType
	Video            string `datastore:",noindex"` // Cover Video
	Price            float64
	DiscountedPrice  float64
	Options          CourseOption
	Contents         CourseContents `datastore:",noindex"`
	EnrollDetail     string         `datastore:",noindex"`
}

// Courses model
type Courses []*Course

// CourseOption type
type CourseOption struct {
	Public     bool
	Enroll     bool `datastore:",noindex"`
	Attend     bool `datastore:",noindex"`
	Assignment bool `datastore:",noindex"`
	Discount   bool
}

// CourseContent type
type CourseContent struct {
	Title       string `datastore:",noindex"`
	Description string `datastore:",noindex"` // Markdown
	Video       string `datastore:",noindex"` // Youtube ID
	DownloadURL string `datastore:",noindex"` // Video download link
}

// CourseContents type
type CourseContents []CourseContent

// CourseType type
type CourseType string

// CourseType
const (
	CourseTypeLive  CourseType = "live"
	CourseTypeVideo CourseType = "video"
	CourseTypeEbook CourseType = "ebook"
)
