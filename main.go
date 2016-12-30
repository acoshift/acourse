package main

import (
	"acourse/app"
	"acourse/ctrl"
	"acourse/store"
	"io/ioutil"
	"log"
	"time"

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

	firApp, err := admin.InitializeApp(admin.ProjectID(cfg.ProjectID), admin.ServiceAccount(serviceAccount))
	if err != nil {
		return
	}
	firAuth := firApp.Auth()

	gin.SetMode(gin.ReleaseMode)
	service := gin.New()

	db := store.NewDB(store.ProjectID(cfg.ProjectID), store.ServiceAccount("private/service_account.json"))

	// globals middlewares
	service.Use(gin.Logger())
	service.Use(gin.Recovery())
	service.Use(cors.New(cors.Config{
		AllowCredentials: false,
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE"},
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
	app.MountHealthController(service.Group("/_ah"), ctrl.NewHealthController())
	app.MountUserController(service.Group("/api/user"), ctrl.NewUserController(db))
	app.MountCourseController(service.Group("/api/course"), ctrl.NewCourseController(db))
	app.MountPaymentController(service.Group("/api/payment"), ctrl.NewPaymentController(db, firAuth))
	app.MountRenderController(service, ctrl.NewRenderController(db))

	if err := service.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
