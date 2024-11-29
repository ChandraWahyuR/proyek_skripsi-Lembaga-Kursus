package voucher

import (
	"time"

	"github.com/labstack/echo/v4"
)

type Voucher struct {
	ID        string
	Nama      string
	Deskripsi string
	Code      string
	Discount  float64
	ExpiredAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type VoucherUsed struct {
	ID        string
	VoucherID string
	UserID    string
}

type VoucherHandlerInterface interface {
	GetAllVoucher() echo.HandlerFunc
	GetByIDVoucher() echo.HandlerFunc
	CreateVoucher() echo.HandlerFunc
	UpdateVoucher() echo.HandlerFunc
	DeleteVoucher() echo.HandlerFunc
}

type VoucherDataInterface interface {
	GetAllVoucher() ([]Voucher, error)
	GetAllVoucherPagination(page, limit int) ([]Voucher, int, error)
	GetByIDVoucher(id string) (Voucher, error)
	CreateVoucher(Voucher) error
	UpdateVoucher(Voucher) error
	DeleteVoucher(id string) error
}

type VoucherServiceInterface interface {
	GetAllVoucher() ([]Voucher, error)
	GetAllVoucherPagination(page, limit int) ([]Voucher, int, error)
	GetByIDVoucher(id string) (Voucher, error)
	CreateVoucher(Voucher) error
	UpdateVoucher(Voucher) error
	DeleteVoucher(id string) error
}
