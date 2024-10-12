package handler

type TransaksiRequest struct {
	// TotalHarga float64 `json:"total_harga"`
	VoucherID string `json:"voucher_id"`
	Kursus    string `json:"kursus"`
}
