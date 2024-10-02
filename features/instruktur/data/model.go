package data

import "gorm.io/gorm"

type Instruktur struct {
	*gorm.Model
	ID     string `gorm:"type:varchar(50);primaryKey;not null;column:id"`
	Name   string `gorm:"type:varchar(255);not null;column:name"`
	Gender string `gorm:"type:varchar(255);not null;column:gender"`
	Email  string `gorm:"type:varchar(255);not null;column:email"`
	Alamat string `gorm:"type:varchar(255);not null;column:alamat"`
	NoHp   string `gorm:"type:varchar(255);not null;column:no_hp"`
}

func (Instruktur) TableName() string {
	return "instrukturs"
}
