package app

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"google.golang.org/grpc"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/logging"
	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"github.com/acoshift/go-firebase-admin"
	"github.com/acoshift/httperror"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/grpc/credentials"
	"gopkg.in/yaml.v2"
)

var (
	firAuth *admin.Auth

	tokenError = httperror.New(http.StatusUnauthorized, "token")
)

// Config type
type Config struct {
	Debug     bool   `yaml:"debug"`
	Port      string `yaml:"port"`
	TLSPort   string `yaml:"tlsPort"`
	TLSCert   string `yaml:"tlsCert"`
	TLSKey    string `yaml:"tlsKey"`
	Domain    string `yaml:"domain"`
	ProjectID string `yaml:"projectId"`
	Services  struct {
		User  string `yaml:"user"`
		Email string `yaml:"email"`
	} `yaml:"services"`
	Email struct {
		From     string `yaml:"from"`
		Server   string `yaml:"server"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	} `yaml:"email"`
	Firebase struct {
		ProjectID      string `yaml:"projectId"`
		ServiceAccount string `yaml:"serviceAccount"`
	} `yaml:"firebase"`
	ServiceAccount string `yaml:"serviceAccount"`
}

// LoadConfig loads config from file
func LoadConfig(filename string) (*Config, error) {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	cfg := Config{}
	err = yaml.Unmarshal(bs, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// InitService inits service
func InitService(auth *admin.Auth) {
	firAuth = auth
}

// MakeTokenSource creates token source from service account
func MakeTokenSource(serviceAccount []byte) (oauth2.TokenSource, error) {
	jwtConfig, err := google.JWTConfigFromJSON([]byte(serviceAccount),
		datastore.ScopeDatastore,
		pubsub.ScopePubSub,
		storage.ScopeFullControl,
		logging.WriteScope,
	)

	if err != nil {
		return nil, err
	}
	return jwtConfig.TokenSource(context.Background()), nil
}

// MakeServiceConnection creates new grpc client connection based on given address
func MakeServiceConnection(addr string, creds credentials.TransportCredentials) (conn *grpc.ClientConn) {
	var err error
	if len(addr) == 0 {
		addr = "localhost:8081"
	} else {
	}
	conn, err = grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatal(err)
	}
	return
}

func validateHeaderToken(header string) (string, error) {
	tk := strings.Split(header, " ")
	if len(tk) != 2 || strings.ToLower(tk[0]) != "bearer" {
		return "", errors.New("invalid authorization header")
	}
	claims, err := firAuth.VerifyIDToken(tk[1])
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}
