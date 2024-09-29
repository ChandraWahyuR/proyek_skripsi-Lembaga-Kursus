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

	InstrukturData "skripsi/features/instruktur/data"
	InstrukturHandler "skripsi/features/instruktur/handler"
	InstrukturService "skripsi/features/instruktur/service"

	KategoriData "skripsi/features/kategori/data"
	KategoriHandler "skripsi/features/kategori/handler"
	KategoriService "skripsi/features/kategori/service"
)

// Ini logout kaya forgot juga diredis aja
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
	helper.InitGCP()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// Feature
	usersData := UsersData.New(db)
	usersService := UsersService.New(usersData, jwt, mailer)
	usersHandler := UsersHandler.New(usersService, jwt)

	adminData := AdminData.New(db)
	adminService := AdminService.New(adminData, jwt)
	adminHandler := AdminHandler.New(adminService, jwt)

	instrukturData := InstrukturData.New(db)
	instrukturService := InstrukturService.New(instrukturData, jwt)
	instrukturHandler := InstrukturHandler.New(instrukturService, jwt)

	kategoriData := KategoriData.New(db)
	kategoriService := KategoriService.New(kategoriData, jwt)
	KategoriHandler := KategoriHandler.New(kategoriService, jwt)

	routes.RouteUser(e, usersHandler, *cfg)
	routes.RouteAdmin(e, adminHandler, *cfg)
	routes.RouteInstruktor(e, instrukturHandler, *cfg)
	routes.RouteKategori(e, KategoriHandler, *cfg)

	// Swagger
	e.Static("/", "static")
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/swagger.yaml", func(c echo.Context) error {
		return c.File("docs/user.yaml")
	})

	e.Logger.Fatal(e.Start(":8080"))
}
