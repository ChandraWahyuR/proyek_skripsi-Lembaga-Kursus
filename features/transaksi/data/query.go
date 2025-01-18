package data

import (
	"log"
	"skripsi/constant"
	"skripsi/features/kursus"
	"skripsi/features/transaksi"
	"skripsi/features/users"
	"skripsi/features/voucher"
	"time"

	"gorm.io/gorm"
)

type TransaksiData struct {
	DB *gorm.DB
}

// ini jika return nya kursusData methodnya ngga dipanggil semua diinterface masih bisa kelemahan ke akses semua kalau pakai ini data method yang dibuat sesaui yang ada di handler interface
func New(db *gorm.DB) transaksi.TransaksiDataInterface {
	return &TransaksiData{
		DB: db,
	}
}

// Transaksi
func (d *TransaksiData) CreateTransaksi(data transaksi.Transaksi) (transaksi.Transaksi, error) {
	if err := d.DB.Create(&data).Error; err != nil {
		return transaksi.Transaksi{}, err
	}

	return data, nil
}

func (d *TransaksiData) ValidateDurationKursus(userID, kursusID string) bool {
	var count int64
	err := d.DB.Model(&TransaksiHistory{}).
		Where("user_id = ?", userID).
		Where("kursus_id = ?", kursusID).
		Where("valid_until >= ?", time.Now()).
		Count(&count).Error

	if err != nil {
		return false
	}

	return count == 0 // Mengembalikan true jika tidak ada kursus yang masih aktif
}

func (d *TransaksiData) GetTotalTransaksiWithDiscount(total float64, voucherId string) (float64, error) {
	var voucherData voucher.Voucher
	if err := d.DB.Where("id = ? AND deleted_at IS NULL", voucherId).First(&voucherData).Error; err != nil {
		return 0, constant.ErrVoucherIDNotFound
	}

	discountedPrice := total - (total * voucherData.Discount / 100)
	if discountedPrice < 0 {
		discountedPrice = 0
	}

	return discountedPrice, nil
}

func (d *TransaksiData) GetAllStatusTransaksi() ([]transaksi.Transaksi, error) {
	var data []transaksi.Transaksi
	if err := d.DB.Preload("User").Preload("Voucher").Preload("Kursus").Where("deleted_at IS NULL").Find(&data).Error; err != nil {
		return nil, constant.ErrGetData
	}
	return data, nil
}

func (d *TransaksiData) GetStatusTransaksiForUser(userID string, page int, limit int) ([]transaksi.Transaksi, int, error) {
	// Pagination
	var total int64
	count := d.DB.Model(&transaksi.Transaksi{}).Where("user_id = ? AND deleted_at IS NULL", userID).Count(&total)
	if count.Error != nil {
		return nil, 0, constant.ErrTransaksiNotFound
	}
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	// data user
	var data []transaksi.Transaksi
	if err := d.DB.Preload("User").Preload("Kursus").Where("user_id = ? AND deleted_at IS NULL", userID).Offset((page - 1) * limit).Limit(limit).Find(&data).Error; err != nil {
		return nil, 0, constant.ErrGetData
	}
	return data, totalPages, nil
}

func (d *TransaksiData) GetStatusTransaksiByID(id string) (transaksi.Transaksi, error) {
	var data transaksi.Transaksi
	if err := d.DB.Preload("User").Preload("Kursus").Where("id = ? AND deleted_at IS NULL", id).Find(&data).Error; err != nil {
		return transaksi.Transaksi{}, constant.ErrGetID
	}
	return data, nil
}

// History Transaksi
func (d *TransaksiData) GetAllTransaksiHistory() ([]transaksi.TransaksiHistory, error) {
	var data []transaksi.TransaksiHistory
	if err := d.DB.Preload("Transaksi").Preload("Kursus").Where("deleted_at IS NULL").Find(&data).Error; err != nil {
		return nil, constant.ErrGetData
	}
	return data, nil
}
func (d *TransaksiData) GetAllTransaksiHistoryForUser(userID string, page, limit int) ([]transaksi.TransaksiHistory, int, error) {
	// Pagination
	var total int64
	count := d.DB.Model(&transaksi.TransaksiHistory{}).Where("user_id = ? AND deleted_at IS NULL", userID).Count(&total)
	if count.Error != nil {
		return nil, 0, constant.ErrTransaksiNotFound
	}
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	// data user
	var data []transaksi.TransaksiHistory
	if err := d.DB.Preload("Transaksi").Preload("Kursus").Where("user_id = ? AND deleted_at IS NULL", userID).Offset((page - 1) * limit).Limit(limit).Find(&data).Error; err != nil {
		return nil, 0, constant.ErrGetData
	}
	return data, totalPages, nil
}

func (d *TransaksiData) GetTransaksiHistoryByID(id string) (transaksi.TransaksiHistory, error) {
	var data transaksi.TransaksiHistory
	if err := d.DB.Preload("User").Preload("Transaksi").Preload("Kursus").Where("id = ? AND deleted_at IS NULL", id).Find(&data).Error; err != nil {
		return transaksi.TransaksiHistory{}, constant.ErrGetID
	}
	return data, nil
}

