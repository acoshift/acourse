package app

import (
	"net/http"
	"strings"

	"github.com/acoshift/acourse/appctx"
	"github.com/acoshift/acourse/view"
)

func fileHandler(name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, name)
	})
}

func courseHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := strings.SplitN(r.URL.Path, "/", 2)
		var p string
		if len(s) > 1 {
			p = strings.TrimSuffix(s[1], "/")
		}

		r = r.WithContext(appctx.NewCourseURLContext(r.Context(), s[0]))
		switch p {
		case "":
			courseView(w, r)
		case "content":
			mustSignedIn(http.HandlerFunc(courseContent)).ServeHTTP(w, r)
		case "enroll":
			mustSignedIn(http.HandlerFunc(courseEnroll)).ServeHTTP(w, r)
		case "assignment":
			mustSignedIn(http.HandlerFunc(courseAssignment)).ServeHTTP(w, r)
		default:
			view.NotFound(w, r)
		}
	})
}
