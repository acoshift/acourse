package controller

import (
	"net/http"

	"github.com/acoshift/go-firebase-admin"

	"github.com/acoshift/acourse/pkg/app"
)

// TODO: fixme
var sessName string

// New creates new app's controller
func New() app.Controller {
	return &ctrl{}
}

type ctrl struct {
	repo    app.Repository
	view    app.View
	firAuth firebase.Auth
}

func back(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, r.RequestURI, http.StatusSeeOther)
}
