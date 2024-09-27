package data

import "gorm.io/gorm"

type Kategori struct {
	*gorm.Model
	ID        string `gorm:"type:varchar(50);not null;column:id"`
	Nama      string `gorm:"type:varchar(255);not null;column:nama"`
	Deskripsi string `gorm:"type:text;not null;column:deskripsi"`
	ImageUrl  string `gorm:"type:text;not null;column:image_url"`
}

func (Kategori) TableName() string {
	return "kategoris"
}
