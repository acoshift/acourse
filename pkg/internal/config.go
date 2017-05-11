package internal

import (
	"context"
	"crypto/rand"
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
	baseURL            = config.String("base_url")
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

	gconf, err := google.JWTConfigFromJSON(serviceAccount, identitytoolkit.CloudPlatformScope)
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
