package database

import (
	Admin "skripsi/features/admin/data"
	Instruktor "skripsi/features/instruktur/data"
	Jadwal "skripsi/features/jadwal_mengajar/data"
	Kategori "skripsi/features/kategori/data"
	Kursus "skripsi/features/kursus/data"
	Transaksi "skripsi/features/transaksi/data"
	Users "skripsi/features/users/data"
	Voucher "skripsi/features/voucher/data"
	Webhook "skripsi/features/webhook/data"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	db.AutoMigrate(&Users.User{})
	db.AutoMigrate(&Users.VerifyOtp{})
	db.AutoMigrate(&Admin.Admin{})
	db.AutoMigrate(&Instruktor.Instruktur{})
	db.AutoMigrate(&Kategori.Kategori{})
	// Kursus
	db.AutoMigrate(&Kursus.Kursus{})
	db.AutoMigrate(&Kursus.ImageKursus{})
	db.AutoMigrate(&Kursus.KategoriKursus{})
	db.AutoMigrate(&Kursus.MateriPembelajaran{})
	db.AutoMigrate(&Kursus.JadwalKursus{})
	// Voucher
	db.AutoMigrate(&Voucher.Voucher{})
	db.AutoMigrate(&Voucher.VoucherUsed{})
	// Transaksi
	db.AutoMigrate(&Transaksi.Transaksi{})
	db.AutoMigrate(&Transaksi.TransaksiHistory{})
	db.AutoMigrate(&Jadwal.JadwalMengajar{})
	// Webhook
	db.AutoMigrate(&Webhook.PaymentNotification{})

	return nil
}
