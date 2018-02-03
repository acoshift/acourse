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

	"github.com/acoshift/acourse/pkg/app"
)

// New creates new view
func New(config Config) app.View {
	return &view{
		baseURL: config.BaseURL,
	}
}

// Config is the view config
type Config struct {
	BaseURL string
}

type view struct {
	baseURL string
}

// Page type provides layout data like title, description, and og
type Page struct {
	Title string
	Desc  string
	Image string
	URL   string
	Me    *app.User
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
	p.Me = app.GetUser(ctx)
	p.Flash = session.Get(ctx, "sess").Flash().Values()
	p.XSRF = template.HTML(`<input type="hidden" name="X" value="` + app.GetXSRFToken(ctx) + `">`)
	return &p
}

// Index renders index view
func (v *view) Index(w http.ResponseWriter, r *http.Request, courses []*app.Course) {
	ctx := r.Context()
	data := struct {
		*Page
		Courses []*app.Course
	}{newPage(ctx), courses}
	render(ctx, w, tmplIndex, &data)
}

var notFoundImages = []string{
	"https://storage.googleapis.com/acourse/static/9961f3c1-575f-4b98-af4f-447566ee1cb3.png",
	"https://storage.googleapis.com/acourse/static/b14a40c9-d3a4-465d-9453-ce7fcfbc594c.png",
}

// NotFound renders not found view
func (v *view) NotFound(w http.ResponseWriter, r *http.Request) {
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
func (v *view) SignIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data := struct {
		*Page
	}{newPage(ctx)}
	render(ctx, w, tmplSignIn, &data)
}

// SignInPassword renders signin-password view
func (v *view) SignInPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data := struct {
		*Page
	}{newPage(ctx)}
	render(ctx, w, tmplSignInPassword, &data)
}

// SignUp renders signup view
func (v *view) SignUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data := struct {
		*Page
	}{newPage(ctx)}
	render(ctx, w, tmplSignUp, &data)
}

// ResetPassword render reset password view
func (v *view) ResetPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	page := newPage(ctx)
	data := struct {
		*Page
	}{page}
	render(ctx, w, tmplResetPassword, &data)
}

// CheckEmail render check email view
func (v *view) CheckEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	page := newPage(ctx)
	data := struct {
		*Page
	}{page}
	render(ctx, w, tmplCheckEmail, &data)
}

// Profile renders profile view
func (v *view) Profile(w http.ResponseWriter, r *http.Request, ownCourses, enrolledCourses []*app.Course) {
	ctx := r.Context()
	page := newPage(ctx)
	me := app.GetUser(ctx)
	page.Title = me.Username + " | " + page.Title

	data := struct {
		*Page
		OwnCourses      []*app.Course
		EnrolledCourses []*app.Course
	}{page, ownCourses, enrolledCourses}
	render(ctx, w, tmplProfile, &data)
}

// ProfileEdit renders profile edit view
func (v *view) ProfileEdit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	me := app.GetUser(ctx)
	page := newPage(ctx)
	page.Title = me.Username + " | " + page.Title

	data := struct {
		*Page
	}{page}
	render(ctx, w, tmplProfileEdit, &data)
}

// Course renders course view
func (v *view) Course(w http.ResponseWriter, r *http.Request, course *app.Course, enrolled bool, owned bool, pendingEnroll bool) {
	ctx := r.Context()
	page := newPage(ctx)
	page.Title = course.Title + " | " + page.Title
	page.Desc = course.ShortDesc
	page.Image = course.Image
	page.URL = v.baseURL + "/course/" + url.PathEscape(course.Link())

	data := struct {
		*Page
		Course        *app.Course
		Enrolled      bool
		Owned         bool
		PendingEnroll bool
	}{page, course, enrolled, owned, pendingEnroll}
	render(ctx, w, tmplCourse, &data)
}

// CourseContent renders course content view
func (v *view) CourseContent(w http.ResponseWriter, r *http.Request, course *app.Course, content *app.CourseContent) {
	ctx := r.Context()
	page := newPage(ctx)
	page.Title = course.Title + " | " + page.Title
	page.Desc = course.ShortDesc
	page.Image = course.Image

	data := struct {
		*Page
		Course  *app.Course
		Content *app.CourseContent
	}{page, course, content}
	render(ctx, w, tmplCourseContent, &data)
}

