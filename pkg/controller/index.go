package controller

import (
	"net/http"

	"github.com/acoshift/acourse/pkg/view"
)

func (c *ctrl) Index(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.URL.Path != "/" {
		view.NotFound(w, r)
		return
	}
	courses, err := c.repo.ListPublicCourses(ctx, cachePool, cachePrefix)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	view.Index(w, r, courses)
}
