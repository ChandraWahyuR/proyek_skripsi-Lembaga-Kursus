package handler

import "time"

type MetadataResponse struct {
	TotalPage int `json:"total_page"`
	Page      int `json:"page"`
}
type ResponseGetAllVoucher struct {
	ID        string    `json:"id"`
	Nama      string    `json:"nama"`
	Code      string    `json:"code"`
	Discount  float64   `json:"discount"`
	ExpiredAt time.Time `json:"expired_at"`
}

type ResponseGetIDVoucher struct {
	ID        string    `json:"id"`
	Nama      string    `json:"nama"`
	Deskripsi string    `json:"deskripsi"`
	Code      string    `json:"code"`
	Discount  float64   `json:"discount"`
	ExpiredAt time.Time `json:"expired_at"`
}
