package data

import (
	"errors"
	"time"

	"skripsi/constant"
	"skripsi/features/users"
	"skripsi/helper"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserData struct {
	DB *gorm.DB
}

func New(db *gorm.DB) users.UserDataInterface {
	return &UserData{
		DB: db,
	}
}

func (u *UserData) Register(user users.User) error {
	isUsername := u.IsUsernameExist(user.Username)
	if isUsername {
		return errors.New("username already exist")
	}

	isEmail := u.IsEmailExist(user.Email)
	if isEmail {
		return errors.New("email already exist")
	}

	userData := User{
		ID:       uuid.New().String(),
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
		NomorHP:  user.NomorHP,
		NIS:      helper.GenerateNis(),
		IsActive: false,
	}

	if err := u.DB.Create(&userData).Error; err != nil {
		return err
	}
	return nil
}

func (u *UserData) Login(user users.User) (users.User, error) {
	var userData users.User
	err := u.DB.Where("email = ?", user.Email).First(&userData).Error
	if err != nil {
		return users.User{}, errors.New("email tidak ada")
	}

	if !userData.IsActive {
		return users.User{}, errors.New("user not active")
	}

	if !helper.CheckPasswordHash(user.Password, userData.Password) {
		return users.User{}, errors.New("wrong password")
	}
	return userData, nil
}

// Forgot Password
func (d *UserData) ForgotPassword(forgot users.ForgotPassword) error {
	d.DB.Where("email = ?", forgot.Email).Delete(&VerifyOtp{})

	forgotData := VerifyOtp{
		ID:        uuid.New().String(),
		Email:     forgot.Email,
		Otp:       forgot.Otp,
		ExpiredAt: time.Now().Add(10 * time.Minute),
	}

	if err := d.DB.Create(&forgotData).Error; err != nil {
		return err
	}
	return nil
}

func (d *UserData) VerifyOTP(verify users.VerifyOtp) error {
	var otp VerifyOtp
	err := d.DB.Model(&VerifyOtp{}).Where("email = ? AND otp = ?", verify.Email, verify.Otp).First(&otp).Error
	if err != nil {
		return err
	}

	if time.Now().After(otp.ExpiredAt) {
		return errors.New("OTP has expired")
	}

	if verify.Otp != otp.Otp {
		return errors.New("wrong otp")
	}

	err = d.DB.Model(&otp).Update("status", verify.Status).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *UserData) ResetPassword(change users.ResetPassword) error {
	var userData User
	err := d.DB.Where("email = ?", change.Email).First(&userData).Error
	if err != nil {
		return err
	}

	if err := d.DB.Model(&User{}).Where("email = ?", change.Email).Update("password", change.Password).Error; err != nil {
		return err
	}

	return nil
}

// ==========================================================================================================
// Check
func (u *UserData) IsEmailExist(email string) bool {
	var userData User
	if err := u.DB.Where("email = ?", email).First(&userData).Error; err != nil {
		return false
	}
	return true
}

func (u *UserData) IsUsernameExist(username string) bool {
	var userData User
	if err := u.DB.Where("username = ?", username).First(&userData).Error; err != nil {
		return false
	}
	return true
}

func (d *UserData) GetByEmail(email string) (users.User, error) {
	var user users.User
	err := d.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return users.User{}, err
	}
	return user, nil
}

func (d *UserData) VerifyEmail(email string, isValid bool) error {
	err := d.DB.Model(&User{}).Where("email = ?", email).Update("is_active", isValid).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *UserData) GetAllUserPagination(page, limit int) ([]users.GetUser, int, error) {
	var data []users.GetUser
	var total int64

	count := d.DB.Model(&users.GetUser{}).Where("deleted_at IS NULL").Count(&total)
	if count.Error != nil {
		return nil, 0, constant.ErrDataNotfound
	}
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	tx := d.DB.
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&data)

	if tx.Error != nil {
		return nil, 0, constant.ErrGetData
	}
	return data, totalPages, nil
}

func (d *UserData) GetUserByID(userId string) (users.User, error) {
	var dataUser users.User
	err := d.DB.Where("id = ?", userId).First(&dataUser).Error
	if err != nil {
		return users.User{}, constant.ErrGetID
	}
	return dataUser, nil
}

func (d *UserData) UpdateUser(data users.EditUser) error {
	tx := d.DB.Begin()

	var dataUsers users.EditUser
	if err := tx.Where("id = ?", data.ID).Where("deleted_at IS NULL").First(&dataUsers).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&dataUsers).Updates(data).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (d *UserData) DeleteUser(userId string) error {
	res := d.DB.Begin()
	if err := res.Where("deleted_at IS NULL").Where("id = ?", userId).Delete(&User{}); err.Error != nil {
		res.Rollback()
		return constant.ErrUserNotFound
	} else if err.RowsAffected == 0 {
		res.Rollback()
		return constant.ErrFailedDelete
	}
	return res.Commit().Error
}
