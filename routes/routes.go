package routes

import (
	"skripsi/config"
	"skripsi/features/admin"
	"skripsi/features/gmaps"
	"skripsi/features/instruktur"
	jadwal "skripsi/features/jadwal_mengajar"
	"skripsi/features/kategori"
	"skripsi/features/kursus"
	"skripsi/features/notification/sse"
	"skripsi/features/transaksi"
	"skripsi/features/users"
	"skripsi/features/voucher"
	"skripsi/features/webhook"
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

	e.GET("/ping", func(c echo.Context) error {
		return c.String(200, "pong")
	})
	// Edit
	e.GET("/api/v1/profile", u.GetUserByUser(), echojwt.WithConfig(jwtConfig))
	e.PUT("/api/v1/profile", u.UpdateUser(), echojwt.WithConfig(jwtConfig))
	// Admin
	e.GET("/api/v1/admin/users", u.GetAllUser(), echojwt.WithConfig(jwtConfig))
	e.GET("/api/v1/admin/users/:id", u.GetUserByID(), echojwt.WithConfig(jwtConfig))
	e.GET("/api/v1/admin/users/search", u.SearchUserByUsernameEmail(), echojwt.WithConfig(jwtConfig))
	e.POST("/api/v1/logout", u.Logout(), echojwt.WithConfig(jwtConfig))
	e.GET("/verify", u.VerifyAccount()) // handler untuk memverifikasi token
}

func RouteAdmin(e *echo.Echo, a admin.AdminHandlerInterface, cfg config.Config) {
	jwtConfig := echojwt.Config{
		SigningKey:   []byte(cfg.JWT_Secret),
		ErrorHandler: helper.JWTErrorHandler,
	}

	e.POST("/api/v1/admin/register", a.RegisterAdmin(), echojwt.WithConfig(jwtConfig))
	e.POST("/api/v1/admin/login", a.LoginAdmin())
	e.POST("/api/v1/admin/laporan-pembelian", a.DownloadLaporanPembelian())
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
	e.GET("/api/v1/GetAllValidVoucher", vc.GetAllValidVoucher(), echojwt.WithConfig(jwtConfig))
}

func RouteTransaksi(e *echo.Echo, tr transaksi.TransaksiHandlerInterface, cfg config.Config) {
	jwtConfig := echojwt.Config{
		SigningKey:   []byte(cfg.JWT_Secret),
		ErrorHandler: helper.JWTErrorHandler,
	}

	e.POST("/api/v1/transaksi", tr.CreateTransaksi(), echojwt.WithConfig(jwtConfig))
	e.GET("/api/v1/admin/list-transaksi", tr.GetAllStatusTransaksi(), echojwt.WithConfig(jwtConfig))
	e.GET("/api/v1/admin/list-transaksi/:id", tr.GetStatusTransaksiByID(), echojwt.WithConfig(jwtConfig))
	e.GET("/api/v1/list-transaksi", tr.GetStatusTransaksiForUser(), echojwt.WithConfig(jwtConfig))
	e.GET("/api/v1/history-transaksi", tr.GetAllTransaksiHistoryForUser(), echojwt.WithConfig(jwtConfig))

	e.GET("/api/v1/list-transaksi/:id", tr.GetStatusTransaksiForUserByID(), echojwt.WithConfig(jwtConfig))
	e.GET("/api/v1/history-transaksi/:id", tr.GetAllTransaksiHistoryForUserByID(), echojwt.WithConfig(jwtConfig))
	//
	e.GET("/api/v1/response-transaksi", tr.GetResponseTransaksi())

	e.GET("/api/v1/admin/history-transaksi", tr.GetAllTransaksiHistory(), echojwt.WithConfig(jwtConfig))
	e.GET("/api/v1/admin/history-transaksi/:id", tr.GetTransaksiHistoryByID(), echojwt.WithConfig(jwtConfig))
	e.GET("/api/v1/admin/users-active", tr.GetActiveUsersFromTransaksiHistory(), echojwt.WithConfig(jwtConfig))
	e.GET("/api/v1/admin/users-new", tr.GetNewUsers(), echojwt.WithConfig(jwtConfig))
}

func RouteJadwal(e *echo.Echo, j jadwal.MengajarHandlerInterface, cfg config.Config) {
	jwtConfig := echojwt.Config{
		SigningKey:   []byte(cfg.JWT_Secret),
		ErrorHandler: helper.JWTErrorHandler,
	}
	e.POST("/api/v1/admin/jadwal", j.CreateJadwalMengajar(), echojwt.WithConfig(jwtConfig))
	e.PUT("/api/v1/admin/jadwal/:id", j.EditJadwalMengajar(), echojwt.WithConfig(jwtConfig))
	e.GET("/api/v1/admin/jadwal", j.GetJadwalMengajar(), echojwt.WithConfig(jwtConfig))
	e.GET("/api/v1/jadwal/:id", j.GetJadwalMengajarByID(), echojwt.WithConfig(jwtConfig))
	e.GET("/api/v1/jadwal", j.GetJadwalMengajarForUser(), echojwt.WithConfig(jwtConfig))
}

func RouteWebhook(e *echo.Echo, w webhook.MidtransNotificationHandler, cfg config.Config) {
	e.POST("/notifikasi-midtrans", w.HandleNotification())
}

func RouteGmaps(e *echo.Echo, g gmaps.GmapsHandlerInterface, cfg config.Config) {
	e.GET("/api/v1/maps/directions", g.GetDirections())
}

func RouteSSE(e *echo.Echo, jwtHelper helper.JWTInterface) {
	sseHandler := sse.NewSSEHandler(jwtHelper)
	e.GET("/api/v1/events", func(c echo.Context) error {
		sseHandler.SseHandler(c.Response(), c.Request())
		return nil
	})
}
