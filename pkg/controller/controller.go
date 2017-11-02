package controller

import (
	"net/http"
	"time"

	"cloud.google.com/go/storage"
	"github.com/acoshift/go-firebase-admin"
	"github.com/garyburd/redigo/redis"
	"gopkg.in/gomail.v2"

	"github.com/acoshift/acourse/pkg/app"
)

// New creates new app's controller
func New() app.Controller {
	return &ctrl{}
}

type ctrl struct {
	repo         app.Repository
	view         app.View
	firAuth      firebase.Auth
	loc          *time.Location
	slackURL     string
	emailFrom    string
	emailDialer  *gomail.Dialer
	baseURL      string
	cachePool    *redis.Pool
	cachePrefix  string
	bucketHandle *storage.BucketHandle
	bucketName   string
}

func back(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, r.RequestURI, http.StatusSeeOther)
}
