package app

import (
	"context"
	"html/template"
	"net/http"
	"net/url"
	"time"

	"github.com/acoshift/header"
	"github.com/acoshift/hime"
	"github.com/acoshift/session"
	"github.com/dustin/go-humanize"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"

	"github.com/acoshift/acourse/appctx"
	"github.com/acoshift/acourse/entity"
)

func loadTemplates(app hime.App) {
	app.
		BeforeRender(func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				appctx.GetSession(r.Context()).Flash().Clear()
				w.Header().Set(header.CacheControl, "no-cache, no-store, must-revalidate")
				h.ServeHTTP(w, r)
			})
		}).
		TemplateFuncs(template.FuncMap{
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
		}).
		Component("layout.tmpl").
		Template("index", "index.tmpl", "app.tmpl", "component/course-card.tmpl").
		Template("error.not-found", "not-found.tmpl", "app.tmpl").
		Template("signin", "signin.tmpl", "auth.tmpl").
		Template("signin.password", "signin-password.tmpl", "auth.tmpl").
		Template("signup", "signup.tmpl", "auth.tmpl").
		Template("reset.password", "reset-password.tmpl", "auth.tmpl").
		Template("check-email", "check-email.tmpl", "auth.tmpl").
		Template("profile", "profile.tmpl", "app.tmpl",
			"component/user-profile.tmpl",
			"component/own-course-card.tmpl",
			"component/enrolled-course-card.tmpl",
		).
		Template("profile.edit", "profile-edit.tmpl", "app.tmpl").
		Template("course", "course.tmpl", "app.tmpl").
		Template("course.content", "course-content.tmpl", "app.tmpl").
		Template("course.enroll", "enroll.tmpl", "app.tmpl").
		Template("assignment", "assignment.tmpl", "app.tmpl").
		Template("editor.create", "editor/create.tmpl", "app.tmpl").
		Template("editor.course", "editor/course.tmpl", "app.tmpl").
		Template("editor.content", "editor/content.tmpl", "app.tmpl").
		Template("editor.content.create", "editor/content-create.tmpl", "app.tmpl").
		Template("editor.content.edit", "editor/content-edit.tmpl", "app.tmpl").
		Template("admin.users", "admin/users.tmpl", "app.tmpl").
		Template("admin.courses", "admin/courses.tmpl", "app.tmpl").
		Template("admin.payments", "admin/payments.tmpl", "app.tmpl").
		Template("admin.payments.reject", "admin/payment-reject.tmpl", "app.tmpl")
}

func newPage(ctx context.Context) map[string]interface{} {
	return map[string]interface{}{
		"Title": "Acourse",
		"Desc":  "แหล่งเรียนรู้ออนไลน์ ที่ทุกคนเข้าถึงได้",
		"Image": "https://storage.googleapis.com/acourse/static/62b9eb0e-3668-4f9f-86b7-a11349938f7a.jpg",
		"URL":   "https://acourse.io",
		"Me":    appctx.GetUser(ctx),
		"Flash": session.Get(ctx, "sess").Flash().Values(),
		"XSRF":  template.HTML(`<input type="hidden" name="X" value="` + appctx.GetXSRFToken(ctx) + `">`),
	}
}

// SignIn renders signin view
func SignIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data := struct {
		*Page
	}{newPage(ctx)}
	render(ctx, w, tmplSignIn, &data)
}

// SignInPassword renders signin-password view
func SignInPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data := struct {
		*Page
	}{newPage(ctx)}
	render(ctx, w, tmplSignInPassword, &data)
}

// SignUp renders signup view
func SignUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data := struct {
		*Page
	}{newPage(ctx)}
	render(ctx, w, tmplSignUp, &data)
}

// ResetPassword render reset password view
func ResetPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	page := newPage(ctx)
	data := struct {
		*Page
	}{page}
	render(ctx, w, tmplResetPassword, &data)
}

// CheckEmail render check email view
func CheckEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	page := newPage(ctx)
	data := struct {
		*Page
	}{page}
	render(ctx, w, tmplCheckEmail, &data)
}

// Profile renders profile view
func Profile(w http.ResponseWriter, r *http.Request, ownCourses, enrolledCourses []*entity.Course) {
	ctx := r.Context()
	page := newPage(ctx)
	me := appctx.GetUser(ctx)
	page.Title = me.Username + " | " + page.Title

	data := struct {
		*Page
		OwnCourses      []*entity.Course
		EnrolledCourses []*entity.Course
	}{page, ownCourses, enrolledCourses}
	render(ctx, w, tmplProfile, &data)
}

// ProfileEdit renders profile edit view
func ProfileEdit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	me := appctx.GetUser(ctx)
	page := newPage(ctx)
	page.Title = me.Username + " | " + page.Title

	data := struct {
		*Page
	}{page}
	render(ctx, w, tmplProfileEdit, &data)
}

// Course renders course view
func Course(w http.ResponseWriter, r *http.Request, course *entity.Course, enrolled bool, owned bool, pendingEnroll bool) {
	ctx := r.Context()
	page := newPage(ctx)
	page.Title = course.Title + " | " + page.Title
	page.Desc = course.ShortDesc
	page.Image = course.Image
	page.URL = baseURL + "/course/" + url.PathEscape(course.Link())

	data := struct {
		*Page
		Course        *entity.Course
		Enrolled      bool
		Owned         bool
		PendingEnroll bool
	}{page, course, enrolled, owned, pendingEnroll}
	render(ctx, w, tmplCourse, &data)
}

