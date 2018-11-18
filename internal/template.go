package internal

import (
	"html/template"
	"time"

	"github.com/dustin/go-humanize"

	"github.com/acoshift/acourse/internal/pkg/markdown"
	"github.com/acoshift/acourse/internal/pkg/model/course"
	"github.com/acoshift/acourse/internal/pkg/model/payment"
)

// TemplateFunc returns template funcs
func TemplateFunc(loc *time.Location) template.FuncMap {
	return template.FuncMap{
		"currency": func(v float64) string {
			return humanize.FormatFloat("#,###.##", v)
		},
		"courseType": func(v int) string {
			switch v {
			case course.Live:
				return "Live"
			case course.Video:
				return "Video"
			case course.EBook:
				return "eBook"
			default:
				return ""
			}
		},
		"date": func(v time.Time) string {
			return v.In(loc).Format("02/01/2006")
		},
		"dateTime": func(v time.Time) string {
			return v.In(loc).Format("02/01/2006 15:04:05")
		},
		"dateInput": func(v time.Time) string {
			return v.Format("2006-01-02")
		},
		"markdown": markdown.HTML,
		"live": func() int {
			return course.Live
		},
		"video": func() int {
			return course.Video
		},
		"eBook": func() int {
			return course.EBook
		},
		"pending": func() int {
			return payment.Pending
		},
		"accepted": func() int {
			return payment.Accepted
		},
		"rejected": func() int {
			return payment.Rejected
		},
		"refunded": func() int {
			return payment.Refunded
		},
		"html": func(v string) template.HTML {
			return template.HTML(v)
		},
		"incr": func(v int) int {
			return v + 1
		},
		"fallbackImage": func() string {
			return "/-/placeholder-img.svg"
		},
	}
}
