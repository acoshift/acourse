package view

import (
	"context"
	"html/template"
	"net/http"
	"net/url"

	"github.com/acoshift/acourse/pkg/appctx"
	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/flash"
)

type (
	keyIndex               struct{}
	keySignIn              struct{}
	keySignUp              struct{}
	keyProfile             struct{}
	keyProfileEdit         struct{}
	keyUser                struct{}
	keyCourse              struct{}
	keyCourseContent       struct{}
	keyAssignment          struct{}
	keyEditorCreate        struct{}
	keyEditorCourse        struct{}
	keyEditorContent       struct{}
	keyEditorContentCreate struct{}
	keyEditorContentEdit   struct{}
	keyCourseEnroll        struct{}
	keyAdminUsers          struct{}
	keyAdminCourses        struct{}
	keyAdminPayments       struct{}
	keyAdminPaymentReject  struct{}
)

// Page type provides layout data like title, description, and og
type Page struct {
	Title string
	Desc  string
	Image string
	URL   string
	Me    *model.User
	Flash flash.Flash
	XSRF  template.HTML
}

var defaultPage = Page{
	Title: "Acourse",
	Desc:  "Online courses for everyone",
	Image: "https://storage.googleapis.com/acourse/static/62b9eb0e-3668-4f9f-86b7-a11349938f7a.jpg",
	URL:   "https://acourse.io",
}

func newPage(ctx context.Context) *Page {
	p := defaultPage
	p.Me = appctx.GetUser(ctx)
	p.Flash = flash.Get(ctx)
	p.XSRF = template.HTML(`<input type="hidden" name="X" value="` + appctx.GetXSRFToken(ctx) + `">`)
	return &p
}

// Index renders index view
func Index(w http.ResponseWriter, r *http.Request, courses []*model.Course) {
	ctx := r.Context()
	data := struct {
		*Page
		Courses []*model.Course
	}{newPage(ctx), courses}
	render(ctx, w, keyIndex{}, &data)
}

// SignIn renders signin view
func SignIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data := struct {
		*Page
	}{newPage(ctx)}
	render(ctx, w, keySignIn{}, &data)
}

// SignUp renders signup view
func SignUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data := struct {
		*Page
	}{newPage(ctx)}
	render(ctx, w, keySignUp{}, &data)
}

// Profile renders profile view
func Profile(w http.ResponseWriter, r *http.Request, ownCourses, enrolledCourses []*model.Course) {
	ctx := r.Context()
	page := newPage(ctx)
	me := appctx.GetUser(ctx)
	page.Title = me.Username + " | " + page.Title

	data := struct {
		*Page
		OwnCourses      []*model.Course
		EnrolledCourses []*model.Course
	}{page, ownCourses, enrolledCourses}
	render(ctx, w, keyProfile{}, &data)
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
	render(ctx, w, keyProfileEdit{}, &data)
}

// Course renders course view
func Course(w http.ResponseWriter, r *http.Request, course *model.Course, enrolled bool, owned bool, pendingEnroll bool) {
	ctx := r.Context()
	page := newPage(ctx)
	page.Title = course.Title + " | " + page.Title
	page.Desc = course.ShortDesc
	page.Image = course.Image
	page.URL = baseURL + "/course/" + url.PathEscape(course.Link())

	data := struct {
		*Page
		Course        *model.Course
		Enrolled      bool
		Owned         bool
		PendingEnroll bool
	}{page, course, enrolled, owned, pendingEnroll}
	render(ctx, w, keyCourse{}, &data)
}

// CourseContent renders course content view
func CourseContent(w http.ResponseWriter, r *http.Request, course *model.Course, content *model.CourseContent) {
	ctx := r.Context()
	page := newPage(ctx)
	page.Title = course.Title + " | " + page.Title
	page.Desc = course.ShortDesc
	page.Image = course.Image

	data := struct {
		*Page
		Course  *model.Course
		Content *model.CourseContent
	}{page, course, content}
	render(ctx, w, keyCourseContent{}, &data)
}

// EditorCreate renders course create view
func EditorCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	page := newPage(ctx)

	data := struct {
		*Page
	}{page}
	render(ctx, w, keyEditorCreate{}, &data)
}

