package internal

import (
	"context"
	"log"
	"time"

	"github.com/acoshift/configfile"
	"github.com/garyburd/redigo/redis"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/identitytoolkit/v3"
)

var config = configfile.NewReader("config")

var (
	redisPrimaryAddr   = config.String("redis_primary_addr")
	redisPrimaryDB     = config.Int("redis_primary_db")
	redisPrimaryPass   = config.String("redis_primary_pass")
	redisSecondaryAddr = config.String("redis_secondary_addr")
	redisSecondaryDB   = config.Int("redis_secondary_db")
	redisSecondaryPass = config.String("redis_secondary_pass")
	xsrfSecret         = config.String("xsrf_secret")
	serviceAccount     = config.Bytes("service_account")
)

var (
	primaryPool = &redis.Pool{
		IdleTimeout: 10 * time.Minute,
		MaxIdle:     10,
		MaxActive:   100,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redisPrimaryAddr,
				redis.DialDatabase(redisPrimaryDB),
				redis.DialPassword(redisPrimaryPass),
			)
		},
	}
	secondaryPool = &redis.Pool{
		IdleTimeout: 10 * time.Minute,
		MaxIdle:     10,
		MaxActive:   100,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redisSecondaryAddr,
				redis.DialDatabase(redisSecondaryDB),
				redis.DialPassword(redisSecondaryPass),
			)
		},
	}
	gitClient *identitytoolkit.RelyingpartyService
)

func init() {
	time.Local = time.UTC

	ctx := context.Background()

	scopes := []string{
		identitytoolkit.CloudPlatformScope,
		identitytoolkit.FirebaseScope,
	}

	gconf, err := google.JWTConfigFromJSON(serviceAccount, scopes...)
	if err != nil {
		log.Fatal(err)
	}
	gitService, err := identitytoolkit.New(gconf.Client(ctx))
	if err != nil {
		log.Fatal(err)
	}
	gitClient = gitService.Relyingparty
}

// GetPrimaryDB returns primary redis connection from pool, use for store app data
func GetPrimaryDB() redis.Conn {
	return primaryPool.Get()
}

// GetSecondaryDB returns secondary redis connection from pool, use for store session
func GetSecondaryDB() redis.Conn {
	return secondaryPool.Get()
}

// GetSecondaryPool returns secondary redis pool
func GetSecondaryPool() *redis.Pool {
	return secondaryPool
}

// GetXSRFSecret returns xsrf secret
func GetXSRFSecret() string {
	return xsrfSecret
}

// SignInUser sign in user
func SignInUser(email, password string) (string, error) {
	req := gitClient.VerifyPassword(&identitytoolkit.IdentitytoolkitRelyingpartyVerifyPasswordRequest{
		Email:    email,
		Password: password,
	})
	resp, err := req.Do()
	if err != nil {
		return "", err
	}
	return resp.LocalId, nil
}
