package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	"cloud.google.com/go/storage"
	"github.com/acoshift/configfile"
	"github.com/acoshift/go-firebase-admin"
	"github.com/acoshift/hime"
	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"gopkg.in/gomail.v2"

	"github.com/acoshift/acourse/app"
)

func main() {
	time.Local = time.UTC

	config := configfile.NewReader("config")

	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// init google cloud config
	gconf, err := google.JWTConfigFromJSON(config.Bytes("service_account"), storage.ScopeReadWrite)
	if err != nil {
		log.Fatal(err)
	}

	firApp, err := firebase.InitializeApp(ctx, firebase.AppOptions{
		ProjectID: config.String("project_id"),
	}, option.WithCredentialsFile("config/service_account"))
	if err != nil {
		log.Fatal(err)
	}
	firAuth := firApp.Auth()

	// init google storage
	storageClient, err := storage.NewClient(ctx, option.WithTokenSource(gconf.TokenSource(ctx)))
	if err != nil {
		log.Fatal(err)
	}
	bucketHandle := storageClient.Bucket(config.String("bucket"))

	// init email client
	emailDialer := gomail.NewPlainDialer(config.String("email_server"), config.Int("email_port"), config.String("email_user"), config.String("email_password"))

	// init redis pool
	redisClient := redis.NewClient(&redis.Options{
		MaxRetries:  3,
		PoolSize:    5,
		IdleTimeout: 60 * time.Minute,
		Addr:        config.String("redis_addr"),
		Password:    config.String("redis_pass"),
	})

	// init databases
	db, err := sql.Open("postgres", config.String("sql_url"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(4)

	himeApp := hime.New()

	h := app.New(himeApp, app.Config{
		DB:            db,
		BaseURL:       config.String("base_url"),
		RedisClient:   redisClient,
		RedisPrefix:   config.String("redis_prefix"),
		SessionSecret: config.Bytes("session_secret"),
		Auth:          firAuth,
		Location:      loc,
		SlackURL:      config.String("slack_url"),
		EmailFrom:     config.String("email_from"),
		EmailDialer:   emailDialer,
		BucketHandle:  bucketHandle,
		BucketName:    config.String("bucket"),
	})

	err = himeApp.
		Handler(h).
		Address(":8080").
		GracefulShutdown().
		Wait(5 * time.Second).
		Timeout(10 * time.Second).
		ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
