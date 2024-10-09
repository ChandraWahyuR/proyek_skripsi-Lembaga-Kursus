package handler

import "time"

type RequestCreateVoucher struct {
	Nama      string    `json:"nama"`
	Deskripsi string    `json:"deskripsi"`
	Code      string    `json:"code"`
	Discount  float64   `json:"discount"`
	ExpiredAt time.Time `json:"expired_at"`
}
