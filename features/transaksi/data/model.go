package data

import (
	Kursus "skripsi/features/kursus/data"
	Users "skripsi/features/users/data"
	"time"

	"gorm.io/gorm"
)

type Transaksi struct {
	*gorm.Model
	ID         string        `gorm:"type:varchar(50);primaryKey;not null;column:id"`
	TotalHarga float64       `gorm:"type:decimal(10,2);not null;column:total_harga"`
	VoucherID  string        `gorm:"type:varchar(50);column:voucher_id"`
	KursusID   string        `gorm:"type:varchar(50);not null;column:kursus_id"`
	Kursus     Kursus.Kursus `gorm:"foreignKey:KursusID;references:ID"`
	UserID     string        `gorm:"type:varchar(50);not null;column:user_id"`
	User       Users.User    `gorm:"foreignKey:UserID;references:ID"`
	SnapURL    string        `gorm:"type:varchar(255);not null;column:snap_url"`
	Status     string        `gorm:"type:varchar(50);not null;column:status"`
}

type TransaksiHistory struct {
	*gorm.Model
	ID          string        `gorm:"type:varchar(50);primaryKey;not null;column:id"`
	TransaksiID string        `gorm:"type:varchar(50);not null;column:transaksi_id"`
	Transaksi   Transaksi     `gorm:"foreignKey:TransaksiID;references:ID"`
	KursusID    string        `gorm:"type:varchar(50);not null;column:kursus_id"`
	Kursus      Kursus.Kursus `gorm:"foreignKey:KursusID;references:ID"`
	UserID      string        `gorm:"type:varchar(50);not null;column:user_id"`
	User        Users.User    `gorm:"foreignKey:UserID;references:ID"`
	Status      string        `gorm:"type:varchar(50);not null;column:status"`
	ValidUntil  time.Time     `gorm:"not null;column:valid_until"`
}

func (Transaksi) TableName() string {
	return "transaksis"
}

func (TransaksiHistory) TableName() string {
	return "transaksi_histories"
}
