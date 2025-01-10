package data

import (
	"errors"
	"log"
	"skripsi/constant"
	"skripsi/features/admin"
	"skripsi/helper"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdminData struct {
	DB *gorm.DB
}

func New(db *gorm.DB) admin.AdminDataInterface {
	return &AdminData{
		DB: db,
	}
}

func (u *AdminData) IsEmailExist(email string) bool {
	var userData Admin
	if err := u.DB.Where("email = ?", email).First(&userData).Error; err != nil {
		return false
	}
	return true
}

func (u *AdminData) IsUsernameExist(username string) bool {
	var userData Admin
	if err := u.DB.Where("username = ?", username).First(&userData).Error; err != nil {
		return false
	}
	return true
}

func (d *AdminData) RegisterAdmin(admins admin.Admin) error {
	if d.IsEmailExist(admins.Email) {
		return errors.New("email already exist")
	}

	if d.IsUsernameExist(admins.Username) {
		return errors.New("username already exist")
	}

	adminData := Admin{
		ID:       uuid.New().String(),
		Username: admins.Username,
		Email:    admins.Email,
		Password: admins.Password,
	}
	err := d.DB.Create(&adminData).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *AdminData) LoginAdmin(admins admin.Admin) (admin.Admin, error) {
	var adminData admin.Admin
	err := d.DB.Where("username = ?", admins.Username).First(&adminData).Error
	if err != nil {
		return admin.Admin{}, constant.ErrAdminNotFound
	}

	if !helper.CheckPasswordHash(admins.Password, adminData.Password) {
		return admin.Admin{}, constant.ErrPasswordNotMatch
	}
	return adminData, nil
}

func (d *AdminData) DownloadLaporanPembelian(startDate, endDate time.Time) ([]map[string]interface{}, error) {
	var histories []map[string]interface{}

	err := d.DB.Table("transaksi_histories").
		Select("transaksi_histories.id, transaksi_histories.transaksi_id, transaksi_histories.kursus_id, transaksi_histories.user_id, users.nama AS user_nama, users.email AS email, kursus.nama AS nama_kursus,transaksi_histories.status, transaksi_histories.valid_until, transaksis.total_harga, transaksis.status AS transaksi_status").
		Joins("JOIN kursus ON kursus.id = transaksi_histories.kursus_id").
		Joins("JOIN transaksis ON transaksis.id = transaksi_histories.transaksi_id").
		Joins("JOIN users ON users.id = transaksi_histories.user_id").
		Where("transaksi_histories.created_at BETWEEN ? AND ?", startDate, endDate).
		Scan(&histories).Error
	if err != nil {
		return nil, err
	}
	log.Println("Histories:", histories)

	return histories, nil
}
