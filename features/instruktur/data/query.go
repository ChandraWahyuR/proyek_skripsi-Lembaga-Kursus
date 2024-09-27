package data

import (
	"skripsi/constant"
	"skripsi/features/instruktur"

	"gorm.io/gorm"
)

type InstrukturData struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *InstrukturData {
	return &InstrukturData{
		DB: db,
	}
}

// Pagination
func (d *InstrukturData) GetInstrukturWithPagination(page int, limit int) ([]instruktur.Instruktur, int, error) {
	var dataInstruktur []instruktur.Instruktur
	var total int64

	// Hitung total instruktur untuk pagination
	count := d.DB.Model(&instruktur.Instruktur{}).Where("deleted_at IS NULL").Count(&total)
	if count.Error != nil {
		return nil, 0, constant.ErrInstrukturNotFound
	}

	// Limit jika lupa
	// if limit <= 0 {
	// 	limit = 20
	// }

	// Hitung total halaman
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	// Ambil instruktur berdasarkan pagination
	tx := d.DB.Where("deleted_at IS NULL").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&dataInstruktur)

	if tx.Error != nil {
		return nil, 0, constant.ErrGetInstruktur
	}

	return dataInstruktur, totalPages, nil
}

func (d InstrukturData) GetAllInstruktur() ([]instruktur.Instruktur, error) {
	var dataInstruktur []instruktur.Instruktur
	if err := d.DB.Where("deleted_at IS NULL").Find(&dataInstruktur).Error; err != nil {
		return nil, err
	}
	return dataInstruktur, nil
}

func (d InstrukturData) GetAllInstrukturByID(id string) (instruktur.Instruktur, error) {
	var dataInstruktur instruktur.Instruktur
	if err := d.DB.Where("id = ?", id).First(&dataInstruktur).Error; err != nil {
		return dataInstruktur, err
	}
	return dataInstruktur, nil
}

func (d InstrukturData) PostInstruktur(data instruktur.Instruktur) error {
	if err := d.DB.Create(&data).Error; err != nil {
		return err
	}
	return nil
}

func (d InstrukturData) UpdateInstruktur(data instruktur.UpdateInstruktur) error {
	var exisistingInstruktur instruktur.Instruktur
	if err := d.DB.Where("id = ?", data.ID).First(&exisistingInstruktur).Error; err != nil {
		return err
	}

	if err := d.DB.Model(&exisistingInstruktur).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

func (d InstrukturData) DeleteInstruktur(id string) error {
	res := d.DB.Begin()

	if err := res.Where("id = ?", id).Delete(&Instruktur{}); err.Error != nil {
		res.Rollback()
		return constant.ErrInstrukturNotFound
	} else if err.RowsAffected == 0 {
		res.Rollback()
		return constant.ErrInstrukturNotFound
	}

	return res.Commit().Error
}

func (d *InstrukturData) GetInstruktorByName(name string, page int, limit int) ([]instruktur.Instruktur, int, error) {
	var instrukturs []instruktur.Instruktur
	var total int64

	// Hitung total data berdasarkan nama
	count := d.DB.Model(&instruktur.Instruktur{}).
		Where("name LIKE ?", "%"+name+"%").
		Count(&total)
	if count.Error != nil {
		return nil, 0, count.Error
	}

	// Hitung total halaman
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	// Query dengan pagination
	tx := d.DB.Where("name LIKE ?", "%"+name+"%").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&instrukturs)

	if tx.Error != nil {
		return nil, 0, tx.Error
	}

	return instrukturs, totalPages, nil
}
