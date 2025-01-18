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

	count := d.DB.Model(&kursus.Kursus{}).Where("deleted_at IS NULL").Count(&total)
	if count.Error != nil {
		return nil, 0, constant.ErrDataNotfound
	}

	// ini menampilkan data didalam array nya misal di satu page nya itu limit 10, contoh (1-1)*10 = 10,  maka ada 10 data yang ditampilkan
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	tx := d.DB.Preload("Image", "deleted_at IS NULL").Preload("Jadwal", "deleted_at IS NULL").Where("deleted_at IS NULL").
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
	if err := d.DB.Preload("Image").Preload("Kategori.Kategori").Preload("Instruktur").Preload("MateriPembelajaran.Materi").Preload("Jadwal").Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (d *KursusData) GetAllKursusById(id string) (kursus.Kursus, error) {
	var dataKursus kursus.Kursus
	if err := d.DB.Model(&kursus.Kursus{}).Preload("Image", "deleted_at IS NULL").Preload("Jadwal", "deleted_at IS NULL").Preload("Instruktur", "deleted_at IS NULL").Preload("Kategori.Kategori", "deleted_at IS NULL").Preload("MateriPembelajaran", "deleted_at IS NULL").Where("id = ? AND deleted_at IS NULL", id).First(&dataKursus).Error; err != nil {
		return kursus.Kursus{}, constant.ErrGetID
	}
	return dataKursus, nil
}

func (d *KursusData) AddKursus(data kursus.Kursus) error {
	// Mulai operasi transaksi, cocok untuk table yang saling berlerasi
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
	tx := d.DB.Begin()

	// Ambil data kursus yang ada
	var dataKursus kursus.Kursus
	if err := tx.Where("id = ?", data.ID).Where("deleted_at IS NULL").First(&dataKursus).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Hapus entri lama menggunakan metode delete yang sudah ada
	if err := d.DeleteMateriKursus(dataKursus.ID); err != nil {
		tx.Rollback()
		return err
	}

	if err := d.DeleteImageKursus(dataKursus.ID); err != nil {
		tx.Rollback()
		return err
	}

	if err := d.DeleteKategoriKursus(dataKursus.ID); err != nil {
		tx.Rollback()
		return err
	}

	if err := d.DeleteJadwalKurusus(dataKursus.ID); err != nil {
		tx.Rollback()
		return err
	}

	// Tambahkan materi baru
	if err := tx.Model(&dataKursus).Association("MateriPembelajaran").Append(&data.MateriPembelajaran); err != nil {
		tx.Rollback()
		return err
	}

	// Tambahkan gambar baru
	if err := tx.Model(&dataKursus).Association("Image").Append(data.Image); err != nil {
		tx.Rollback()
		return err
	}

	// Tambahkan kategori baru
	if err := tx.Model(&dataKursus).Association("Kategori").Append(data.Kategori); err != nil {
		tx.Rollback()
		return err
	}

	// Tambahan Jadwal
	if err := tx.Model(&dataKursus).Association("Jadwal").Append(data.Jadwal); err != nil {
		tx.Rollback()
		return err
	}

	// Update kursus
	if err := tx.Model(&dataKursus).Updates(data).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (d *KursusData) DeleteKursus(id string) error {
	res := d.DB.Begin()

	if err := res.Where("kursus_id = ?", id).Delete(&ImageKursus{}); err.Error != nil {
		res.Rollback()
		return constant.ErrImageKursusNotFound
	}

	if err := res.Where("kursus_id = ?", id).Delete(&KategoriKursus{}); err.Error != nil {
		res.Rollback()
		return constant.ErrKategoriKursusNotFound
	}
	if err := res.Where("kursus_id = ?", id).Delete(&MateriPembelajaran{}); err.Error != nil {
		res.Rollback()
		return constant.ErrMateriKursusNotFound
	}

	if err := res.Where("deleted_at IS NULL").Where("id = ?", id).Delete(&Kursus{}); err.Error != nil {
		res.Rollback()
		return constant.ErrKursusNotFound
	} else if err.RowsAffected == 0 {
		res.Rollback()
		return constant.ErrKursusNotFound
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
		Where("nama LIKE ?", "%"+name+"%").
		Count(&total)
	if count.Error != nil {
		return nil, 0, count.Error
	}
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	// Query dengan pagination
	tx := d.DB.Where("nama LIKE ?", "%"+name+"%").Preload("Image", "deleted_at IS NULL").Preload("Jadwal", "deleted_at IS NULL").Where("deleted_at IS NULL").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&dataKursus)

	if tx.Error != nil {
		return nil, 0, tx.Error
	}

	return dataKursus, totalPages, nil
}

func (d *KursusData) DeleteImageKursus(id string) error {
	res := d.DB.Begin()
	if err := res.Where("kursus_id = ?", id).Delete(&ImageKursus{}); err.Error != nil {
		res.Rollback()
		return constant.ErrImageKursusNotFound
	}

	return res.Commit().Error
}
func (d *KursusData) DeleteMateriKursus(id string) error {
	res := d.DB.Begin()
	if err := res.Where("kursus_id = ?", id).Delete(&MateriPembelajaran{}); err.Error != nil {
		res.Rollback()
		return constant.ErrMateriKursusNotFound
	}

	return res.Commit().Error
}

func (d *KursusData) DeleteKategoriKursus(id string) error {
	res := d.DB.Begin()
	if err := res.Where("kursus_id = ?", id).Delete(&KategoriKursus{}); err.Error != nil {
		res.Rollback()
		return constant.ErrKategoriKursusNotFound
	}

	return res.Commit().Error
}

func (d *KursusData) DeleteJadwalKurusus(id string) error {
	res := d.DB.Begin()
	if err := res.Where("kursus_id = ?", id).Delete(&JadwalKursus{}); err.Error != nil {
		res.Rollback()
		return constant.ErrKategoriKursusNotFound
	}

	return res.Commit().Error
}
