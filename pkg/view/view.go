package view

import (
	"context"
	"html/template"
	"math/rand"
	"net/http"
	"net/url"

	"github.com/acoshift/flash"
	"github.com/acoshift/header"
	"github.com/acoshift/session"

	"github.com/acoshift/acourse/pkg/appctx"
	"github.com/acoshift/acourse/pkg/entity"
)

// view's config
var (
	BaseURL string
)

// Page type provides layout data like title, description, and og
type Page struct {
	Title string
	Desc  string
	Image string
	URL   string
	Me    *entity.User
	Flash flash.Data
	XSRF  template.HTML
}

var defaultPage = Page{
	Title: "Acourse",
	Desc:  "แหล่งเรียนรู้ออนไลน์ ที่ทุกคนเข้าถึงได้",
	Image: "https://storage.googleapis.com/acourse/static/62b9eb0e-3668-4f9f-86b7-a11349938f7a.jpg",
	URL:   "https://acourse.io",
}

func newPage(ctx context.Context) *Page {
	p := defaultPage
	p.Me = appctx.GetUser(ctx)
	p.Flash = session.Get(ctx, "sess").Flash().Values()
	p.XSRF = template.HTML(`<input type="hidden" name="X" value="` + appctx.GetXSRFToken(ctx) + `">`)
	return &p
}

// Index renders index view
func Index(w http.ResponseWriter, r *http.Request, courses []*entity.Course) {
	ctx := r.Context()
	data := struct {
		*Page
		Courses []*entity.Course
	}{newPage(ctx), courses}
	render(ctx, w, tmplIndex, &data)
}

var notFoundImages = []string{
	"https://storage.googleapis.com/acourse/static/9961f3c1-575f-4b98-af4f-447566ee1cb3.png",
	"https://storage.googleapis.com/acourse/static/b14a40c9-d3a4-465d-9453-ce7fcfbc594c.png",
}

// NotFound renders not found view
func NotFound(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	page := newPage(ctx)
	page.Image = notFoundImages[rand.Intn(len(notFoundImages))]

	data := struct {
		*Page
	}{page}

	w.Header().Set(header.XContentTypeOptions, "nosniff")
	renderWithStatusCode(ctx, w, http.StatusNotFound, tmplNotFound, &data)
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
	page.URL = BaseURL + "/course/" + url.PathEscape(course.Link())

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
