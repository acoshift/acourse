package view

import "net/http"

type viewKey int

const (
	_ viewKey = iota
	keyIndex
	keySignIn
	keySignUp
	keyProfile
	keyProfileEdit
	keyUser
	keyCourse
	keyCourseEdit
)

// Index renders index view
func Index(w http.ResponseWriter, r *http.Request, data *IndexData) {
	render(w, r, keyIndex, data)
}

// SignIn renders signin view
func SignIn(w http.ResponseWriter, r *http.Request, data *AuthData) {
	render(w, r, keySignIn, data)
}

// SignUp renders signup view
func SignUp(w http.ResponseWriter, r *http.Request, data *AuthData) {
	render(w, r, keySignUp, data)
}

// Profile renders profile view
func Profile(w http.ResponseWriter, r *http.Request, data *ProfileData) {
	render(w, r, keyProfile, data)
}

// Course renders course view
func Course(w http.ResponseWriter, r *http.Request, data *CourseData) {
	render(w, r, keyCourse, data)
}
