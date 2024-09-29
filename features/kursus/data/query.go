package data

import (
	"gorm.io/gorm"
)

type KursusData struct {
	DB *gorm.DB
}

// func New(db *gorm.DB) kursus.KursusDataInterface {
// 	return &KursusData{
// 		DB: db,
// 	}
// }

// func (d *KursusData) GetKursusPagination(page int, limit int) ([]kursus.Kursus, int, error) {
// 	var dataKursus []kursus.Kursus
// 	var total int64

// 	count := d.DB.Model(&kursus.Kursus{}).Where("deleted_at IS NULL").Count(&total)
// 	if count.Error != nil {
// 		return nil, 0, constant.ErrDataNotfound
// 	}

// 	// ini menampilkan data didalam array nya misal di satu page nya itu limit 10, contoh (1-1)*10 = 10,  maka ada 10 data yang ditampilkan
// 	totalPages := int((total + int64(limit) - 1) / int64(limit))
// 	tx := d.DB.Where("deleted_at IS NULL").
// 		Offset((page - 1) * limit).
// 		Limit(limit).
// 		Find(&dataKursus)

// 	if tx.Error != nil {
// 		return nil, 0, constant.ErrGetData
// 	}
// 	return dataKursus, totalPages, nil
// }

// func (d *KursusData) GetAllKursus() ([]kursus.Kursus, error) {
// 	var data []kursus.Kursus
// 	if err := d.DB.Preload("Image").Preload("KategoriKursus.Kategori").Preload("Instruktur.Instruktur").Preload("MateriPembelajaran.Materi").Find(&data).Error; err != nil {
// 		return nil, err
// 	}
// 	return data, nil
// }

// func (d *KursusData) GetAllKursusById(id string) (kursus.Kursus, error) {
// 	var dataKursus kursus.Kursus
// 	if err := d.DB.Where("deleted_at IS NULL").Where("id = ?", id).First(&dataKursus).Error; err != nil {
// 		return kursus.Kursus{}, nil
// 	}
// 	return dataKursus, nil
// }

// func (d *KursusData) AddKursus(data kursus.Kursus) error {
// 	if err := d.DB.Create(&data).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (d *KursusData) UpdateKursus(data kursus.Kursus) error {
// 	var dataKursus kursus.Kursus
// 	if err := d.DB.Where("deleted_at IS NULL").Where("id = ?", data.ID).First(&dataKursus).Error; err != nil {
// 		return err
// 	}
// 	// update hanya satu value atau field, jika banyak Updates
// 	if err := d.DB.Model(&dataKursus).Updates(data).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }
