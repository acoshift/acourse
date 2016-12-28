package main

import (
	"acourse/app"
	"acourse/ctrl"
	"acourse/store"
	"log"

	"gopkg.in/gin-contrib/cors.v1"
	"gopkg.in/gin-gonic/gin.v1"
)

const projectID = "acourse-d9d0a"

func main() {
	gin.SetMode(gin.ReleaseMode)
	service := gin.New()

	db := store.NewDB(store.ProjectID(projectID), store.ServiceAccount("private/service_account.json"))

	// globals middlewares
	service.Use(gin.Logger())
	service.Use(gin.Recovery())
	service.Use(cors.New(cors.Config{
		AllowCredentials: false,
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE"},
		AllowOrigins:     []string{"https://acourse.io"},
		MaxAge:           3600,
	}))

	if err := app.InitService(service, projectID); err != nil {
		log.Fatal(err)
	}

	// mount controllers
	app.MountHealthController(service.Group("/_ah"), ctrl.NewHealthController())
	app.MountUserController(service.Group("/api/user"), ctrl.NewUserController(db))
	app.MountCourseController(service.Group("/api/course"), ctrl.NewCourseController(db))
	app.MountPaymentController(service.Group("/api/payment"), ctrl.NewPaymentController(db))
	app.MountRenderController(service, ctrl.NewRenderController(db))

	if err := service.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
