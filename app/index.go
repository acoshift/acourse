package app

import (
	"net/http"

	"github.com/acoshift/acourse/repository"
	"github.com/acoshift/acourse/view"
)

func index(w http.ResponseWriter, r *http.Request) {
	courses, err := repository.ListPublicCourses(db, cachePool, cachePrefix)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	view.Index(w, r, courses)
}
