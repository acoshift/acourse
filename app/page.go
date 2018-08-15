package app

import (
	"context"
	"html/template"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"

	"github.com/acoshift/acourse/appctx"
	"github.com/acoshift/acourse/entity"
)

// TemplateFunc returns app's template funcs
func TemplateFunc() template.FuncMap {
	return template.FuncMap{
		"currency": func(v float64) string {
			return humanize.FormatFloat("#,###.##", v)
		},
		"paginate": func(p, n int) []int {
			r := make([]int, 0, 7)
			r = append(r, 1)
			if n <= 1 {
				return r
			}
			if n <= 2 {
				return append(r, 2)
			}
			if p <= 3 {
				r = append(r, 2, 3)
			}
			if p > 3 {
				r = append(r, -1, p-1)
				if p < n {
					r = append(r, p)
				}
			}
			if n-p+1 >= 3 && p >= 3 {
				r = append(r, p+1)
			}
			if n-p >= 3 {
				r = append(r, -1)
			}
			if n >= 4 {
				r = append(r, n)
			}
			return r
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
		"decr": func(v int) int {
			return v - 1
		},
	}
}

func newPage(ctx context.Context) map[string]interface{} {
	return map[string]interface{}{
		"Title": "Acourse",
		"Desc":  "แหล่งเรียนรู้ออนไลน์ ที่ทุกคนเข้าถึงได้",
		"Image": "https://storage.googleapis.com/acourse/static/62b9eb0e-3668-4f9f-86b7-a11349938f7a.jpg",
		"URL":   "https://acourse.io",
		"Me":    appctx.GetUser(ctx),
		"Flash": appctx.GetSession(ctx).Flash(),
	}
}
