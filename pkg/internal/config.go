package internal

import (
	"context"
	"database/sql"
	"log"
	"time"

	"cloud.google.com/go/storage"
	"github.com/acoshift/configfile"
	"github.com/garyburd/redigo/redis"
	_ "github.com/lib/pq"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/identitytoolkit/v3"
	"google.golang.org/api/option"
	"gopkg.in/gomail.v2"
)

var config = configfile.NewReader("config")

var (
	redisAddr      = config.String("redis_addr")
	redisDB        = config.Int("redis_db")
	redisPass      = config.String("redis_pass")
	xsrfSecret     = config.String("xsrf_secret")
	serviceAccount = config.Bytes("service_account")
	baseURL        = config.String("base_url")
	sqlURL         = config.String("sql_url")
	bucket         = config.String("bucket")
	emailFrom      = config.String("email_from")
)

var (
	redisPool = &redis.Pool{
		IdleTimeout: 30 * time.Minute,
		MaxIdle:     10,
		MaxActive:   100,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redisAddr,
				redis.DialDatabase(redisDB),
				redis.DialPassword(redisPass),
			)
		},
	}
	gitClient    *identitytoolkit.RelyingpartyService
	db           *sql.DB
	bucketHandle *storage.BucketHandle
	emailDialer  = gomail.NewPlainDialer(
		config.String("email_server"),
		config.Int("email_port"),
		config.String("email_user"),
		config.String("email_password"),
	)
)

func init() {
	time.Local = time.UTC

	ctx := context.Background()

	gconf, err := google.JWTConfigFromJSON(serviceAccount, identitytoolkit.CloudPlatformScope, storage.ScopeReadWrite)
	if err != nil {
		log.Fatal(err)
	}
	gitService, err := identitytoolkit.New(gconf.Client(ctx))
	if err != nil {
		log.Fatal(err)
	}
	gitClient = gitService.Relyingparty

	storageClient, err := storage.NewClient(ctx, option.WithTokenSource(gconf.TokenSource(ctx)))
	if err != nil {
		log.Fatal(err)
	}
	bucketHandle = storageClient.Bucket(bucket)

	db, err = sql.Open("postgres", sqlURL)
	if err != nil {
		log.Fatal(err)
	}
}

// GetDB returns sql db
func GetDB() *sql.DB {
	return db
}

// GetRedisDB returns redis connection from pool
func GetRedisDB() redis.Conn {
	return redisPool.Get()
}

// GetRedisPool returns redis pool
func GetRedisPool() *redis.Pool {
	return redisPool
}

// GetXSRFSecret returns xsrf secret
func GetXSRFSecret() string {
	return xsrfSecret
}

// GetBaseURL returns base url
func GetBaseURL() string {
	return baseURL
}
