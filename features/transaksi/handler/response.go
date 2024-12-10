package handler

import (
	"time"
)

type PaymentResponse struct {
	Amount  int    `json:"amount"`
	SnapURL string `json:"snap_url"`
}

type MidtransNotification struct {
	OrderID           string `json:"order_id"`
	TransactionStatus string `json:"transaction_status"`
	FraudStatus       string `json:"fraud_status"`
}

type MetadataResponse struct {
	TotalPage int `json:"total_page"`
	Page      int `json:"current_page"`
}

// Transaksi
type GetAllTransaksiAdminResponse struct {
	ID         string  `json:"id"`
	User       User    `json:"user"`
	Kursus     Kursus  `json:"kursus"`
	TotalHarga float64 `json:"total_harga"`
	Status     string  `json:"status"`
}

type GetAllTransaksiByIDAdminResponse struct {
	ID         string  `json:"id"`
	User       User    `json:"user"`
	Kursus     Kursus  `json:"kursus"`
	VoucherID  string  `json:"voucher_id"`
	TotalHarga float64 `json:"total_harga"`
	Status     string  `json:"status"`
}

type GetAllTransaksiUserResponse struct {
	ID         string  `json:"id"`
	Kursus     Kursus  `json:"kursus"`
	TotalHarga float64 `json:"total_harga"`
	Status     string  `json:"status"`
}

type Kursus struct {
	ID   string `json:"id"`
	Nama string `json:"nama"`
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// History Transaksi
type GetAllHistoryAdminResponse struct {
	ID         string    `json:"id"`
	Transaksi  Transaksi `json:"transaksi"`
	KursusID   string    `json:"kursus_id"`
	UserID     string    `json:"user_id"`
	Status     string    `json:"status"`
	Voucher    string    `json:"voucher"`
	ValidUntil time.Time `json:"valid_until"`
}

type GetHistoryAdminByIDResponse struct {
	ID         string    `json:"id"`
	Transaksi  Transaksi `json:"transaksi"`
	Kursus     Kursus    `json:"kursus"`
	User       User      `json:"user"`
	VoucherID  string    `json:"voucher_id"`
	TotalHarga float64   `json:"total_harga"`
	Status     string    `json:"status"`
	ValidUntil time.Time `json:"valid_until"`
}
type GetTransaksiResponse struct {
	ID         string    `json:"id"`
	Transaksi  Transaksi `json:"transaksi"`
	Kursus     Kursus    `json:"kursus"`
	User       User      `json:"user"`
	VoucherID  string    `json:"voucher_id"`
	TotalHarga float64   `json:"total_harga"`
}

type Transaksi struct {
	ID         string  `json:"id"`
	TotalHarga float64 `json:"total_harga"`
}

type GetAllHistoryUserResponse struct {
	ID         string    `json:"id"`
	Transaksi  Transaksi `json:"transaksi"`
	KursusID   string    `json:"kursus_id"`
	UserID     string    `json:"user_id"`
	ValidUntil time.Time `json:"valid_until"`
}
