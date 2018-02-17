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
	"github.com/garyburd/redigo/redis"
	_ "github.com/lib/pq"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"gopkg.in/gomail.v2"

	"github.com/acoshift/acourse/app"
	"github.com/acoshift/acourse/view"
)

func main() {
	time.Local = time.UTC

	config := configfile.NewReader("config")

	// email location
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
	redisPool := &redis.Pool{
		MaxIdle:     5,
		IdleTimeout: 5 * time.Minute,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", config.String("redis_addr"), redis.DialPassword(config.String("redis_pass")))
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) > time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}

	// init cache pool
	// TODO: use in-memory redis for caching
	cachePool := redisPool

	// init databases
	db, err := sql.Open("postgres", config.String("sql_url"))
	if err != nil {
		log.Fatal(err)
	}

	view.BaseURL = config.String("base_url")

	app := app.New(app.Config{
		DB:            db,
		BaseURL:       config.String("base_url"),
		XSRFSecret:    config.String("xsrf_key"),
		RedisPool:     redisPool,
		RedisPrefix:   config.String("redis_prefix"),
		CachePool:     cachePool,
		CachePrefix:   config.String("redis_prefix"),
		SessionSecret: config.Bytes("session_secret"),
		Auth:          firAuth,
		Location:      loc,
		SlackURL:      config.String("slack_url"),
		EmailFrom:     config.String("email_from"),
		EmailDialer:   emailDialer,
		BucketHandle:  bucketHandle,
		BucketName:    config.String("bucket"),
	})

	err = hime.New().
		Handler(hime.Factory(app)).
		GracefulShutdown().
		ListenAndServe(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
