package main

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/acoshift/acourse/pkg/acourse"
	"github.com/acoshift/acourse/pkg/app"
	"github.com/acoshift/acourse/pkg/ctrl"
	"github.com/acoshift/acourse/pkg/service/course"
	"github.com/acoshift/acourse/pkg/service/email"
	"github.com/acoshift/acourse/pkg/service/payment"
	"github.com/acoshift/acourse/pkg/service/user"
	"github.com/acoshift/acourse/pkg/store"
	"github.com/acoshift/go-firebase-admin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"gopkg.in/gin-contrib/cors.v1"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/yaml.v2"
)

// Config type
type Config struct {
	Debug     bool   `yaml:"debug"`
	Port      string `yaml:"port"`
	ProjectID string `yaml:"projectId"`
	Email     struct {
		From     string `yaml:"from"`
		Server   string `yaml:"server"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	} `yaml:"email"`
	ServiceAccountString string `yaml:"serviceAccount"`
	ServiceAccount       []byte `yaml:"-"`
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
	cfg.ServiceAccount = []byte(cfg.ServiceAccountString)
	return &cfg, nil
}

func main() {
	grpclog.SetLogger(app.NewLogger())

	configFile := os.Getenv("CONFIG")
	if configFile == "" {
		configFile = "config.yaml"
	}

	cfg, err := LoadConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	firApp, err := admin.InitializeApp(admin.AppOptions{
		ProjectID:      cfg.ProjectID,
		ServiceAccount: cfg.ServiceAccount,
	})
	if err != nil {
		return
	}
	firAuth := firApp.Auth()

	httpServer := gin.New()

	db := store.NewDB(store.ProjectID(cfg.ProjectID), store.ServiceAccount(cfg.ServiceAccount))

	// globals middlewares
	httpServer.Use(gin.Logger())
	httpServer.Use(gin.Recovery())
	httpServer.Use(cors.New(cors.Config{
		AllowCredentials: false,
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodDelete},
		AllowOrigins:     []string{"https://acourse.io", "http://localhost:9000"},
		MaxAge:           12 * time.Hour,
	}))

	if err := app.InitService(httpServer, firAuth); err != nil {
		log.Fatal(err)
	}

	// create service clients
	conn, err := grpc.Dial("127.0.0.1:8081", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
		return
	}
	userServiceClient := acourse.NewUserServiceClient(conn)
	courseServiceClient := acourse.NewCourseServiceClient(conn)
	emailServiceClient := acourse.NewEmailServiceClient(conn)
	paymentServiceClient := acourse.NewPaymentServiceClient(conn)

	// register service clients to http server
	app.RegisterUserServiceClient(httpServer, userServiceClient)
	app.RegisterCourseServiceClient(httpServer, courseServiceClient)
	// app.RegisterEmailServiceClient(httpServer, emailService) // do not expose email service to the world right now
	app.RegisterPaymentServiceClient(httpServer, paymentServiceClient)

	// mount controllers
	app.MountHealthController(httpServer, ctrl.NewHealth())
	app.MountRenderController(httpServer, ctrl.NewRenderController(db, courseServiceClient))

	// run grpc server
	go func() {
		grpcListener, err := net.Listen("tcp", ":8081")
		if err != nil {
			log.Fatal(err)
		}
		grpcServer := grpc.NewServer(grpc.UnaryInterceptor(app.UnaryInterceptors))
		acourse.RegisterUserServiceServer(grpcServer, user.New(db))
		acourse.RegisterCourseServiceServer(grpcServer, course.New(db))
		acourse.RegisterPaymentServiceServer(grpcServer, payment.New(db, firAuth, emailServiceClient))
		acourse.RegisterEmailServiceServer(grpcServer, email.New(email.Config{
			From:     cfg.Email.From,
			Server:   cfg.Email.Server,
			Port:     cfg.Email.Port,
			User:     cfg.Email.User,
			Password: cfg.Email.Password,
		}))
		if err = grpcServer.Serve(grpcListener); err != nil {
			log.Fatal(err)
		}
	}()

	hostPort := net.JoinHostPort("0.0.0.0", cfg.Port)
	log.Printf("Listening on %s", hostPort)

	if !cfg.Debug {
		go payment.StartNotification(db, emailServiceClient)
	}

	if err := httpServer.Run(hostPort); err != nil {
		log.Fatal(err)
	}
}
