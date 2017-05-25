package app

import (
	"context"
	"database/sql"
	"encoding/gob"

	"cloud.google.com/go/storage"
	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/acourse/pkg/view"
	"github.com/acoshift/go-firebase-admin"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"gopkg.in/gomail.v2"
)

// app shared vars
var (
	bucketName   string
	bucketHandle *storage.BucketHandle
	emailDialer  *gomail.Dialer
	emailFrom    string
	baseURL      string
	xsrfSecret   string
	db           *sql.DB
	firAuth      *admin.Auth
)

// Config use to init app package
type Config struct {
	ProjectID      string
	ServiceAccount []byte
	BucketName     string
	EmailServer    string
	EmailPort      int
	EmailUser      string
	EmailPassword  string
	EmailFrom      string
	BaseURL        string
	XSRFSecret     string
	SQLURL         string
}

func init() {
	gob.Register(sessionKey(0))
}

// Init inits app package with given config
func Init(config Config) error {
	ctx := context.Background()

	// init google cloud config
	gconf, err := google.JWTConfigFromJSON(config.ServiceAccount, storage.ScopeReadWrite)
	if err != nil {
		return err
	}

	firApp, err := admin.InitializeApp(ctx, admin.AppOptions{
		ProjectID:      config.ProjectID,
		ServiceAccount: config.ServiceAccount,
	})
	if err != nil {
		return err
	}
	firAuth = firApp.Auth()

	// init google storage
	storageClient, err := storage.NewClient(ctx, option.WithTokenSource(gconf.TokenSource(ctx)))
	if err != nil {
		return err
	}
	bucketName = config.BucketName
	bucketHandle = storageClient.Bucket(config.BucketName)

	// init email client
	emailDialer = gomail.NewPlainDialer(config.EmailServer, config.EmailPort, config.EmailUser, config.EmailPassword)
	emailFrom = config.EmailFrom

	baseURL = config.BaseURL
	xsrfSecret = config.XSRFSecret

	// init databases
	db, err = sql.Open("postgres", config.SQLURL)
	if err != nil {
		return err
	}

	// init other packages
	err = model.Init(model.Config{DB: db})
	if err != nil {
		return err
	}
	err = view.Init(view.Config{XSRFSecret: xsrfSecret})
	if err != nil {
		return err
	}

	return nil
}
