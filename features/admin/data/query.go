package data

import (
	"errors"
	"skripsi/features/admin"
	"skripsi/helper"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdminData struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *AdminData {
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
	err := d.DB.Where("email = ?", admins.Email).First(&adminData).Error
	if err != nil {
		return admin.Admin{}, err
	}

	if !helper.CheckPasswordHash(admins.Password, adminData.Password) {
		return admin.Admin{}, errors.New("wrong password")
	}
	return adminData, nil
}
