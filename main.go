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
	service.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: false,
		AllowHeaders:     []string{echo.HeaderAuthorization, echo.HeaderContentType},
		AllowMethods:     []string{echo.GET, echo.POST, echo.PATCH, echo.PUT, echo.DELETE},
		AllowOrigins:     []string{"http://localhost:9000", "https://acourse.io"},
		MaxAge:           3600,
	}))

	if err := app.InitService(service, projectID); err != nil {
		service.Logger.Fatal(err)
	}

	// mount controllers
	app.MountHealthController(service, ctrl.NewHealthController())
	app.MountUserController(service, ctrl.NewUserController(db))
	app.MountCourseController(service, ctrl.NewCourseController(db))

	if err := service.Start(":8080"); err != nil {
		service.Logger.Fatal(err)
	}
}
