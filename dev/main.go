package main

import (
	"acourse/app"
	"acourse/ctrl"
	"acourse/store"
	"log"

	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
)

const projectID = "acourse-d9d0a"

func main() {
	gin.SetMode(gin.DebugMode)
	service := gin.New()

	db := store.NewDB(store.ProjectID(projectID), store.ServiceAccount("private/service_account.json"))

	// globals middlewares
	service.Use(gin.Logger())
	service.Use(gin.Recovery())
	service.Use(cors.New(cors.Config{
		AllowCredentials: false,
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE"},
		AllowedOrigins:   []string{"http://localhost:9000", "http://localhost:8080"},
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
