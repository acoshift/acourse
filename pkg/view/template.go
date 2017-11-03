package view

import (
	"bytes"
	"context"
	"html/template"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/acoshift/header"
	"github.com/dustin/go-humanize"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
	"gopkg.in/yaml.v2"

	"github.com/acoshift/acourse/pkg/app"
)

// templates
var (
	tmplIndex          = parse("index.tmpl", "app.tmpl", "layout.tmpl", "component/course-card.tmpl")
	tmplNotFound       = parse("not-found.tmpl", "app.tmpl", "layout.tmpl")
	tmplSignIn         = parse("signin.tmpl", "auth.tmpl", "layout.tmpl")
	tmplSignInPassword = parse("signin-password.tmpl", "auth.tmpl", "layout.tmpl")
	tmplSignUp         = parse("signup.tmpl", "auth.tmpl", "layout.tmpl")
	tmplResetPassword  = parse("reset-password.tmpl", "auth.tmpl", "layout.tmpl")
	tmplCheckEmail     = parse("check-email.tmpl", "auth.tmpl", "layout.tmpl")
	tmplProfile        = parse(
		"profile.tmpl", "app.tmpl", "layout.tmpl",
		"component/user-profile.tmpl",
		"component/own-course-card.tmpl",
		"component/enrolled-course-card.tmpl",
	)
	tmplProfileEdit = parse("profile-edit.tmpl", "app.tmpl", "layout.tmpl")
	// tmplUser                = parse()
	tmplCourse              = parse("course.tmpl", "app.tmpl", "layout.tmpl")
	tmplCourseContent       = parse("course-content.tmpl", "app.tmpl", "layout.tmpl")
	tmplCourseEnroll        = parse("enroll.tmpl", "app.tmpl", "layout.tmpl")
	tmplAssignment          = parse("assignment.tmpl", "app.tmpl", "layout.tmpl")
	tmplEditorCreate        = parse("editor/create.tmpl", "app.tmpl", "layout.tmpl")
	tmplEditorCourse        = parse("editor/course.tmpl", "app.tmpl", "layout.tmpl")
	tmplEditorContent       = parse("editor/content.tmpl", "app.tmpl", "layout.tmpl")
	tmplEditorContentCreate = parse("editor/content-create.tmpl", "app.tmpl", "layout.tmpl")
	tmplEditorContentEdit   = parse("editor/content-edit.tmpl", "app.tmpl", "layout.tmpl")
	tmplAdminUsers          = parse("admin/users.tmpl", "app.tmpl", "layout.tmpl")
	tmplAdminCourses        = parse("admin/courses.tmpl", "app.tmpl", "layout.tmpl")
	tmplAdminPayments       = parse("admin/payments.tmpl", "app.tmpl", "layout.tmpl")
	tmplAdminPaymentReject  = parse("admin/payment-reject.tmpl", "app.tmpl", "layout.tmpl")
)

const templateDir = "template"

var (
	m          = minify.New()
	loc        *time.Location
	staticConf = make(map[string]string)
)

type tmpl struct {
	*template.Template
	set []string
}

func init() {
	var err error
	loc, err = time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Fatal(err)
	}

	// add mime types
	mime.AddExtensionType(".js", "text/javascript")

	// add minifier functions
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/javascript", js.Minify)

	// load static config
	{
		bs, _ := ioutil.ReadFile("static.yaml")
		yaml.Unmarshal(bs, &staticConf)
	}
}

func joinTemplateDir(files []string) []string {
	r := make([]string, len(files))
	for i, f := range files {
		r[i] = filepath.Join(templateDir, f)
	}
	return r
}

func parse(set ...string) *tmpl {
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
			case app.Live:
				return "Live"
			case app.Video:
				return "Video"
			case app.EBook:
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
			return app.Live
		},
		"video": func() int {
			return app.Video
		},
		"eBook": func() int {
			return app.EBook
		},
		"pending": func() int {
			return app.Pending
		},
		"accepted": func() int {
			return app.Accepted
		},
		"rejected": func() int {
			return app.Rejected
		},
		"refunded": func() int {
			return app.Refunded
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
	return &tmpl{
		Template: t,
		set:      set,
	}
}

func renderWithStatusCode(ctx context.Context, w http.ResponseWriter, code int, t *tmpl, data interface{}) {
	if dev {
		// reload template for dev env
		t = parse(t.set...)
	}

	// clear flash after render
	app.GetSession(ctx).Flash().Clear()

	// set header for html
	w.Header().Set(header.ContentType, "text/html; charset=utf-8")
	w.Header().Set(header.CacheControl, "no-cache, no-store, must-revalidate, max-age=0")
	w.WriteHeader(code)

	// use buffer is faster than pipe stream for this case
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
}

func render(ctx context.Context, w http.ResponseWriter, t *tmpl, data interface{}) {
	renderWithStatusCode(ctx, w, http.StatusOK, t, data)
}
