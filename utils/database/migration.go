package database

import (
	Admin "skripsi/features/admin/data"
	Instruktor "skripsi/features/instruktur/data"
	Kategori "skripsi/features/kategori/data"
	Kursus "skripsi/features/kursus/data"
	Users "skripsi/features/users/data"
	Voucher "skripsi/features/voucher/data"

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

	return nil
}
