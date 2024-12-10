package transaksi

import (
	"skripsi/features/kursus"
	"skripsi/features/users"
	"skripsi/features/voucher"
	"time"

	"github.com/labstack/echo/v4"
)

// Belum kelar yang get
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
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type TransaksiHistory struct {
	ID          string
	TransaksiID string
	Transaksi   Transaksi
	KursusID    string
	Kursus      kursus.Kursus
	UserID      string
	User        users.User
	Status      string
	ValidUntil  time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UpdateTransaksiStatus struct {
	ID     string
	Status string
}

type UpdateHistoryStatus struct {
	ID         string
	Status     string
	ValidUntil time.Time
}

type TransaksiHandlerInterface interface {
	CreateTransaksi() echo.HandlerFunc
	GetAllStatusTransaksi() echo.HandlerFunc
	GetStatusTransaksiForUser() echo.HandlerFunc
	GetStatusTransaksiByID() echo.HandlerFunc
	GetResponseTransaksi() echo.HandlerFunc

	GetStatusTransaksiForUserByID() echo.HandlerFunc
	// Tansaksi History
	GetAllTransaksiHistory() echo.HandlerFunc
	GetAllTransaksiHistoryForUser() echo.HandlerFunc
	GetTransaksiHistoryByID() echo.HandlerFunc

	GetAllTransaksiHistoryForUserByID() echo.HandlerFunc
}

type TransaksiDataInterface interface {
	CreateTransaksi(Transaksi) (Transaksi, error)
	GetTotalTransaksiWithDiscount(total float64, voucherId string) (float64, error)

	GetAllStatusTransaksi() ([]Transaksi, error)
	GetStatusTransaksiForUser(userID string, page int, limit int) ([]Transaksi, int, error)
	GetStatusTransaksiByID(id string) (Transaksi, error)

	// Tansaksi History
	CreateTransaksiHistory(TransaksiHistory) error
	GetAllTransaksiHistory() ([]TransaksiHistory, error)
	GetAllTransaksiHistoryForUser(userID string, page, limit int) ([]TransaksiHistory, int, error)
	GetTransaksiHistoryByID(id string) (TransaksiHistory, error) // Belum kepakai
	// Pagination
	GetAllTransaksiPagination(page, limit int) ([]Transaksi, int, error)
	GetAllHistoryTransaksiPagination(page, limit int) ([]TransaksiHistory, int, error)

	GetByIDVoucher(id string) (voucher.Voucher, error)
	GetKursusByID(kursusID string) (kursus.Kursus, error)
	GetUserByID(userID string) (users.User, error)
	ValidateUserDokumentation(userId string) bool
	//
	UsedVoucher(voucher.VoucherUsed) error
	UsedVoucherCheck(userID, voucherID string) bool
	CheckVoucherExists(voucherID string) (bool, error)
}

type TransaksiServiceInterface interface {
	CreateTransaksi(Transaksi) (Transaksi, error)
	GetStatusTransaksiForUser(userID string, page int, limit int) ([]Transaksi, int, error)
	GetStatusTransaksiByID(id string) (Transaksi, error)
	UpdateTransaksiStatus(id string) (Transaksi, error)
	// // Tansaksi History
	CreateTransaksiHistory(TransaksiHistory) error
	GetAllTransaksiHistoryForUser(userID string, page, limit int) ([]TransaksiHistory, int, error)
	GetTransaksiHistoryByID(id string) (TransaksiHistory, error)
	// // Pagination
	GetAllTransaksiPagination(page, limit int) ([]Transaksi, int, error)
	GetAllHistoryTransaksiPagination(page, limit int) ([]TransaksiHistory, int, error)
	//
	UsedVoucher(voucher.VoucherUsed) error
	CheckVoucherExists(voucherID string) (bool, error)
}
