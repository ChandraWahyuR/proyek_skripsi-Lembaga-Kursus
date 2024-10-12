package transaksi

import (
	"skripsi/features/kursus"
	"skripsi/features/users"
	"skripsi/features/voucher"
	"time"

	"github.com/labstack/echo/v4"
)

type Transaksi struct {
	ID         string
	TotalHarga float64
	VoucherID  string
	KursusID   string
	Kursus     kursus.Kursus
	UserID     string
	User       users.User
	SnapURL    string
	Status     string
}

type TransaksiHistory struct {
	ID          string
	TransaksiID string
	Transaksi   Transaksi
	KursusID    string
	Kursus      kursus.Kursus
	UserID      string
	User        users.User
	VoucherID   string
	Voucher     voucher.Voucher
	TotalHarga  float64
	Status      string
	ValidUntil  time.Time
}

type UpdateTransaksiStatus struct {
	ID     string
	Status string
}

type TransaksiHandlerInterface interface {
	CreateTransaksi() echo.HandlerFunc
	// GetAllStatusTransaksi() echo.HandlerFunc
	// GetStatusTransaksiForUser() echo.HandlerFunc
	// GetStatusTransaksiByID() echo.HandlerFunc
	// Tansaksi History
	// GetAllTransaksiHistory() echo.HandlerFunc
	// GetAllTransaksiHistoryForUser() echo.HandlerFunc
	// GetTransaksiHistoryByID() echo.HandlerFunc
}

type TransaksiDataInterface interface {
	CreateTransaksi(Transaksi) (Transaksi, error)
	GetTotalTransaksiWithDiscount(total float64, voucherId string) (float64, error)

	GetAllStatusTransaksi() ([]Transaksi, error)
	GetStatusTransaksiForUser(userID string) ([]Transaksi, error)
	GetStatusTransaksiByID(id string) (Transaksi, error)
	// Tansaksi History
	GetAllTransaksiHistory() ([]TransaksiHistory, error)
	GetAllTransaksiHistoryForUser(userID string) ([]TransaksiHistory, error)
	GetTransaksiHistoryByID(id string) (TransaksiHistory, error)
	// Pagination
	GetAllTransaksiPagination(page, limit int) ([]Transaksi, int, error)
	GetAllHistoryTransaksiPagination(page, limit int) ([]TransaksiHistory, int, error)

	GetByIDVoucher(id string) (voucher.Voucher, error)
	GetKursusByID(kursusID string) (kursus.Kursus, error)
	GetUserByID(userID string) (users.User, error)
}

type TransaksiServiceInterface interface {
	CreateTransaksi(Transaksi) (Transaksi, error)
	// GetAllStatusTransaksi() ([]Transaksi, error)
	// GetTotalTransaksiWithDiscount(total float64, voucherId string) (float64, error)
	// GetStatusTransaksiForUser(userID string) ([]Transaksi, error)
	// GetStatusTransaksiByID(id string) (Transaksi, error)
	// // Tansaksi History
	// GetAllTransaksiHistory() ([]TransaksiHistory, error)
	// GetAllTransaksiHistoryForUser(userID string) ([]TransaksiHistory, error)
	// GetTransaksiHistoryByID(id string) (TransaksiHistory, error)
	// // Pagination
	// GetAllTransaksiPagination(page, limit int) ([]Transaksi, int, error)
	// GetAllHistoryTransaksiPagination(page, limit int) ([]TransaksiHistory, int, error)

}
