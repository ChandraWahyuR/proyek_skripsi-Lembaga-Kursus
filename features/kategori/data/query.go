package data

import (
	"skripsi/constant"
	"skripsi/features/kategori"

	"gorm.io/gorm"
)

type KategoriData struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *KategoriData {
	return &KategoriData{
		DB: db,
	}
}

func (d KategoriData) GetKategoriWithPagination(page int, limit int) ([]kategori.Kategori, int, error) {
	var dataKategori []kategori.Kategori
	var total int64

	count := d.DB.Model(&kategori.Kategori{}).Where("deleted_at IS NULL").Count(&total)
	if count.Error != nil {
		return nil, 0, constant.ErrDataNotfound
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	tx := d.DB.Where("deleted_at IS NULL").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&dataKategori)

	if tx.Error != nil {
		return nil, 0, constant.ErrGetData
	}

	return dataKategori, totalPages, nil
}

func (d KategoriData) GetAllKategori() ([]kategori.Kategori, error) {
	var dataKategori []kategori.Kategori
	if err := d.DB.Where("deleted_at IS NULL").Find(&dataKategori).Error; err != nil {
		return nil, err
	}
	return dataKategori, nil
}

func (d KategoriData) GetKategoriById(id string) (kategori.Kategori, error) {
	var dataKategori kategori.Kategori
	if err := d.DB.Where("deleted_at IS NULL").Where("id = ?", id).Find(&dataKategori).Error; err != nil {
		return dataKategori, nil
	}
	return dataKategori, nil
}

func (d KategoriData) CreateKategori(data kategori.Kategori) error {
	if err := d.DB.Create(&data).Error; err != nil {
		return err
	}

	return nil
}

func (d KategoriData) UpdateKategori(data kategori.Kategori) error {
	var existingKategori kategori.Kategori

	if err := d.DB.Where("id = ?", data.ID).Where("deleted_at IS NULL").First(&existingKategori).Error; err != nil {
		return err
	}

	if err := d.DB.Model(&existingKategori).Updates(data).Error; err != nil {
		return err
	}

	return nil
}

func (d KategoriData) DeleteKategori(id string) error {
	res := d.DB.Begin()

	if err := res.Where("id = ?", id).Delete(&Kategori{}); err.Error != nil {
		res.Rollback()
		return constant.ErrKategoriNotFound
	} else if err.RowsAffected == 0 {
		res.Rollback()
		return constant.ErrKategoriNotFound
	}

	return res.Commit().Error
}