// Pagination
func (d *TransaksiData) GetAllTransaksiPagination(page, limit int) ([]transaksi.Transaksi, int, error) {
	var data []transaksi.Transaksi
	var total int64

	count := d.DB.Model(&transaksi.Transaksi{}).Where("deleted_at IS NULL").Count(&total)
	if count.Error != nil {
		return nil, 0, constant.ErrDataNotfound
	}
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	tx := d.DB.Preload("User", "deleted_at IS NULL").Preload("Kursus", "deleted_at IS NULL").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&data)

	if tx.Error != nil {
		return nil, 0, constant.ErrGetData
	}
	return data, totalPages, nil
}

func (d *TransaksiData) GetAllHistoryTransaksiPagination(page, limit int) ([]transaksi.TransaksiHistory, int, error) {
	var data []transaksi.TransaksiHistory
	var total int64

	count := d.DB.Model(&transaksi.TransaksiHistory{}).Where("deleted_at IS NULL").Count(&total)
	if count.Error != nil {
		return nil, 0, constant.ErrDataNotfound
	}
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	tx := d.DB.Preload("Transaksi", "deleted_at IS NULL").Preload("Kursus", "deleted_at IS NULL").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&data)

	if tx.Error != nil {
		return nil, 0, constant.ErrGetData
	}
	return data, totalPages, nil
}

func (d *TransaksiData) GetByIDVoucher(id string) (voucher.Voucher, error) {
	var dataVoucher voucher.Voucher
	if err := d.DB.Where("id = ?", id).Where("deleted_at IS NULL AND expired_at > ?", time.Now()).First(&dataVoucher).Error; err != nil {
		return voucher.Voucher{}, err
	}
	return dataVoucher, nil
}

func (d *TransaksiData) GetKursusByID(kursusID string) (kursus.Kursus, error) {
	var payments kursus.Kursus
	err := d.DB.Where("id = ?", kursusID).First(&payments).Error
	if err != nil {
		return kursus.Kursus{}, err
	}
	return payments, nil
}

func (d *TransaksiData) GetUserByID(userID string) (users.User, error) {
	var user users.User
	err := d.DB.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return users.User{}, err
	}
	return user, nil
}

func (d *TransaksiData) CreateTransaksiHistory(data transaksi.TransaksiHistory) error {
	if err := d.DB.Create(&data).Error; err != nil {
		return err
	}

	return nil
}

func (d *TransaksiData) ValidateUserDokumentation(userId string) bool {
	var user users.User
	err := d.DB.Where("id = ?", userId).First(&user).Error
	if err != nil {
		return false
	}

	if user.KTP == "" || user.KartuKeluarga == "" || user.Ijazah == "" {
		return false
	}

	return true
}

// Voucher

func (d *TransaksiData) UsedVoucher(data voucher.VoucherUsed) error {
	if err := d.DB.Create(&data).Error; err != nil {
		return err
	}

	return nil
}

func (d *TransaksiData) UsedVoucherCheck(userID, voucherID string) bool {
	var count int64
	if err := d.DB.Model(&voucher.VoucherUsed{}).
		Where("user_id = ? AND voucher_id = ?", userID, voucherID).
		Count(&count).Error; err != nil {
		log.Printf("Error checking voucher usage: %v", err)
		return false
	}

	return count > 0
}

func (d *TransaksiData) CheckVoucherExists(voucherID string) (bool, error) {
	var count int64
	err := d.DB.Model(&voucher.Voucher{}).Where("id = ?", voucherID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Cron Job
func (d *TransaksiData) FindExpiredTransactions(now time.Time) ([]transaksi.TransaksiHistory, error) {
	var transaksi []transaksi.TransaksiHistory
	if err := d.DB.Where("valid_until < ? AND status != ?", now, "tidak aktif").Find(&transaksi).Error; err != nil {
		return nil, err
	}

	return transaksi, nil
}

func (d *TransaksiData) UpdateStatus(transaksiID string, status string) error {
	if err := d.DB.Model(&transaksi.TransaksiHistory{}).Where("id = ?", transaksiID).Update("status", status).Error; err != nil {
		return err
	}

	return nil
}

// View User Active
func (d *TransaksiData) GetActiveUsersFromTransaksiHistory(page, limit int) ([]transaksi.TransaksiHistory, int, error) {
	var total int64
	var histories []transaksi.TransaksiHistory
	// Hitung total jumlah data dengan status "Active" tanpa memuat data relasi
	if err := d.DB.Preload("User").Where("status = ?", "Active").Find(&histories).Count(&total).Error; err != nil {
		return nil, 0, constant.ErrDataNotfound
	}

	// Ambil data dengan pagination dan preload relasi
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	// Query data dengan pagination dan preload relasi
	err := d.DB.Preload("User"). // Preload data User
					Preload("Transaksi", "deleted_at IS NULL"). // Preload data Transaksi dengan kondisi
					Preload("Kursus", "deleted_at IS NULL").    // Preload data Kursus dengan kondisi
					Where("status = ?", "Active").              // Filter berdasarkan status
					Offset((page - 1) * limit).                 // Mulai data berdasarkan halaman
					Limit(limit).                               // Batas jumlah data per halaman
					Find(&histories).Error                      // Eksekusi query dan simpan hasil ke histories
	if err != nil {
		return nil, 0, constant.ErrGetData
	}
	return histories, totalPages, nil
}