// EditorCourse renders course edit view
func EditorCourse(w http.ResponseWriter, r *http.Request, course *model.Course) {
	ctx := r.Context()
	page := newPage(ctx)

	data := struct {
		*Page
		Course *model.Course
	}{page, course}
	render(ctx, w, keyEditorCourse{}, &data)
}

// EditorContent renders editor content view
func EditorContent(w http.ResponseWriter, r *http.Request, course *model.Course) {
	ctx := r.Context()
	page := newPage(ctx)

	data := struct {
		*Page
		Course *model.Course
	}{page, course}
	render(ctx, w, keyEditorContent{}, &data)
}

// EditorContentCreate renders editor content create view
func EditorContentCreate(w http.ResponseWriter, r *http.Request, course *model.Course) {
	ctx := r.Context()
	page := newPage(ctx)

	data := struct {
		*Page
		Course *model.Course
	}{page, course}
	render(ctx, w, keyEditorContentCreate{}, &data)
}

// EditorContentEdit renders editor content edit view
func EditorContentEdit(w http.ResponseWriter, r *http.Request, course *model.Course, content *model.CourseContent) {
	ctx := r.Context()
	page := newPage(ctx)

	data := struct {
		*Page
		Course  *model.Course
		Content *model.CourseContent
	}{page, course, content}
	render(ctx, w, keyEditorContentEdit{}, &data)
}

// CourseEnroll renders course enroll view
func CourseEnroll(w http.ResponseWriter, r *http.Request, course *model.Course) {
	ctx := r.Context()
	page := newPage(ctx)
	page.Title = course.Title + " | " + page.Title
	page.Desc = course.ShortDesc
	page.Image = course.Image
	page.URL = baseURL + "/course/" + url.PathEscape(course.Link())

	data := struct {
		*Page
		Course *model.Course
	}{page, course}
	render(ctx, w, keyCourseEnroll{}, &data)
}

// Assignment render assignment view
func Assignment(w http.ResponseWriter, r *http.Request, course *model.Course, assignments []*model.Assignment) {
	ctx := r.Context()
	page := newPage(ctx)
	page.Title = course.Title + " | " + page.Title
	page.Desc = course.ShortDesc
	page.Image = course.Image
	page.URL = baseURL + "/course/" + url.PathEscape(course.Link())

	data := struct {
		*Page
		Course      *model.Course
		Assignments []*model.Assignment
	}{page, course, assignments}
	render(ctx, w, keyAssignment{}, &data)
}

// AdminUsers renders admin users view
func AdminUsers(w http.ResponseWriter, r *http.Request, users []*model.User, currentPage, totalPage int) {
	ctx := r.Context()
	page := newPage(ctx)

	data := struct {
		*Page
		Users       []*model.User
		CurrentPage int
		TotalPage   int
	}{page, users, currentPage, totalPage}
	render(ctx, w, keyAdminUsers{}, &data)
}

// AdminCourses renders admin courses view
func AdminCourses(w http.ResponseWriter, r *http.Request, courses []*model.Course, currentPage, totalPage int) {
	ctx := r.Context()
	page := newPage(ctx)

	data := struct {
		*Page
		Courses     []*model.Course
		CurrentPage int
		TotalPage   int
	}{page, courses, currentPage, totalPage}
	render(ctx, w, keyAdminCourses{}, &data)
}

// AdminPayments renders admin payments view
func AdminPayments(w http.ResponseWriter, r *http.Request, payments []*model.Payment, currentPage, totalPage int) {
	ctx := r.Context()
	page := newPage(ctx)

	data := struct {
		*Page
		Payments    []*model.Payment
		CurrentPage int
		TotalPage   int
	}{page, payments, currentPage, totalPage}
	render(ctx, w, keyAdminPayments{}, &data)
}

// AdminPaymentReject renders admin payment reject view
func AdminPaymentReject(w http.ResponseWriter, r *http.Request, payment *model.Payment, message string) {
	ctx := r.Context()
	page := newPage(ctx)

	data := struct {
		*Page
		Payment *model.Payment
		Message string
	}{page, payment, message}
	render(ctx, w, keyAdminPaymentReject{}, &data)
}