// EditorCreate renders course create view
func (v *view) EditorCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	page := newPage(ctx)

	data := struct {
		*Page
	}{page}
	render(ctx, w, tmplEditorCreate, &data)
}

// EditorCourse renders course edit view
func (v *view) EditorCourse(w http.ResponseWriter, r *http.Request, course *app.Course) {
	ctx := r.Context()
	page := newPage(ctx)

	data := struct {
		*Page
		Course *app.Course
	}{page, course}
	render(ctx, w, tmplEditorCourse, &data)
}

// EditorContent renders editor content view
func (v *view) EditorContent(w http.ResponseWriter, r *http.Request, course *app.Course) {
	ctx := r.Context()
	page := newPage(ctx)

	data := struct {
		*Page
		Course *app.Course
	}{page, course}
	render(ctx, w, tmplEditorContent, &data)
}

// EditorContentCreate renders editor content create view
func (v *view) EditorContentCreate(w http.ResponseWriter, r *http.Request, course *app.Course) {
	ctx := r.Context()
	page := newPage(ctx)

	data := struct {
		*Page
		Course *app.Course
	}{page, course}
	render(ctx, w, tmplEditorContentCreate, &data)
}

// EditorContentEdit renders editor content edit view
func (v *view) EditorContentEdit(w http.ResponseWriter, r *http.Request, course *app.Course, content *app.CourseContent) {
	ctx := r.Context()
	page := newPage(ctx)

	data := struct {
		*Page
		Course  *app.Course
		Content *app.CourseContent
	}{page, course, content}
	render(ctx, w, tmplEditorContentEdit, &data)
}

// CourseEnroll renders course enroll view
func (v *view) CourseEnroll(w http.ResponseWriter, r *http.Request, course *app.Course) {
	ctx := r.Context()
	page := newPage(ctx)
	page.Title = course.Title + " | " + page.Title
	page.Desc = course.ShortDesc
	page.Image = course.Image
	page.URL = v.baseURL + "/course/" + url.PathEscape(course.Link())

	data := struct {
		*Page
		Course *app.Course
	}{page, course}
	render(ctx, w, tmplCourseEnroll, &data)
}

// Assignment render assignment view
func (v *view) Assignment(w http.ResponseWriter, r *http.Request, course *app.Course, assignments []*app.Assignment) {
	ctx := r.Context()
	page := newPage(ctx)
	page.Title = course.Title + " | " + page.Title
	page.Desc = course.ShortDesc
	page.Image = course.Image
	page.URL = v.baseURL + "/course/" + url.PathEscape(course.Link())

	data := struct {
		*Page
		Course      *app.Course
		Assignments []*app.Assignment
	}{page, course, assignments}
	render(ctx, w, tmplAssignment, &data)
}

// AdminUsers renders admin users view
func (v *view) AdminUsers(w http.ResponseWriter, r *http.Request, users []*app.User, currentPage, totalPage int) {
	ctx := r.Context()
	page := newPage(ctx)

	data := struct {
		*Page
		Users       []*app.User
		CurrentPage int
		TotalPage   int
	}{page, users, currentPage, totalPage}
	render(ctx, w, tmplAdminUsers, &data)
}

// AdminCourses renders admin courses view
func (v *view) AdminCourses(w http.ResponseWriter, r *http.Request, courses []*app.Course, currentPage, totalPage int) {
	ctx := r.Context()
	page := newPage(ctx)

	data := struct {
		*Page
		Courses     []*app.Course
		CurrentPage int
		TotalPage   int
	}{page, courses, currentPage, totalPage}
	render(ctx, w, tmplAdminCourses, &data)
}

// AdminPayments renders admin payments view
func (v *view) AdminPayments(w http.ResponseWriter, r *http.Request, payments []*app.Payment, currentPage, totalPage int) {
	ctx := r.Context()
	page := newPage(ctx)

	data := struct {
		*Page
		Payments    []*app.Payment
		CurrentPage int
		TotalPage   int
	}{page, payments, currentPage, totalPage}
	render(ctx, w, tmplAdminPayments, &data)
}

// AdminPaymentReject renders admin payment reject view
func (v *view) AdminPaymentReject(w http.ResponseWriter, r *http.Request, payment *app.Payment, message string) {
	ctx := r.Context()
	page := newPage(ctx)

	data := struct {
		*Page
		Payment *app.Payment
		Message string
	}{page, payment, message}
	render(ctx, w, tmplAdminPaymentReject, &data)
}
