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
	"github.com/acoshift/acourse/pkg/acourse"
	"github.com/acoshift/acourse/pkg/app"
	"github.com/acoshift/acourse/pkg/ctrl/health"
	"github.com/acoshift/acourse/pkg/ctrl/render"
	"github.com/acoshift/acourse/pkg/service/assignment"
	"github.com/acoshift/acourse/pkg/service/course"
	"github.com/acoshift/acourse/pkg/service/email"
	"github.com/acoshift/acourse/pkg/service/payment"
	"github.com/acoshift/acourse/pkg/service/user"
	"github.com/acoshift/acourse/pkg/store"
	"github.com/acoshift/cors"
	"github.com/acoshift/ds"
	"github.com/acoshift/go-firebase-admin"
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
		panic(err)
	}

	firApp, err := admin.InitializeApp(admin.AppOptions{
		ProjectID:      cfg.Firebase.ProjectID,
		ServiceAccount: []byte(cfg.Firebase.ServiceAccount),
	})
	if err != nil {
		panic(err)
	}
	firAuth := firApp.Auth()

	db := store.NewDB(store.ProjectID(cfg.ProjectID), store.ServiceAccount([]byte(cfg.ServiceAccount)))

	jwtConfig, err := google.JWTConfigFromJSON([]byte(cfg.ServiceAccount),
		datastore.ScopeDatastore,
		pubsub.ScopePubSub,
		storage.ScopeFullControl,
	)
	if err != nil {
		panic(err)
	}
	tokenSource := jwtConfig.TokenSource(context.Background())

	client, err := ds.NewClient(context.Background(), cfg.ProjectID, option.WithTokenSource(tokenSource))
	if err != nil {
		panic(err)
	}

	httpServer := chain(
		app.Logger,
		app.RequestID,
		app.HSTS,
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
	)

	mux := http.NewServeMux()

	app.InitService(firAuth)

	// create service clients
	conn, err := grpc.Dial("127.0.0.1:8081", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	userServiceClient := acourse.NewUserServiceClient(conn)
	courseServiceClient := acourse.NewCourseServiceClient(conn)
	emailServiceClient := acourse.NewEmailServiceClient(conn)
	paymentServiceClient := acourse.NewPaymentServiceClient(conn)
	assignmentServiceClient := acourse.NewAssignmentServiceClient(conn)

	// register service clients to http server
	app.RegisterUserServiceClient(mux, userServiceClient)
	app.RegisterCourseServiceClient(mux, courseServiceClient)
	// app.RegisterEmailServiceClient(mux, emailService) // do not expose email service to the world right now
	app.RegisterPaymentServiceClient(mux, paymentServiceClient)
	app.RegisterAssignmentServiceClient(mux, assignmentServiceClient)

	// mount controllers
	app.MountHealthController(mux, health.New())
	app.MountRenderController(mux, render.New(db, courseServiceClient))

	// run grpc server
	go func() {
		grpcListener, err := net.Listen("tcp", ":8081")
		if err != nil {
			panic(err)
		}
		grpcServer := grpc.NewServer(grpc.UnaryInterceptor(app.UnaryInterceptors))
		acourse.RegisterUserServiceServer(grpcServer, user.New(client))
		acourse.RegisterCourseServiceServer(grpcServer, course.New(db, userServiceClient))
		acourse.RegisterPaymentServiceServer(grpcServer, payment.New(db, client, userServiceClient, firAuth, emailServiceClient))
		acourse.RegisterEmailServiceServer(grpcServer, email.New(email.Config{
			From:     cfg.Email.From,
			Server:   cfg.Email.Server,
			Port:     cfg.Email.Port,
			User:     cfg.Email.User,
			Password: cfg.Email.Password,
		}))
		acourse.RegisterAssignmentServiceServer(grpcServer, assignment.New(db, client))
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
