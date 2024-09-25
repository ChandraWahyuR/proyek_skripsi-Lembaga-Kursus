package routes

import (
	"skripsi/config"
	"skripsi/features/admin"
	"skripsi/features/users"
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

	e.POST("/api/v1/admin/register", a.RegisterAdmin())
	e.POST("/api/v1/admin/login", a.LoginAdmin(), echojwt.WithConfig(jwtConfig))
}
