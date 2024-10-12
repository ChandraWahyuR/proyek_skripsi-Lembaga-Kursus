package routes

import (
	"skripsi/config"
	"skripsi/features/admin"
	"skripsi/features/instruktur"
	"skripsi/features/kategori"
	"skripsi/features/kursus"
	"skripsi/features/transaksi"
	"skripsi/features/users"
	"skripsi/features/voucher"
	"skripsi/helper"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func RouteUser(e *echo.Echo, u users.UserHandlerInterface, cfg config.Config) {
	jwtConfig := echojwt.Config{
		SigningKey:   []byte(cfg.JWT_Secret),
		ErrorHandler: helper.JWTErrorHandler,
	}

	e.POST("/api/v1/register", u.Register())
	e.POST("/api/v1/login", u.Login())

	e.POST("/api/v1/forgot", u.ForgotPassword())
	e.POST("/api/v1/otp", u.VerifyOTP(), echojwt.WithConfig(jwtConfig))
	e.POST("/api/v1/reset", u.ResetPassword(), echojwt.WithConfig(jwtConfig))
}

func RouteAdmin(e *echo.Echo, a admin.AdminHandlerInterface, cfg config.Config) {
	jwtConfig := echojwt.Config{
		SigningKey:   []byte(cfg.JWT_Secret),
		ErrorHandler: helper.JWTErrorHandler,
	}

	e.POST("/api/v1/admin/register", a.RegisterAdmin(), echojwt.WithConfig(jwtConfig))
	e.POST("/api/v1/admin/login", a.LoginAdmin())
}

func RouteInstruktor(e *echo.Echo, i instruktur.InstrukturHandlerInterface, cfg config.Config) {
	jwtConfig := echojwt.Config{
		SigningKey:   []byte(cfg.JWT_Secret),
		ErrorHandler: helper.JWTErrorHandler,
	}

	e.POST("/api/v1/instruktur", i.PostInstruktur(), echojwt.WithConfig(jwtConfig))
	e.GET("/api/v1/instruktur", i.GetAllInstruktur(), echojwt.WithConfig(jwtConfig))
	e.GET("/api/v1/instruktur/:id", i.GetAllInstrukturByID(), echojwt.WithConfig(jwtConfig))
	e.PUT("/api/v1/instruktur/:id", i.UpdateInstruktur(), echojwt.WithConfig(jwtConfig))
	e.DELETE("/api/v1/instruktur/:id", i.DeleteInstruktur(), echojwt.WithConfig(jwtConfig))
	e.GET("/api/v1/instruktur/search", i.GetInstruktorByName(), echojwt.WithConfig(jwtConfig))
}

func RouteKategori(e *echo.Echo, k kategori.KategoriHandlerInterface, cfg config.Config) {
	jwtConfig := echojwt.Config{
		SigningKey:   []byte(cfg.JWT_Secret),
		ErrorHandler: helper.JWTErrorHandler,
	}

	e.GET("/api/v1/kategori", k.GetAllKategori(), echojwt.WithConfig(jwtConfig))
	e.GET("/api/v1/kategori/:id", k.GetKategoriById(), echojwt.WithConfig(jwtConfig))
	e.POST("/api/v1/kategori", k.CreateKategori(), echojwt.WithConfig(jwtConfig))
	e.PUT("/api/v1/kategori/:id", k.UpdateKategori(), echojwt.WithConfig(jwtConfig))
	e.DELETE("/api/v1/kategori/:id", k.DeleteKategori(), echojwt.WithConfig(jwtConfig))
}

func RouteKursus(e *echo.Echo, kr kursus.KursusHandlerInterface, cfg config.Config) {
	jwtConfig := echojwt.Config{
		SigningKey:   []byte(cfg.JWT_Secret),
		ErrorHandler: helper.JWTErrorHandler,
	}

	e.GET("/api/v1/kursus", kr.GetAllKursus(), echojwt.WithConfig(jwtConfig))
	e.GET("/api/v1/kursus/:id", kr.GetAllKursusById(), echojwt.WithConfig(jwtConfig))
	e.PUT("/api/v1/kursus/:id", kr.UpdateKursus(), echojwt.WithConfig(jwtConfig))
	e.DELETE("/api/v1/kursus/:id", kr.DeleteKursus(), echojwt.WithConfig(jwtConfig))
	e.POST("/api/v1/kursus", kr.AddKursus(), echojwt.WithConfig(jwtConfig))
	e.GET("/api/v1/kursus/search", kr.GetAllKursusByName(), echojwt.WithConfig(jwtConfig))
}

func RouteVoucher(e *echo.Echo, vc voucher.VoucherHandlerInterface, cfg config.Config) {
	jwtConfig := echojwt.Config{
		SigningKey:   []byte(cfg.JWT_Secret),
		ErrorHandler: helper.JWTErrorHandler,
	}

	e.GET("/api/v1/voucher", vc.GetAllVoucher(), echojwt.WithConfig(jwtConfig))
	e.GET("/api/v1/voucher/:id", vc.GetByIDVoucher(), echojwt.WithConfig(jwtConfig))
	e.POST("/api/v1/voucher", vc.CreateVoucher(), echojwt.WithConfig(jwtConfig))
	e.PUT("/api/v1/voucher/:id", vc.UpdateVoucher(), echojwt.WithConfig(jwtConfig))
	e.DELETE("/api/v1/voucher/:id", vc.DeleteVoucher(), echojwt.WithConfig(jwtConfig))
}

func RouteTransaksi(e *echo.Echo, tr transaksi.TransaksiHandlerInterface, cfg config.Config) {
	jwtConfig := echojwt.Config{
		SigningKey:   []byte(cfg.JWT_Secret),
		ErrorHandler: helper.JWTErrorHandler,
	}

	e.POST("/api/v1/transaksi", tr.CreateTransaksi(), echojwt.WithConfig(jwtConfig))
}
