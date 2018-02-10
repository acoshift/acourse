package app

import (
	"net/http"

	"github.com/acoshift/acourse/repository"
	"github.com/acoshift/acourse/view"
)

func index(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.URL.Path != "/" {
		view.NotFound(w, r)
		return
	}
	courses, err := repository.ListPublicCourses(ctx, cachePool, cachePrefix)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	view.Index(w, r, courses)
}
