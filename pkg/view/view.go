package view

import "net/http"

type (
	keyIndex         struct{}
	keySignIn        struct{}
	keySignUp        struct{}
	keyProfile       struct{}
	keyProfileEdit   struct{}
	keyUser          struct{}
	keyCourse        struct{}
	keyEditorCreate  struct{}
	keyEditorCourse  struct{}
	keyEditorContent struct{}
	keyCourseEnroll  struct{}
	keyAdminUsers    struct{}
	keyAdminCourses  struct{}
	keyAdminPayments struct{}
)

// Index renders index view
func Index(w http.ResponseWriter, r *http.Request, data *IndexData) {
	render(w, r, keyIndex{}, data)
}

// SignIn renders signin view
func SignIn(w http.ResponseWriter, r *http.Request, data *AuthData) {
	render(w, r, keySignIn{}, data)
}

// SignUp renders signup view
func SignUp(w http.ResponseWriter, r *http.Request, data *AuthData) {
	render(w, r, keySignUp{}, data)
}

// Profile renders profile view
func Profile(w http.ResponseWriter, r *http.Request, data *ProfileData) {
	render(w, r, keyProfile{}, data)
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

// EditorContents renders course content edit view
func EditorContents(w http.ResponseWriter, r *http.Request, data *CourseEditData) {
	render(w, r, keyEditorContent{}, data)
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
