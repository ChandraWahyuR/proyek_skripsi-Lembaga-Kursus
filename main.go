package main

import (
	"skripsi/config"
	"skripsi/helper"
	"skripsi/routes"
	"skripsi/utils/database"
	"skripsi/utils/database/seeders"

	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	UsersData "skripsi/features/users/data"
	UsersHandler "skripsi/features/users/handler"
	UsersService "skripsi/features/users/service"

	AdminData "skripsi/features/admin/data"
	AdminHandler "skripsi/features/admin/handler"
	AdminService "skripsi/features/admin/service"
)

func main() {
	cfg := config.InitConfig()
	db, err := database.InitDB(*cfg)
	if err != nil {
		return
	}

	err = database.Migrate(db)
	if err != nil {
		return
	}
	seeder := seeders.NewSeeder(db)
	seeder.Seed()
	e := echo.New()

	// JWT and Mailer
	jwt := helper.NewJWT(cfg.JWT_Secret)
	mailer := helper.NewMailer(cfg.SMTP)

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	usersData := UsersData.New(db)
	usersService := UsersService.New(usersData, jwt, mailer)
	usersHandler := UsersHandler.New(usersService, jwt)

	adminData := AdminData.New(db)
	adminService := AdminService.New(adminData, jwt)
	adminHandler := AdminHandler.New(adminService, jwt)

	routes.RouteUser(e, usersHandler, *cfg)
	routes.RouteAdmin(e, adminHandler, *cfg)

	// Swagger
	e.Static("/", "static")
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/swagger.yaml", func(c echo.Context) error {
		return c.File("docs/user.yaml")
	})

	e.Logger.Fatal(e.Start(":8080"))
}
