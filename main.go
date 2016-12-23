package main

import (
	"acourse/app"
	"acourse/ctrl"
	"acourse/store"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const projectID = "acourse-d9d0a"

func main() {
	service := echo.New()

	db := store.NewDB(store.ProjectID(projectID))

	// globals middlewares
	service.Use(middleware.Recover())
	service.Use(middleware.Logger())

	app.InitService(service, projectID)

	// mount controllers
	app.MountHealthController(service, ctrl.NewHealthController())
	app.MountUserController(service, ctrl.NewUserController(db))
	app.MountCourseController(service, ctrl.NewCourseController(db))

	if err := service.Start(":8080"); err != nil {
		service.Logger.Error(err)
	}
}
