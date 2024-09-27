package database

import (
	Admin "skripsi/features/admin/data"
	Instruktor "skripsi/features/instruktur/data"
	Users "skripsi/features/users/data"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	db.AutoMigrate(&Users.User{})
	db.AutoMigrate(&Users.VerifyOtp{})
	db.AutoMigrate(&Admin.Admin{})
	db.AutoMigrate(&Instruktor.Instruktur{})
	return nil
}
