package controller

import (
	"net/http"
	"time"

	"cloud.google.com/go/storage"
	"github.com/acoshift/go-firebase-admin"
	"gopkg.in/gomail.v2"

	"github.com/acoshift/acourse/pkg/app"
)

// New creates new app's controller
func New(config Config) app.Controller {
	return &ctrl{
		repo:         config.Repository,
		auth:         config.Auth,
		loc:          config.Location,
		slackURL:     config.SlackURL,
		emailFrom:    config.EmailFrom,
		emailDialer:  config.EmailDialer,
		baseURL:      config.BaseURL,
		cachePrefix:  config.CachePrefix,
		bucketHandle: config.BucketHandle,
		bucketName:   config.BucketName,
	}
}

// Config is the controller config
type Config struct {
	Repository   app.Repository
	Auth         *firebase.Auth
	Location     *time.Location
	SlackURL     string
	EmailFrom    string
	EmailDialer  *gomail.Dialer
	BaseURL      string
	CachePrefix  string
	BucketHandle *storage.BucketHandle
	BucketName   string
}

type ctrl struct {
	repo         app.Repository
	auth         *firebase.Auth
	loc          *time.Location
	slackURL     string
	emailFrom    string
	emailDialer  *gomail.Dialer
	baseURL      string
	cachePrefix  string
	bucketHandle *storage.BucketHandle
	bucketName   string
}

func back(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, r.RequestURI, http.StatusSeeOther)
}
