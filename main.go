package main

import (
	"fmt"
	"os"
	"skripsi/config"
	"skripsi/helper"
	"skripsi/routes"
	"skripsi/utils"
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

	KursusData "skripsi/features/kursus/data"
	KursusHandler "skripsi/features/kursus/handler"
	KursusService "skripsi/features/kursus/service"

	VoucherData "skripsi/features/voucher/data"
	VoucherHandler "skripsi/features/voucher/handler"
	VoucherService "skripsi/features/voucher/service"

	TransaksiData "skripsi/features/transaksi/data"
	TransaksiHandler "skripsi/features/transaksi/handler"
	TransaksiService "skripsi/features/transaksi/service"

	JadwalData "skripsi/features/jadwal_mengajar/data"
	JadwalHandler "skripsi/features/jadwal_mengajar/handler"
	JadwalService "skripsi/features/jadwal_mengajar/service"

	WebhookData "skripsi/features/webhook/data"
	WebhookHandler "skripsi/features/webhook/handler"
	WebhookService "skripsi/features/webhook/service"

	GmapsData "skripsi/features/gmaps/data"
	GmapsHandler "skripsi/features/gmaps/handler"
	GmapsService "skripsi/features/gmaps/service"
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

	// Redis
	// redisClient := helper.InitRedis(cfg)
	// redisHelper := helper.NewRedisHelper(redisClient)
	// jwt := helper.NewJWT(cfg.JWT_Secret, redisHelper)
	// JWT and Mailer
	jwt := helper.NewJWT(cfg.JWT_Secret)
	mailer := helper.NewMailer(cfg.SMTP)
	helper.InitGCP()
	midtransClient := utils.NewMidtransClient(cfg.Midtrans.ServerKey, cfg.Midtrans.ClientKey)

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
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

	kursusData := KursusData.New(db)
	kursusService := KursusService.New(kursusData, jwt)
	kursusHandler := KursusHandler.New(kursusService, jwt)

	voucherData := VoucherData.New(db)
	voucherService := VoucherService.New(voucherData, jwt)
	voucherHandler := VoucherHandler.New(voucherService, jwt)

	transaksiData := TransaksiData.New(db)
	transaksiService := TransaksiService.New(transaksiData, jwt, midtransClient)
	transaksiHandler := TransaksiHandler.New(transaksiService, jwt)

	jadwalData := JadwalData.New(db)
	jadwalService := JadwalService.NewServiceCatatan(jadwalData, jwt)
	jadwalHandler := JadwalHandler.New(jadwalService, jwt)

	webhookData := WebhookData.New(db)
	webhookService := WebhookService.New(webhookData)
	webhookHandler := WebhookHandler.New(webhookService)

	gmapsData := GmapsData.New(cfg.Gmaps.GOOGLE_MAPS_API_KEY)
	gmapsService := GmapsService.New(gmapsData)
	gmapsHandler := GmapsHandler.New(gmapsService)

	routes.RouteUser(e, usersHandler, *cfg)
	routes.RouteAdmin(e, adminHandler, *cfg)
	routes.RouteInstruktor(e, instrukturHandler, *cfg)
	routes.RouteKategori(e, KategoriHandler, *cfg)
	routes.RouteKursus(e, kursusHandler, *cfg)
	routes.RouteVoucher(e, voucherHandler, *cfg)
	routes.RouteTransaksi(e, transaksiHandler, *cfg)
	routes.RouteJadwal(e, jadwalHandler, *cfg)
	routes.RouteWebhook(e, webhookHandler, *cfg)
	routes.RouteGmaps(e, gmapsHandler, *cfg)
	routes.RouteSSE(e, jwt)
	// Cron job
	utils.StartCronJob(transaksiService)
	utils.CronJobMinute()

	// Redirect
	// http://localhost:8080/halaman/example.html
	e.Static("/assets", "assets")

	e.File("/verification-success", "assets/verifikasi_berhasil.html")
	// e.File("/verification-failed", "assets/verifikasi_gagal.html")

	// Swagger
	e.Static("/static", "static")
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/index.html", func(c echo.Context) error {
		return c.String(403, "Akses langsung ke halaman ini tidak diperbolehkan.")
	})

	e.GET("/dokumen", func(c echo.Context) error {
		return c.File("static/index.html")
	})

	e.GET("/docs/users.yaml", func(c echo.Context) error {
		return c.File("docs/users.yaml")
	})

	e.GET("/docs/admins.yaml", func(c echo.Context) error {
		return c.File("docs/admin.yaml")
	})

	// e.Logger.Fatal(e.Start(":8080"))
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
