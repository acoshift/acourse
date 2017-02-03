package main

import (
	"context"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"cloud.google.com/go/trace"
	"github.com/acoshift/acourse/pkg/acourse"
	"github.com/acoshift/acourse/pkg/app"
	"github.com/acoshift/acourse/pkg/ctrl/health"
	"github.com/acoshift/acourse/pkg/ctrl/render"
	"github.com/acoshift/acourse/pkg/service/assignment"
	"github.com/acoshift/acourse/pkg/service/course"
	"github.com/acoshift/acourse/pkg/service/email"
	"github.com/acoshift/acourse/pkg/service/payment"
	"github.com/acoshift/acourse/pkg/service/user"
	"github.com/acoshift/cors"
	"github.com/acoshift/ds"
	"github.com/acoshift/go-firebase-admin"
	"github.com/acoshift/gzip"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"gopkg.in/yaml.v2"
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
	Email     struct {
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

func chain(hs ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		for i := len(hs); i > 0; i-- {
			h = hs[i-1](h)
		}
		return h
	}
}

func main() {
	grpclog.SetLogger(app.NewNoFatalLogger())

	configFile := os.Getenv("CONFIG")
	if configFile == "" {
		configFile = "config.yaml"
	}

	cfg, err := LoadConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	firApp, err := admin.InitializeApp(admin.AppOptions{
		ProjectID:      cfg.Firebase.ProjectID,
		ServiceAccount: []byte(cfg.Firebase.ServiceAccount),
	})
	if err != nil {
		log.Fatal(err)
	}
	firAuth, err := firApp.Auth()
	if err != nil {
		log.Fatal(err)
	}

	jwtConfig, err := google.JWTConfigFromJSON([]byte(cfg.ServiceAccount),
		datastore.ScopeDatastore,
		pubsub.ScopePubSub,
		storage.ScopeFullControl,
	)
	if err != nil {
		log.Fatal(err)
	}
	tokenSource := jwtConfig.TokenSource(context.Background())

	client, err := ds.NewClient(context.Background(), cfg.ProjectID, option.WithTokenSource(tokenSource))
	if err != nil {
		log.Fatal(err)
	}

	traceClient, err := trace.NewClient(context.Background(), cfg.ProjectID, option.WithTokenSource(tokenSource))
	if err != nil {
		log.Fatal(err)
	}

	middlewares := []func(http.Handler) http.Handler{
		app.Trace(traceClient),
		app.Logger,
		app.RequestID,
	}
	if len(cfg.TLSPort) > 0 {
		middlewares = append(middlewares, app.HSTS)
	}
	middlewares = append(middlewares,
		app.Recovery,
		cors.New(cors.Config{
			AllowCredentials: false,
			AllowOrigins: []string{
				"https://acourse.io",
				"http://localhost:8080",
				"https://localhost:8080",
				"http://localhost:9000",
			},
			AllowMethods: []string{
				http.MethodGet,
				http.MethodPost,
			},
			AllowHeaders: []string{
				"Content-Type",
				"Authorization",
			},
			MaxAge: 12 * time.Hour,
		}),
		gzip.New(gzip.Config{Level: gzip.DefaultCompression}),
	)

	httpServer := chain(middlewares...)

	mux := http.NewServeMux()

	app.InitService(firAuth)

	// create service clients
	conn, err := grpc.Dial("127.0.0.1:8081", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	userServiceClient := acourse.NewUserServiceClient(conn)
	courseServiceClient := acourse.NewCourseServiceClient(conn)
	emailServiceClient := acourse.NewEmailServiceClient(conn)
	paymentServiceClient := acourse.NewPaymentServiceClient(conn)
	assignmentServiceClient := acourse.NewAssignmentServiceClient(conn)

	// register service clients to http server
	app.RegisterUserServiceClient(mux, userServiceClient)
	app.RegisterCourseServiceClient(mux, courseServiceClient)
	app.RegisterPaymentServiceClient(mux, paymentServiceClient)
	app.RegisterAssignmentServiceClient(mux, assignmentServiceClient)

	// mount controllers
	app.MountHealthController(mux, health.New())
	app.MountRenderController(mux, render.New(courseServiceClient))

	// run grpc server
	go func() {
		grpcListener, err := net.Listen("tcp", ":8081")
		if err != nil {
			log.Fatal(err)
		}
		grpcServer := grpc.NewServer(grpc.UnaryInterceptor(app.UnaryInterceptors))
		acourse.RegisterUserServiceServer(grpcServer, user.New(client))
		acourse.RegisterCourseServiceServer(grpcServer, course.New(client, userServiceClient, paymentServiceClient))
		acourse.RegisterPaymentServiceServer(grpcServer, payment.New(client, userServiceClient, courseServiceClient, firAuth, emailServiceClient))
		acourse.RegisterEmailServiceServer(grpcServer, email.New(email.Config{
			From:     cfg.Email.From,
			Server:   cfg.Email.Server,
			Port:     cfg.Email.Port,
			User:     cfg.Email.User,
			Password: cfg.Email.Password,
		}))
		acourse.RegisterAssignmentServiceServer(grpcServer, assignment.New(client, courseServiceClient))
		log.Fatal(grpcServer.Serve(grpcListener))
	}()

	if !cfg.Debug {
		go payment.StartNotification(client, emailServiceClient)
	}

	serverHandler := httpServer(mux)
	addr := net.JoinHostPort("0.0.0.0", cfg.Port)

	if cfg.TLSPort != "" {
		tlsAddr := net.JoinHostPort("0.0.0.0", cfg.TLSPort)
		go func() {
			log.Printf("Listening Redirect on %s", addr)
			log.Fatal(http.ListenAndServe(addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, "https://"+cfg.Domain, http.StatusMovedPermanently)
			})))
		}()
		log.Printf("Listening TLS on %s", tlsAddr)
		log.Fatal(http.ListenAndServeTLS(tlsAddr, cfg.TLSCert, cfg.TLSKey, serverHandler))
	} else {
		log.Printf("Listening on %s", addr)
		log.Fatal(http.ListenAndServe(addr, serverHandler))
	}
}
