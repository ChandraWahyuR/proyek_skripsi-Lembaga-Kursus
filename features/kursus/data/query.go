package data

import (
	"skripsi/constant"
	"skripsi/features/kursus"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type KursusData struct {
	DB *gorm.DB
}

// ini jika return nya kursusData methodnya ngga dipanggil semua diinterface masih bisa kelemahan ke akses semua kalau pakai ini data method yang dibuat sesaui yang ada di handler interface
func New(db *gorm.DB) kursus.KursusDataInterface {
	return &KursusData{
		DB: db,
	}
}

func (d *KursusData) GetKursusPagination(page int, limit int) ([]kursus.Kursus, int, error) {
	var dataKursus []kursus.Kursus
	var total int64

	count := d.DB.Model(&Kursus{}).Where("deleted_at IS NULL").Count(&total)
	if count.Error != nil {
		return nil, 0, constant.ErrDataNotfound
	}

	// ini menampilkan data didalam array nya misal di satu page nya itu limit 10, contoh (1-1)*10 = 10,  maka ada 10 data yang ditampilkan
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	tx := d.DB.Preload("Image").Where("deleted_at IS NULL").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&dataKursus)

	if tx.Error != nil {
		return nil, 0, constant.ErrGetData
	}
	return dataKursus, totalPages, nil
}

func (d *KursusData) GetAllKursus() ([]kursus.Kursus, error) {
	var data []kursus.Kursus
	if err := d.DB.Preload("Image").Preload("Kategori.Kategori").Preload("Instruktur").Preload("MateriPembelajaran.Materi").Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (d *KursusData) GetAllKursusById(id string) (kursus.Kursus, error) {
	var dataKursus kursus.Kursus
	if err := d.DB.Model(&Kursus{}).Preload("Image").Preload("Kategori.Kategori").Preload("MateriPembelajaran.Materi").Where("deleted_at IS NULL").Where("id = ?", id).First(&dataKursus).Error; err != nil {
		return kursus.Kursus{}, err
	}

	return dataKursus, nil
}

func (d *KursusData) AddKursus(data kursus.Kursus) error {
	// Mulai transaksi
	tx := d.DB.Begin()

	// Generate UUID untuk kursus jika belum di-set
	if data.ID == "" {
		data.ID = uuid.New().String()
	}

	// Menyimpan data kursus utama
	if err := tx.Create(&data).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (d *KursusData) UpdateKursus(data kursus.Kursus) error {
	var dataKursus kursus.Kursus
	if err := d.DB.Where("deleted_at IS NULL").Where("id = ?", data.ID).First(&dataKursus).Error; err != nil {
		return err
	}
	// update hanya satu value atau field, jika banyak Updates
	if err := d.DB.Model(&dataKursus).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

func (d *KursusData) DeleteKursus(id string) error {
	res := d.DB.Begin()

	if err := res.Where("deleted_at IS NULL").Where("id = ?", id).Delete(&Kursus{}); err.Error != nil {
		res.Rollback()
		return constant.ErrInstrukturNotFound
	} else if err.RowsAffected == 0 {
		res.Rollback()
		return constant.ErrInstrukturNotFound
	}

	return res.Commit().Error
}

func (d *KursusData) GetAllKursusByName(name string, page int, limit int) ([]kursus.Kursus, int, error) {
	var dataKursus []kursus.Kursus
	var total int64

	// Hitung jumlah data yang sesuai dengan filter di model Kursus,
	// lalu simpan hasil perhitungannya ke variabel `total` dalam bentuk integer.
	// Jika ada error selama proses, kembalikan error tersebut.
	count := d.DB.Model(&kursus.Kursus{}).
		Where("name LIKE ?", "%"+name+"%").
		Count(&total)
	if count.Error != nil {
		return nil, 0, count.Error
	}
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	// Query dengan pagination
	tx := d.DB.Where("name LIKE ?", "%"+name+"%").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&dataKursus)

	if tx.Error != nil {
		return nil, 0, tx.Error
	}

	return dataKursus, totalPages, nil
}
