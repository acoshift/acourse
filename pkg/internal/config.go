package internal

import (
	"context"
	"crypto/rand"
	"database/sql"
	"log"
	"time"

	"github.com/acoshift/configfile"
	"github.com/garyburd/redigo/redis"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/identitytoolkit/v3"
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
	gitClient *identitytoolkit.RelyingpartyService
	db        *sql.DB
)

func init() {
	time.Local = time.UTC

	ctx := context.Background()

	gconf, err := google.JWTConfigFromJSON(serviceAccount, identitytoolkit.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}
	gitService, err := identitytoolkit.New(gconf.Client(ctx))
	if err != nil {
		log.Fatal(err)
	}
	gitClient = gitService.Relyingparty

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

// SignInUser sign in user with email and password
func SignInUser(email, password string) (string, error) {
	resp, err := gitClient.VerifyPassword(&identitytoolkit.IdentitytoolkitRelyingpartyVerifyPasswordRequest{
		Email:    email,
		Password: password,
	}).Do()
	if err != nil {
		return "", err
	}
	return resp.LocalId, nil
}

func generateSessionID() string {
	b := make([]byte, 24)
	rand.Read(b)
	return string(b)
}

// SignInUserProvider sign in user with open id provider
func SignInUserProvider(provider string) (redirectURI string, sessionID string, err error) {
	sessID := generateSessionID()
	resp, err := gitClient.CreateAuthUri(&identitytoolkit.IdentitytoolkitRelyingpartyCreateAuthUriRequest{
		ProviderId:   provider,
		ContinueUri:  baseURL + "/openid/callback",
		AuthFlowType: "CODE_FLOW",
		SessionId:    sessID,
	}).Do()
	if err != nil {
		return "", "", err
	}
	return resp.AuthUri, sessID, nil
}

// SignInUserProviderCallback sign in user with open id provider callback
func SignInUserProviderCallback(callbackURI string, sessID string) (string, error) {
	resp, err := gitClient.VerifyAssertion(&identitytoolkit.IdentitytoolkitRelyingpartyVerifyAssertionRequest{
		RequestUri: baseURL + callbackURI,
		SessionId:  sessID,
	}).Do()
	if err != nil {
		return "", err
	}
	return resp.LocalId, nil
}

// SignUpUser creates new user
func SignUpUser(email, password string) (string, error) {
	resp, err := gitClient.SignupNewUser(&identitytoolkit.IdentitytoolkitRelyingpartySignupNewUserRequest{
		Email:    email,
		Password: password,
	}).Do()
	if err != nil {
		return "", err
	}
	return resp.LocalId, nil
}

// GetVerifyEmailCode gets out-of-band confirmation code for verify email
func GetVerifyEmailCode(email string) (string, error) {
	resp, err := gitClient.GetOobConfirmationCode(&identitytoolkit.Relyingparty{
		Kind:        "identitytoolkit#relyingparty",
		RequestType: "VERIFY_EMAIL",
		Email:       email,
	}).Do()
	if err != nil {
		return "", err
	}
	return resp.OobCode, nil
}

// GetResetPasswordCode gets out-of-band confirmation code for reset password
func GetResetPasswordCode(email string) (string, error) {
	resp, err := gitClient.GetOobConfirmationCode(&identitytoolkit.Relyingparty{
		Kind:        "identitytoolkit#relyingparty",
		RequestType: "PASSWORD_RESET",
		Email:       email,
	}).Do()
	if err != nil {
		return "", err
	}
	return resp.OobCode, nil
}
