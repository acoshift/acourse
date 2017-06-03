package view

import (
	"net/http"

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
	keyEditorCreate        struct{}
	keyEditorCourse        struct{}
	keyEditorContent       struct{}
	keyEditorContentCreate struct{}
	keyEditorContentEdit   struct{}
	keyCourseEnroll        struct{}
	keyAdminUsers          struct{}
	keyAdminCourses        struct{}
	keyAdminPayments       struct{}
)

var defaultPage = Page{
	Title: "Acourse",
	Desc:  "Online courses for everyone",
	Image: "https://storage.googleapis.com/acourse/static/62b9eb0e-3668-4f9f-86b7-a11349938f7a.jpg",
	URL:   "https://acourse.io",
}

// Index renders index view
func Index(w http.ResponseWriter, r *http.Request, courses []*model.Course) {
	data := struct {
		Page    *Page
		Courses []*model.Course
	}{&defaultPage, courses}
	render(w, r, keyIndex{}, &data)
}

// SignIn renders signin view
func SignIn(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Page  *Page
		Flash flash.Flash
	}{&defaultPage, flash.Get(r.Context())}
	render(w, r, keySignIn{}, &data)
}

// SignUp renders signup view
func SignUp(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Page  *Page
		Flash flash.Flash
	}{&defaultPage, flash.Get(r.Context())}
	render(w, r, keySignUp{}, &data)
}

// Profile renders profile view
func Profile(w http.ResponseWriter, r *http.Request, ownCourses, enrolledCourses []*model.Course) {
	me := appctx.GetUser(r.Context())
	page := defaultPage
	page.Title = me.Username + " | " + page.Title

	data := struct {
		Page            *Page
		Navbar          string
		Me              *model.User
		OwnCourses      []*model.Course
		EnrolledCourses []*model.Course
	}{&page, "profile", me, ownCourses, enrolledCourses}
	render(w, r, keyProfile{}, &data)
}

// ProfileEdit renders profile edit view
func ProfileEdit(w http.ResponseWriter, r *http.Request, data *ProfileEditData) {
	render(w, r, keyProfileEdit{}, data)
}

// Course renders course view
func Course(w http.ResponseWriter, r *http.Request, data *CourseData) {
	render(w, r, keyCourse{}, data)
}

// EditorCreate renders course create view
func EditorCreate(w http.ResponseWriter, r *http.Request, data *CourseCreateData) {
	render(w, r, keyEditorCreate{}, data)
}

// EditorCourse renders course edit view
func EditorCourse(w http.ResponseWriter, r *http.Request, data *CourseEditData) {
	render(w, r, keyEditorCourse{}, data)
}

// EditorContent renders editor content view
func EditorContent(w http.ResponseWriter, r *http.Request, data *CourseEditData) {
	render(w, r, keyEditorContent{}, data)
}

// EditorContentCreate renders editor content create view
func EditorContentCreate(w http.ResponseWriter, r *http.Request, data *EditorContentCreateData) {
	data.Page = &defaultPage
	render(w, r, keyEditorContentCreate{}, data)
}

// EditorContentEdit renders editor content edit view
func EditorContentEdit(w http.ResponseWriter, r *http.Request, data *EditorContentCreateData) {
	data.Page = &defaultPage
	render(w, r, keyEditorContentEdit{}, data)
}

// CourseEnroll renders course enroll view
func CourseEnroll(w http.ResponseWriter, r *http.Request, data *CourseData) {
	render(w, r, keyCourseEnroll{}, data)
}

// AdminUsers renders admin users view
func AdminUsers(w http.ResponseWriter, r *http.Request, data *AdminUsersData) {
	render(w, r, keyAdminUsers{}, data)
}

// AdminCourses renders admin courses view
func AdminCourses(w http.ResponseWriter, r *http.Request, data *AdminCoursesData) {
	render(w, r, keyAdminCourses{}, data)
}

// AdminPayments renders admin payments view
func AdminPayments(w http.ResponseWriter, r *http.Request, data *AdminPaymentsData) {
	render(w, r, keyAdminPayments{}, data)
}
