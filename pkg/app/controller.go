package app

import "net/http"

// Controller is the app's controller
type Controller interface {
	// Index
	Index(w http.ResponseWriter, r *http.Request)

	// Auth
	SignIn(w http.ResponseWriter, r *http.Request)
	SignInPassword(w http.ResponseWriter, r *http.Request)
	CheckEmail(w http.ResponseWriter, r *http.Request)
	OpenID(w http.ResponseWriter, r *http.Request)
	OpenIDCallback(w http.ResponseWriter, r *http.Request)
	SignUp(w http.ResponseWriter, r *http.Request)
	SignOut(w http.ResponseWriter, r *http.Request)
	ResetPassword(w http.ResponseWriter, r *http.Request)

	// Course
	CourseView(w http.ResponseWriter, r *http.Request)
	CourseContent(w http.ResponseWriter, r *http.Request)
	CourseEnroll(w http.ResponseWriter, r *http.Request)
	CourseAssignment(w http.ResponseWriter, r *http.Request)

	// Admin
	AdminUsers(w http.ResponseWriter, r *http.Request)
	AdminCourses(w http.ResponseWriter, r *http.Request)
	AdminPendingPayments(w http.ResponseWriter, r *http.Request)
	AdminHistoryPayments(w http.ResponseWriter, r *http.Request)
	AdminRejectPayment(w http.ResponseWriter, r *http.Request)

	// Profile
	Profile(w http.ResponseWriter, r *http.Request)
	ProfileEdit(w http.ResponseWriter, r *http.Request)

	// Editor
	EditorCreate(w http.ResponseWriter, r *http.Request)
	EditorCourse(w http.ResponseWriter, r *http.Request)
	EditorContent(w http.ResponseWriter, r *http.Request)
	EditorContentCreate(w http.ResponseWriter, r *http.Request)
	EditorContentEdit(w http.ResponseWriter, r *http.Request)
}
