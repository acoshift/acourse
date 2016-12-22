package main

import (
	"acourse/app"

	"acourse/store"

	"github.com/labstack/echo"
)

func main() {
	service := echo.New()

	db := store.NewDB(store.ProjectID("acourse-d9d0a"))

	app.MountUserController(service, &UserController{db})

	if err := service.Start(":8080"); err != nil {
		service.Logger.Error(err)
	}
}
