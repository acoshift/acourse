package internal

import (
	"html/template"
	"time"

	"github.com/dustin/go-humanize"

	"github.com/acoshift/acourse/entity"
	"github.com/acoshift/acourse/view"
)

// TemplateFunc returns template funcs
func TemplateFunc(loc *time.Location) template.FuncMap {
	return template.FuncMap{
		"currency": func(v float64) string {
			return humanize.FormatFloat("#,###.##", v)
		},
		"courseType": func(v int) string {
			switch v {
			case entity.Live:
				return "Live"
			case entity.Video:
				return "Video"
			case entity.EBook:
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
		"markdown": view.MarkdownHTML,
		"live": func() int {
			return entity.Live
		},
		"video": func() int {
			return entity.Video
		},
		"eBook": func() int {
			return entity.EBook
		},
		"pending": func() int {
			return entity.Pending
		},
		"accepted": func() int {
			return entity.Accepted
		},
		"rejected": func() int {
			return entity.Rejected
		},
		"refunded": func() int {
			return entity.Refunded
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