// CourseContent renders course content view
func CourseContent(w http.ResponseWriter, r *http.Request, course *entity.Course, content *entity.CourseContent) {
	ctx := r.Context()
	page := newPage(ctx)
	page.Title = course.Title + " | " + page.Title
	page.Desc = course.ShortDesc
	page.Image = course.Image

	data := struct {
		*Page
		Course  *entity.Course
		Content *entity.CourseContent
	}{page, course, content}
	render(ctx, w, tmplCourseContent, &data)
}

// EditorCreate renders course create view
func EditorCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	page := newPage(ctx)

	data := struct {
		*Page
	}{page}
	render(ctx, w, tmplEditorCreate, &data)
}

// EditorCourse renders course edit view
func EditorCourse(w http.ResponseWriter, r *http.Request, course *entity.Course) {
	ctx := r.Context()
	page := newPage(ctx)

	data := struct {
		*Page
		Course *entity.Course
	}{page, course}
	render(ctx, w, tmplEditorCourse, &data)
}

// EditorContent renders editor content view
func EditorContent(w http.ResponseWriter, r *http.Request, course *entity.Course) {
	ctx := r.Context()
	page := newPage(ctx)

	data := struct {
		*Page
		Course *entity.Course
	}{page, course}
	render(ctx, w, tmplEditorContent, &data)
}

// EditorContentCreate renders editor content create view
func EditorContentCreate(w http.ResponseWriter, r *http.Request, course *entity.Course) {
	ctx := r.Context()
	page := newPage(ctx)

	data := struct {
		*Page
		Course *entity.Course
	}{page, course}
	render(ctx, w, tmplEditorContentCreate, &data)
}

// EditorContentEdit renders editor content edit view
func EditorContentEdit(w http.ResponseWriter, r *http.Request, course *entity.Course, content *entity.CourseContent) {
	ctx := r.Context()
	page := newPage(ctx)

	data := struct {
		*Page
		Course  *entity.Course
		Content *entity.CourseContent
	}{page, course, content}
	render(ctx, w, tmplEditorContentEdit, &data)
}

// CourseEnroll renders course enroll view
func CourseEnroll(w http.ResponseWriter, r *http.Request, course *entity.Course) {
	ctx := r.Context()
	page := newPage(ctx)
	page.Title = course.Title + " | " + page.Title
	page.Desc = course.ShortDesc
	page.Image = course.Image
	page.URL = BaseURL + "/course/" + url.PathEscape(course.Link())

	data := struct {
		*Page
		Course *entity.Course
	}{page, course}
	render(ctx, w, tmplCourseEnroll, &data)
}

// Assignment render assignment view
func Assignment(w http.ResponseWriter, r *http.Request, course *entity.Course, assignments []*entity.Assignment) {
	ctx := r.Context()
	page := newPage(ctx)
	page.Title = course.Title + " | " + page.Title
	page.Desc = course.ShortDesc
	page.Image = course.Image
	page.URL = BaseURL + "/course/" + url.PathEscape(course.Link())

	data := struct {
		*Page
		Course      *entity.Course
		Assignments []*entity.Assignment
	}{page, course, assignments}
	render(ctx, w, tmplAssignment, &data)
}

// AdminUsers renders admin users view
func AdminUsers(w http.ResponseWriter, r *http.Request, users []*entity.User, currentPage, totalPage int) {
	ctx := r.Context()
	page := newPage(ctx)

	data := struct {
		*Page
		Users       []*entity.User
		CurrentPage int
		TotalPage   int
	}{page, users, currentPage, totalPage}
	render(ctx, w, tmplAdminUsers, &data)
}

// AdminCourses renders admin courses view
func AdminCourses(w http.ResponseWriter, r *http.Request, courses []*entity.Course, currentPage, totalPage int) {
	ctx := r.Context()
	page := newPage(ctx)

	data := struct {
		*Page
		Courses     []*entity.Course
		CurrentPage int
		TotalPage   int
	}{page, courses, currentPage, totalPage}
	render(ctx, w, tmplAdminCourses, &data)
}

// AdminPayments renders admin payments view
func AdminPayments(w http.ResponseWriter, r *http.Request, payments []*entity.Payment, currentPage, totalPage int) {
	ctx := r.Context()
	page := newPage(ctx)

	data := struct {
		*Page
		Payments    []*entity.Payment
		CurrentPage int
		TotalPage   int
	}{page, payments, currentPage, totalPage}
	render(ctx, w, tmplAdminPayments, &data)
}

// AdminPaymentReject renders admin payment reject view
func AdminPaymentReject(w http.ResponseWriter, r *http.Request, payment *entity.Payment, message string) {
	ctx := r.Context()
	page := newPage(ctx)

	data := struct {
		*Page
		Payment *entity.Payment
		Message string
	}{page, payment, message}
	render(ctx, w, tmplAdminPaymentReject, &data)
}
