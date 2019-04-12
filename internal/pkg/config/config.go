package config

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/errorreporting"
	"cloud.google.com/go/storage"
	"github.com/acoshift/configfile"
	"github.com/acoshift/go-firebase-admin"
	"google.golang.org/api/option"
)

var cfg = configfile.NewReader("config")

var (
	StringDefault   = cfg.StringDefault
	String          = cfg.String
	IntDefault      = cfg.IntDefault
	Int             = cfg.Int
	DurationDefault = cfg.DurationDefault
	Bytes           = cfg.Bytes
)

var (
	firebaseApp   *firebase.App
	storageClient *storage.Client
	errorClient   *errorreporting.Client
	location      *time.Location
)

func init() {
	time.Local = time.UTC

	var err error

	location, err = time.LoadLocation(StringDefault("location", "Asia/Bangkok"))
	must(err)

	ctx := context.Background()
	var googleClientOpts []option.ClientOption

	if len(cfg.Bytes("service_account")) > 0 {
		googleClientOpts = append(googleClientOpts, option.WithCredentialsFile("config/service_account"))
	}

	serviceName := StringDefault("service", "acourse")
	projectID := String("project_id")

	// init error reporting, ignore error
	errorClient, _ = errorreporting.NewClient(ctx, projectID, errorreporting.Config{
		ServiceName: serviceName,
		OnError: func(err error) {
			log.Println(err)
		},
	}, googleClientOpts...)

	firebaseApp, err = firebase.InitializeApp(ctx, firebase.AppOptions{
		ProjectID: projectID,
	}, googleClientOpts...)
	must(err)

	// init google storage
	storageClient, err = storage.NewClient(ctx, googleClientOpts...)
	must(err)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func ErrorClient() *errorreporting.Client {
	return errorClient
}

func StorageClient() *storage.Client {
	return storageClient
}

func FirebaseApp() *firebase.App {
	return firebaseApp
}

func Location() *time.Location {
	return location
}
