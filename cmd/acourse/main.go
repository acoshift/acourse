package main

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/acoshift/acourse/pkg/app"
	"github.com/acoshift/acourse/pkg/ctrl"
	"github.com/acoshift/acourse/pkg/service/health"
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

	gin.SetMode(gin.ReleaseMode)
	service := gin.New()

	db := store.NewDB(store.ProjectID(cfg.ProjectID), store.ServiceAccount(serviceAccount))

	// globals middlewares
	service.Use(gin.Logger())
	service.Use(gin.Recovery())
	service.Use(cors.New(cors.Config{
		AllowCredentials: false,
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodDelete},
		AllowOrigins:     []string{"https://acourse.io", "http://localhost:9000"},
		MaxAge:           12 * time.Hour,
	}))

	if err := app.InitService(service, firAuth); err != nil {
		log.Fatal(err)
	}

	ctrl.InitMail(ctrl.EmailConfig{
		From:     cfg.Email.From,
		Server:   cfg.Email.Server,
		Port:     cfg.Email.Port,
		User:     cfg.Email.User,
		Password: cfg.Email.Password,
	})

	ctrl.StartNotiPayment(db)

	// mount controllers
	courseCtrl := ctrl.NewCourseController(db)
	app.MountHealthService(service, health.New())
	app.MountUserService(service, user.New(db))
	app.MountCourseController(service.Group("/api/course"), courseCtrl)
	app.MountPaymentController(service.Group("/api/payment"), ctrl.NewPaymentController(db, firAuth))
	app.MountRenderController(service, ctrl.NewRenderController(db, courseCtrl))

	hostPort := net.JoinHostPort("0.0.0.0", os.Getenv("PORT"))
	log.Printf("Listening on %s", hostPort)

	if err := service.Run(hostPort); err != nil {
		log.Fatal(err)
	}
}
