package app

import "net/http"

// View is the app's view renderer
type View interface {
	Index(w http.ResponseWriter, r *http.Request, courses []*Course)
	NotFound(w http.ResponseWriter, r *http.Request)
	SignInPassword(w http.ResponseWriter, r *http.Request)
	SignUp(w http.ResponseWriter, r *http.Request)
	ResetPassword(w http.ResponseWriter, r *http.Request)
	Profile(w http.ResponseWriter, r *http.Request, ownCourses, enrolledCourses []*Course)
	ProfileEdit(w http.ResponseWriter, r *http.Request)
	Course(w http.ResponseWriter, r *http.Request, course *Course, enrolled bool, owned bool, pendingEnroll bool)
	CourseContent(w http.ResponseWriter, r *http.Request, course *Course, content *CourseContent)
	EditorCreate(w http.ResponseWriter, r *http.Request)
	EditorCourse(w http.ResponseWriter, r *http.Request, course *Course)
	EditorContent(w http.ResponseWriter, r *http.Request, course *Course)
	EditorContentCreate(w http.ResponseWriter, r *http.Request, course *Course)
	EditorContentEdit(w http.ResponseWriter, r *http.Request, course *Course, content *CourseContent)
	CourseEnroll(w http.ResponseWriter, r *http.Request, course *Course)
	Assignment(w http.ResponseWriter, r *http.Request, course *Course, assignments []*Assignment)
	AdminUsers(w http.ResponseWriter, r *http.Request, users []*User, currentPage, totalPage int)
	AdminCourses(w http.ResponseWriter, r *http.Request, courses []*Course, currentPage, totalPage int)
	AdminPayments(w http.ResponseWriter, r *http.Request, payments []*Payment, currentPage, totalPage int)
	AdminPaymentReject(w http.ResponseWriter, r *http.Request, payment *Payment, message string)
}
