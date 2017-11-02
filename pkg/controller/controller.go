package controller

import (
	"net/http"
	"time"

	"github.com/acoshift/go-firebase-admin"
	"gopkg.in/gomail.v2"

	"github.com/acoshift/acourse/pkg/app"
)

// New creates new app's controller
func New() app.Controller {
	return &ctrl{}
}

type ctrl struct {
	repo        app.Repository
	view        app.View
	firAuth     firebase.Auth
	loc         *time.Location
	slackURL    string
	emailFrom   string
	emailDialer *gomail.Dialer
	baseURL     string
}

func back(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, r.RequestURI, http.StatusSeeOther)
}
