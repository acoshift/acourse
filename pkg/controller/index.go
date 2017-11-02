package controller

import (
	"net/http"
)

func (c *ctrl) Index(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.URL.Path != "/" {
		c.view.NotFound(w, r)
		return
	}
	courses, err := c.repo.ListPublicCourses(ctx, c.cachePool, c.cachePrefix)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	c.view.Index(w, r, courses)
}
