package seeders

import (
	"fmt"
	"skripsi/features/voucher/data"
	"skripsi/helper"
	"time"
)

func ParseDate(dateString string) (time.Time, error) {
	return time.Parse(time.RFC3339, dateString)
}

func (s *Seeder) SeedVoucher() {
	vouchers := []data.Voucher{
		{
			ID:        "b1e974db-16aa-4f12-a999-8e13ebe5ab06",
			Nama:      "Voucher Akhir Tahun",
			Deskripsi: "Voucher memeringati akhir tahun",
			Code:      helper.GenerateCode(),
			Discount:  10,
			ExpiredAt: time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC),
		},
		{
			ID:        "3dad39d4-4dcb-4d00-a83c-39d10a9c99a1",
			Nama:      "Voucher Diskon Hari Kemerdekaan",
			Deskripsi: "Voucher untuk diskon dalam memperingati kemerdekaan indonesia",
			Code:      helper.GenerateCode(),
			Discount:  15,
			ExpiredAt: time.Date(2025, 6, 30, 23, 59, 59, 0, time.UTC),
		},
		{
			ID:        "15270598-8731-41ea-b3a3-12e18ef14206",
			Nama:      "Voucher Awal Tahun",
			Deskripsi: "Voucher untuk diskon awal tahun",
			Code:      helper.GenerateCode(),
			Discount:  5,
			ExpiredAt: time.Date(2025, 1, 07, 23, 59, 59, 0, time.UTC),
		},
	}
	for _, voucher := range vouchers {
		result := s.db.FirstOrCreate(&voucher, data.Voucher{ID: voucher.ID})
		if result.Error != nil {
			fmt.Printf("Failed to seed voucher %v: %v\n", voucher.ID, result.Error)
		}
	}
}
