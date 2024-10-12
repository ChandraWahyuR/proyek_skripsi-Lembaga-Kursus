package data

import (
	"time"

	"gorm.io/gorm"
)

type Voucher struct {
	*gorm.Model
	ID        string    `gorm:"type:varchar(50);primaryKey;not null;column:id"`
	Nama      string    `gorm:"type:varchar(255);not null;column:nama"`
	Deskripsi string    `gorm:"type:text;not null;column:deskripsi"`
	Code      string    `gorm:"type:varchar(255);not null;column:code"`
	Discount  float64   `gorm:"type:decimal(10,2);not null;column:discount"`
	ExpiredAt time.Time `gorm:"type:datetime;not null;column:expired_at"`
}

func (Voucher) TableName() string {
	return "vouchers"
}
