package main

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"

	"github.com/acoshift/acourse/pkg/acourse"
	"github.com/acoshift/acourse/pkg/app"
	"github.com/acoshift/acourse/pkg/service/course"
	"github.com/acoshift/acourse/pkg/service/email"
	"github.com/acoshift/acourse/pkg/service/health"
	"github.com/acoshift/acourse/pkg/service/payment"
	"github.com/acoshift/acourse/pkg/service/user"
	"github.com/acoshift/acourse/pkg/store"
	"github.com/acoshift/go-firebase-admin"
	"gopkg.in/gin-contrib/cors.v1"
	"gopkg.in/gin-gonic/gin.v1"
)

func main() {
	cfg, err := app.LoadConfig("private/config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	serviceAccount, err := ioutil.ReadFile("private/service_account.json")
	if err != nil {
		log.Fatal(err)
	}

	firApp, err := admin.InitializeApp(admin.AppOptions{
		ProjectID:      cfg.ProjectID,
		ServiceAccount: serviceAccount,
	})
	if err != nil {
		return
	}
	firAuth := firApp.Auth()

	gin.SetMode(gin.DebugMode)
	httpServer := gin.New()

	db := store.NewDB(store.ProjectID(cfg.ProjectID), store.ServiceAccount(serviceAccount))

	// globals middlewares
	httpServer.Use(gin.Logger())
	httpServer.Use(gin.Recovery())
	httpServer.Use(cors.New(cors.Config{
		AllowCredentials: false,
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodDelete},
		AllowAllOrigins:  true,
	}))

	if err := app.InitService(httpServer, firAuth); err != nil {
		log.Fatal(err)
	}

	emailService := email.New(email.Config{
		From:     cfg.Email.From,
		Server:   cfg.Email.Server,
		Port:     cfg.Email.Port,
		User:     cfg.Email.User,
		Password: cfg.Email.Password,
	})

	// run grpc server
	go func() {
		grpcListener, err := net.Listen("tcp", ":8081")
		if err != nil {
			log.Fatal(err)
		}
		grpcServer := grpc.NewServer()
		acourse.RegisterUserServiceServer(grpcServer, user.New(db))
		grpcServer.Serve(grpcListener)
	}()

	// mount services
	// courseCtrl := ctrl.NewCourseController(db)
	app.RegisterHealthService(httpServer, health.New())
	app.RegisterUserServiceServer(httpServer, "127.0.0.1:8081")
	app.RegisterCourseService(httpServer, course.New(db))
	app.RegisterPaymentService(httpServer, payment.New(db, firAuth, emailService))
	// app.MountCourseController(service.Group("/api/course"), courseCtrl)
	// app.MountRenderController(service, ctrl.NewRenderController(db, courseCtrl))

	if err := httpServer.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
