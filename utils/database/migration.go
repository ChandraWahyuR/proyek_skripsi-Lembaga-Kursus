package database

import (
	Admin "skripsi/features/admin/data"
	Users "skripsi/features/users/data"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	db.AutoMigrate(&Users.User{})
	db.AutoMigrate(&Users.VerifyOtp{})
	db.AutoMigrate(&Admin.Admin{})
	return nil
}
