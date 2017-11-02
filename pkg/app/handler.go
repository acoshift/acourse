package app

import (
	"net/http"
	"os"

	"github.com/acoshift/header"
	"github.com/acoshift/middleware"
)

// TODO: fixme
var ctrl Controller

// Handler returns app's handlers
func Handler() http.Handler {
	mux := http.NewServeMux()

	editor := http.NewServeMux()
	editor.Handle("/create", onlyInstructor(http.HandlerFunc(ctrl.EditorCreate)))
	editor.Handle("/course", isCourseOwner(http.HandlerFunc(ctrl.EditorCourse)))
	editor.Handle("/content", isCourseOwner(http.HandlerFunc(ctrl.EditorContent)))
	editor.Handle("/content/create", isCourseOwner(http.HandlerFunc(ctrl.EditorContentCreate)))
	editor.Handle("/content/edit", http.HandlerFunc(ctrl.EditorContentEdit))

	admin := http.NewServeMux()
	admin.Handle("/users", http.HandlerFunc(ctrl.AdminUsers))
	admin.Handle("/courses", http.HandlerFunc(ctrl.AdminCourses))
	admin.Handle("/payments/pending", http.HandlerFunc(ctrl.AdminPendingPayments))
	admin.Handle("/payments/history", http.HandlerFunc(ctrl.AdminHistoryPayments))
	admin.Handle("/payments/reject", http.HandlerFunc(ctrl.AdminRejectPayment))

	main := http.NewServeMux()
	main.Handle("/", http.HandlerFunc(ctrl.Index))
	main.Handle("/signin", mustNotSignedIn(http.HandlerFunc(ctrl.SignIn)))
	main.Handle("/openid", mustNotSignedIn(http.HandlerFunc(ctrl.OpenID)))
	main.Handle("/openid/callback", mustNotSignedIn(http.HandlerFunc(ctrl.OpenIDCallback)))
	main.Handle("/signup", mustNotSignedIn(http.HandlerFunc(ctrl.SignUp)))
	main.Handle("/signout", http.HandlerFunc(ctrl.SignOut))
	main.Handle("/reset/password", mustNotSignedIn(http.HandlerFunc(ctrl.ResetPassword)))
	main.Handle("/profile", mustSignedIn(http.HandlerFunc(ctrl.Profile)))
	main.Handle("/profile/edit", mustSignedIn(http.HandlerFunc(ctrl.ProfileEdit)))
	main.Handle("/course/", http.StripPrefix("/course/", http.HandlerFunc(ctrl.Course)))
	main.Handle("/admin/", http.StripPrefix("/admin", onlyAdmin(admin)))
	main.Handle("/editor/", http.StripPrefix("/editor", editor))

	mux.Handle("/", Middleware(main))
	mux.Handle("/~/", http.StripPrefix("/~", cache(http.FileServer(&fileFS{http.Dir("static")}))))
	mux.Handle("/favicon.ico", fileHandler("static/favicon.ico"))

	return middleware.Chain(
		setHeaders,
	)(mux)
}

type fileFS struct {
	http.FileSystem
}

func (fs *fileFS) Open(name string) (http.File, error) {
	f, err := fs.FileSystem.Open(name)
	if err != nil {
		return nil, err
	}
	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if stat.IsDir() {
		return nil, os.ErrNotExist
	}
	return f, nil
}

func cache(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(header.CacheControl, "public, max-age=31536000")
		h.ServeHTTP(w, r)
	})
}

func fileHandler(name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, name)
	})
}
