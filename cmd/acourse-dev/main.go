package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/acoshift/acourse/pkg/app"
	"github.com/acoshift/acourse/pkg/ctrl"
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
	service := gin.New()

	db := store.NewDB(store.ProjectID(cfg.ProjectID), store.ServiceAccount(serviceAccount))

	// globals middlewares
	service.Use(gin.Logger())
	service.Use(gin.Recovery())
	service.Use(cors.New(cors.Config{
		AllowCredentials: false,
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodDelete},
		AllowAllOrigins:  true,
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
	emailService := email.New(email.Config{
		From:     cfg.Email.From,
		Server:   cfg.Email.Server,
		Port:     cfg.Email.Port,
		User:     cfg.Email.User,
		Password: cfg.Email.Password,
	})

	// mount controllers
	courseCtrl := ctrl.NewCourseController(db)
	app.MountHealthService(service, health.New())
	app.MountUserService(service, user.New(db))
	app.MountPaymentService(service, payment.New(db, firAuth, emailService))
	app.MountCourseController(service.Group("/api/course"), courseCtrl)
	app.MountRenderController(service, ctrl.NewRenderController(db, courseCtrl))

	if err := service.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
