package view

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/flash"
	"github.com/acoshift/header"
	"github.com/dustin/go-humanize"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
	"gopkg.in/yaml.v2"
)

const templateDir = "template"

var (
	m          = minify.New()
	muExecute  = &sync.Mutex{}
	templates  = make(map[interface{}]*templateStruct)
	loc        *time.Location
	staticConf = make(map[string]string)
)

type templateStruct struct {
	*template.Template
	set []string
}

func init() {
	var err error
	loc, err = time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Fatal(err)
	}

	m.AddFunc("text/html", html.Minify)
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/javascript", js.Minify)

	// load static config
	{
		bs, _ := ioutil.ReadFile("static.yaml")
		yaml.Unmarshal(bs, &staticConf)
	}

	parseTemplate(keyIndex{}, []string{"index.tmpl", "app.tmpl", "layout.tmpl", "component/course-card.tmpl"})
	parseTemplate(keySignIn{}, []string{"signin.tmpl", "auth.tmpl", "layout.tmpl"})
	parseTemplate(keySignUp{}, []string{"signup.tmpl", "auth.tmpl", "layout.tmpl"})
	parseTemplate(keyProfile{}, []string{
		"profile.tmpl", "app.tmpl", "layout.tmpl",
		"component/user-profile.tmpl",
		"component/own-course-card.tmpl",
		"component/enrolled-course-card.tmpl",
	})
	parseTemplate(keyProfileEdit{}, []string{"profile-edit.tmpl", "app.tmpl", "layout.tmpl"})
	parseTemplate(keyCourse{}, []string{"course.tmpl", "app.tmpl", "layout.tmpl"})
	parseTemplate(keyCourseContent{}, []string{"course-content.tmpl", "app.tmpl", "layout.tmpl"})
	parseTemplate(keyAssignment{}, []string{"assignment.tmpl", "app.tmpl", "layout.tmpl"})
	parseTemplate(keyCourseEnroll{}, []string{"enroll.tmpl", "app.tmpl", "layout.tmpl"})
	parseTemplate(keyEditorCreate{}, []string{"editor/create.tmpl", "app.tmpl", "layout.tmpl"})
	parseTemplate(keyEditorCourse{}, []string{"editor/course.tmpl", "app.tmpl", "layout.tmpl"})
	parseTemplate(keyEditorContent{}, []string{"editor/content.tmpl", "app.tmpl", "layout.tmpl"})
	parseTemplate(keyEditorContentCreate{}, []string{"editor/content-create.tmpl", "app.tmpl", "layout.tmpl"})
	parseTemplate(keyEditorContentEdit{}, []string{"editor/content-edit.tmpl", "app.tmpl", "layout.tmpl"})
	parseTemplate(keyAdminUsers{}, []string{"admin/users.tmpl", "app.tmpl", "layout.tmpl"})
	parseTemplate(keyAdminCourses{}, []string{"admin/courses.tmpl", "app.tmpl", "layout.tmpl"})
	parseTemplate(keyAdminPayments{}, []string{"admin/payments.tmpl", "app.tmpl", "layout.tmpl"})
	parseTemplate(keyAdminPaymentReject{}, []string{"admin/payment-reject.tmpl", "app.tmpl", "layout.tmpl"})
	parseTemplate(keyNotFound{}, []string{"not-found.tmpl", "app.tmpl", "layout.tmpl"})
}

func joinTemplateDir(files []string) []string {
	r := make([]string, len(files))
	for i, f := range files {
		r[i] = filepath.Join(templateDir, f)
	}
	return r
}

func parseTemplate(key interface{}, set []string) {
	templateName := strings.TrimSuffix(set[0], ".tmpl")
	t := template.New("")
	t.Funcs(template.FuncMap{
		"templateName": func() string {
			return templateName
		},
		"currency": func(v float64) string {
			return humanize.FormatFloat("#,###.##", v)
		},
		"static": func(s string) string {
			return "/~/" + staticConf[s]
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
			case model.Live:
				return "Live"
			case model.Video:
				return "Video"
			case model.EBook:
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
			return model.Live
		},
		"video": func() int {
			return model.Video
		},
		"eBook": func() int {
			return model.EBook
		},
		"pending": func() int {
			return model.Pending
		},
		"accepted": func() int {
			return model.Accepted
		},
		"rejected": func() int {
			return model.Rejected
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
	})
	_, err := t.ParseFiles(joinTemplateDir(set)...)
	if err != nil {
		log.Fatalf("view: parse template %s error; %v", templateName, err)
	}
	t = t.Lookup("root")
	if t == nil {
		log.Fatalf("view: root template not found in %s", templateName)
	}
	templates[key] = &templateStruct{
		Template: t,
		set:      set,
	}
}

func render(ctx context.Context, w http.ResponseWriter, key, data interface{}) {
	t := templates[key]
	if t == nil {
		http.Error(w, fmt.Sprintf("template not found"), http.StatusInternalServerError)
		return
	}
	if dev {
		muExecute.Lock()
		defer muExecute.Unlock()
		parseTemplate(key, t.set)
		t = templates[key]
	}

	w.Header().Set(header.ContentType, "text/html; charset=utf-8")
	w.Header().Set(header.CacheControl, "no-cache, no-store, must-revalidate, max-age=0")
	pipe := &bytes.Buffer{}
	err := t.Execute(pipe, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = m.Minify("text/html", w, pipe)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	f := flash.Get(ctx)
	f.Clear()
}
