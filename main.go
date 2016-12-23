package main

import (
	"acourse/app"
	"acourse/ctrl"
	"acourse/store"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	service := echo.New()

	db := store.NewDB(store.ProjectID("acourse-d9d0a"))

	// globals middlewares
	service.Use(middleware.Recover())
	service.Use(middleware.Logger())

	// mount controllers
	app.MountHealthController(service, ctrl.NewHealthController())
	app.MountUserController(service, ctrl.NewUserController(db))

	if err := service.Start(":8080"); err != nil {
		service.Logger.Error(err)
	}
}
