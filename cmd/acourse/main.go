package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/acoshift/acourse/pkg/acourse"
	"github.com/acoshift/acourse/pkg/app"
	"github.com/acoshift/acourse/pkg/ctrl/health"
	"github.com/acoshift/acourse/pkg/ctrl/render"
	"github.com/acoshift/acourse/pkg/service/course"
	"github.com/acoshift/acourse/pkg/service/email"
	"github.com/acoshift/acourse/pkg/service/payment"
	"github.com/acoshift/acourse/pkg/service/user"
	"github.com/acoshift/acourse/pkg/store"
	"github.com/acoshift/cors"
	"github.com/acoshift/go-firebase-admin"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
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

type loggerWriter struct {
	http.ResponseWriter
	header int
}

func (w *loggerWriter) WriteHeader(header int) {
	w.header = header
	w.ResponseWriter.WriteHeader(header)
}

func logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		path := r.URL.Path
		tw := &loggerWriter{w, 0}
		ip := r.Header.Get("X-Real-IP")
		if ip == "" {
			ip = r.RemoteAddr
		}
		h.ServeHTTP(tw, r)
		end := time.Now()
		fmt.Printf("%v | %3d | %13v | %s | %s | %s | %s\n",
			end.Format(time.RFC3339),
			tw.header,
			end.Sub(start),
			ip,
			w.Header().Get("X-Request-ID"),
			r.Method,
			path,
		)
	})
}

func recovery(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e := recover(); e != nil {
				log.Println(e)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "%v", e)
			}
		}()
		h.ServeHTTP(w, r)
	})
}

func requestID(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := uuid.New().String()
		w.Header().Set("X-Request-ID", rid)
		h.ServeHTTP(w, r)
	})
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

	db := store.NewDB(store.ProjectID(cfg.ProjectID), store.ServiceAccount(cfg.ServiceAccount))

	httpServer := chain(
		logger,
		requestID,
		recovery,
		cors.New(cors.Config{
			AllowCredentials: false,
			AllowOrigins: []string{
				"https://acourse.io",
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

	if err := app.InitService(firAuth); err != nil {
		log.Fatal(err)
	}

	// create service clients
	conn, err := grpc.Dial("127.0.0.1:8081", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	userServiceClient := acourse.NewUserServiceClient(conn)
	courseServiceClient := acourse.NewCourseServiceClient(conn)
	emailServiceClient := acourse.NewEmailServiceClient(conn)
	paymentServiceClient := acourse.NewPaymentServiceClient(conn)

	// register service clients to http server
	app.RegisterUserServiceClient(mux, userServiceClient)
	app.RegisterCourseServiceClient(mux, courseServiceClient)
	// app.RegisterEmailServiceClient(mux, emailService) // do not expose email service to the world right now
	app.RegisterPaymentServiceClient(mux, paymentServiceClient)

	// mount controllers
	app.MountHealthController(mux, health.New())
	app.MountRenderController(mux, render.New(db, courseServiceClient))

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

	if !cfg.Debug {
		go payment.StartNotification(db, emailServiceClient)
	}

	addr := net.JoinHostPort("0.0.0.0", cfg.Port)
	log.Printf("Listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, httpServer(mux)))
}
