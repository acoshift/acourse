package app

import (
	"context"
	"database/sql"
	"encoding/gob"
	"time"

	"cloud.google.com/go/storage"
	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/acourse/pkg/view"
	"github.com/acoshift/go-firebase-admin"
	"github.com/garyburd/redigo/redis"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"gopkg.in/gomail.v2"
)

// app shared vars
var (
	bucketName    string
	bucketHandle  *storage.BucketHandle
	emailDialer   *gomail.Dialer
	emailFrom     string
	baseURL       string
	xsrfSecret    string
	db            *sql.DB
	firAuth       *admin.Auth
	redisAddr     string
	redisPass     string
	redisPrefix   string
	slackURL      string
	sessionSecret []byte
	loc           *time.Location
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
	RedisAddr      string
	RedisPass      string
	RedisPrefix    string
	SessionSecret  []byte
	SlackURL       string
}

func init() {
	gob.Register(sessionKey(0))
}

// Init inits app package with given config
func Init(config Config) error {
	ctx := context.Background()

	var err error
	loc, err = time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return err
	}

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
	redisAddr = config.RedisAddr
	redisPass = config.RedisPass
	redisPrefix = config.RedisPrefix
	slackURL = config.SlackURL
	sessionSecret = config.SessionSecret

	// init databases
	db, err = sql.Open("postgres", config.SQLURL)
	if err != nil {
		return err
	}
	db.SetMaxIdleConns(10)

	// TODO: use in-memory redis for caching
	redisPool := &redis.Pool{
		MaxIdle:     100,
		IdleTimeout: 10 * time.Minute,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redisAddr, redis.DialPassword(redisPass))
		},
	}

	// init other packages
	err = model.Init(model.Config{DB: db, RedisPool: redisPool, RedisPrefix: redisPrefix})
	if err != nil {
		return err
	}
	err = view.Init(view.Config{BaseURL: baseURL})
	if err != nil {
		return err
	}

	return nil
}
