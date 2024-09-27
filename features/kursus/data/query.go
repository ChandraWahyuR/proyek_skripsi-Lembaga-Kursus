package data

import (
	"gorm.io/gorm"
)

type KursusData struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *KursusData {
	return &KursusData{
		DB: db,
	}
}

// func (d *KursusData) GetAllKursus() ([]Kursus, error) {
// 	var data []kursus.Kursus
// 	if err := d.DB.Preload("ImageKursus").Preload("MateriPembelajaran").Preload("Instruktur").Preload("KategoriKursus").Find(&data).Error; err != nil {

// 	}
// }
