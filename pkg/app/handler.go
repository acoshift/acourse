package app

import (
	"net/http"
	"os"
	"strings"

	"github.com/acoshift/header"
)

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

func courseHandler(ctrl Controller, view View) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := strings.SplitN(r.URL.Path, "/", 2)
		var p string
		if len(s) > 1 {
			p = strings.TrimSuffix(s[1], "/")
		}

		r = r.WithContext(WithCourseURL(r.Context(), s[0]))
		switch p {
		case "":
			ctrl.CourseView(w, r)
		case "content":
			mustSignedIn(http.HandlerFunc(ctrl.CourseContent)).ServeHTTP(w, r)
		case "enroll":
			mustSignedIn(http.HandlerFunc(ctrl.CourseEnroll)).ServeHTTP(w, r)
		case "assignment":
			mustSignedIn(http.HandlerFunc(ctrl.CourseAssignment)).ServeHTTP(w, r)
		default:
			view.NotFound(w, r)
		}
	})
}
