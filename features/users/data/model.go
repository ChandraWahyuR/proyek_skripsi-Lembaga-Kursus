package data

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	*gorm.Model
	ID         string `gorm:"primary_key;type:varchar(50);not null;column:id"`
	Username   string `gorm:"type:varchar(255);not null;column:username"`
	Email      string `gorm:"type:varchar(255);not null;column:email"`
	Password   string `gorm:"type:varchar(255);not null;column:password"`
	NomorHP    string `gorm:"type:varchar(255);not null;column:nomor_hp"`
	ProfileUrl string `gorm:"type:varchar(255);not null;column:profile_url"`
	IsActive   bool   `gorm:"not null;column:is_active"`
	//
	Nama          string    `gorm:"type:varchar(255);not null;column:nama"`
	NIS           string    `gorm:"type:varchar(255);not null;column:nis"`
	Agama         string    `gorm:"type:varchar(255);not null;column:agama"`
	Gender        string    `gorm:"type:varchar(255);not null;column:gender"`
	TempatLahir   string    `gorm:"type:varchar(255);not null;column:tempat_lahir"`
	TanggalLahir  time.Time `gorm:"type:TIMESTAMP;null;default:null"`
	OrangTua      string    `gorm:"type:varchar(255);not null;column:orang_tua"`
	Profesi       string    `gorm:"type:varchar(255);not null;column:profesi"`
	Ijazah        string    `gorm:"type:varchar(255);not null;column:ijazah"`
	KTP           string    `gorm:"type:varchar(255);not null;column:ktp"`
	KartuKeluarga string    `gorm:"type:varchar(255);not null;column:kartu_keluarga"`
}

type VerifyOtp struct {
	*gorm.Model
	ID        string    `gorm:"type:varchar(50);not null;column:id"`
	Email     string    `gorm:"type:varchar(255);not null;column:email"`
	Otp       string    `gorm:"type:varchar(255);not null;column:otp"`
	Status    string    `gorm:"type:varchar(50);not null;column:status"`
	ExpiredAt time.Time `gorm:"not null;column:expired_at"`
}

func (u *User) TableName() string {
	return "users"
}
func (u *VerifyOtp) TableName() string {
	return "verify_otps"
}
