package course

import (
	"time"

	"github.com/acoshift/acourse/pkg/acourse"
	"github.com/acoshift/ds"
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

type courseModel struct {
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

const (
	kindAttend = "Attend"
	kindEnroll = "Enroll"
	kindCourse = "Course"
)

func formatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

func toCourse(x *courseModel) *acourse.Course {
	r := &acourse.Course{
		Id:               x.ID(),
		CreatedAt:        formatTime(x.CreatedAt),
		UpdatedAt:        formatTime(x.UpdatedAt),
		Title:            x.Title,
		ShortDescription: x.ShortDescription,
		Description:      x.Description,
		Photo:            x.Photo,
		Owner:            x.Owner,
		Start:            formatTime(x.Start),
		Url:              x.URL,
		Type:             string(x.Type),
		Video:            x.Video,
		Price:            x.Price,
		DiscountedPrice:  x.DiscountedPrice,
		Options: &acourse.Course_Option{
			Public:     x.Options.Public,
			Enroll:     x.Options.Enroll,
			Attend:     x.Options.Attend,
			Assignment: x.Options.Assignment,
			Discount:   x.Options.Discount,
		},
		EnrollDetail: x.EnrollDetail,
	}

	r.Contents = make([]*acourse.Course_Content, len(x.Contents))
	for i, c := range x.Contents {
		r.Contents[i] = &acourse.Course_Content{
			Title:       c.Title,
			Description: c.Description,
			Video:       c.Video,
			DownloadURL: c.DownloadURL,
		}
	}

	return r
}

func toCourses(xs courseModels) []*acourse.Course {
	rs := make([]*acourse.Course, len(xs))
	for i, x := range xs {
		rs[i] = toCourse(x)
	}
	return rs
}
