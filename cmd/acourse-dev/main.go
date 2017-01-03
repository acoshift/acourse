package main

import (
	"io/ioutil"
	"log"

	"github.com/acoshift/acourse/pkg/app"
	"github.com/acoshift/acourse/pkg/ctrl"
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

	firApp, err := admin.InitializeApp(admin.ProjectID(cfg.ProjectID), admin.ServiceAccount(serviceAccount))
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
		AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE"},
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

	// mount controllers
	courseCtrl := ctrl.NewCourseController(db)
	app.MountHealthController(service.Group("/_ah"), ctrl.NewHealthController())
	app.MountUserController(service.Group("/api/user"), ctrl.NewUserController(db))
	app.MountCourseController(service.Group("/api/course"), courseCtrl)
	app.MountPaymentController(service.Group("/api/payment"), ctrl.NewPaymentController(db, firAuth))
	app.MountRenderController(service, ctrl.NewRenderController(db, courseCtrl))

	if err := service.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
