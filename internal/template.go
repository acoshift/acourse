package internal

import (
	"html/template"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"

	"github.com/acoshift/acourse/entity"
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
		"markdown": func(s string) template.HTML {
			renderer := blackfriday.HtmlRenderer(
				0|
					blackfriday.HTML_USE_XHTML|
					blackfriday.HTML_USE_SMARTYPANTS|
					blackfriday.HTML_SMARTYPANTS_FRACTIONS|
					blackfriday.HTML_SMARTYPANTS_DASHES|
					blackfriday.HTML_SMARTYPANTS_LATEX_DASHES|
					blackfriday.HTML_HREF_TARGET_BLANK,
				"", "")
			md := blackfriday.MarkdownOptions([]byte(s), renderer, blackfriday.Options{
				Extensions: 0 |
					blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
					blackfriday.EXTENSION_TABLES |
					blackfriday.EXTENSION_FENCED_CODE |
					blackfriday.EXTENSION_AUTOLINK |
					blackfriday.EXTENSION_STRIKETHROUGH |
					blackfriday.EXTENSION_SPACE_HEADERS |
					blackfriday.EXTENSION_HEADER_IDS |
					blackfriday.EXTENSION_BACKSLASH_LINE_BREAK |
					blackfriday.EXTENSION_DEFINITION_LISTS,
			})
			p := bluemonday.UGCPolicy()
			p.AllowAttrs("target").OnElements("a")
			r := p.SanitizeBytes(md)
			return template.HTML(r)
		},
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
			return "https://storage.googleapis.com/acourse/static/d509b7d8-88ad-478c-aa40-2984878c87cd.svg"
		},
	}
}
