package admin

import (
	"io"
	"time"

	"github.com/labstack/echo/v4"
)

type Admin struct {
	ID              string
	Username        string
	Email           string
	Password        string
	ConfirmPassword string
}

type Login struct {
	Username string
	Password string
	Token    string
}

type AdminHandlerInterface interface {
	RegisterAdmin() echo.HandlerFunc
	LoginAdmin() echo.HandlerFunc
	DownloadLaporanPembelian() echo.HandlerFunc
}

type AdminServiceInterface interface {
	RegisterAdmin(admin Admin) error
	LoginAdmin(admin Admin) (Login, error)
	//
	DownloadLaporanPembelian(startDate, endDate time.Time) ([]map[string]interface{}, error)
	GenerateLaporanCSV(w io.Writer, histories []map[string]interface{}, startDate, endDate time.Time) error
}

type AdminDataInterface interface {
	RegisterAdmin(admin Admin) error
	LoginAdmin(admin Admin) (Admin, error)
	IsEmailExist(email string) bool
	IsUsernameExist(username string) bool
	//
	DownloadLaporanPembelian(startDate, endDate time.Time) ([]map[string]interface{}, error)
}
