package data

import (
	"skripsi/constant"
	"skripsi/features/voucher"
	"time"

	"gorm.io/gorm"
)

// VoucherData adalah struct yang berfungsi sebagai data layer yang menghubungkan logika aplikasi dengan database menggunakan GORM
type VoucherData struct {
	DB *gorm.DB
}

// Fungsi New memastikan bahwa struct VoucherData
// Memiliki semua metode yang sesuai dengan kontrak interface.
// Siap digunakan sebagai implementasi dari interface.
// Sebaliknya, fungsi ini mengembalikan objek VoucherData sebagai implementasi dari VoucherDataInterface.
// Atau Menciptakan objek VoucherData dan Mengembalikannya dalam bentuk interface (VoucherDataInterface) agar memenuhi kontrak yang diharapkan
func New(db *gorm.DB) voucher.VoucherDataInterface {
	return &VoucherData{
		DB: db,
	}
}

func (d *VoucherData) GetAllVoucher() ([]voucher.Voucher, error) {
	var dataVoucher []voucher.Voucher
	if err := d.DB.Where("deleted_at IS NULL deleted_at IS null AND expired_at > ?", time.Now()).First(&dataVoucher).Error; err != nil {
		return nil, err
	}
	return dataVoucher, nil
}

func (d *VoucherData) GetAllVoucherPagination(page, limit int) ([]voucher.Voucher, int, error) {
	var dataVoucher []voucher.Voucher
	var total int64
	// Count ambil data ditable sesuai where nya apa.
	// Count menghitung baris data yang ada di table. jadi misal ada 50 data. misal menampilkan data 1 halaman 10 berarti 1 dari 5 halaman.
	count := d.DB.Model(voucher.Voucher{}).Where("deleted_at IS NULL AND expired_at > ?", time.Now()).Count(&total)
	if count.Error != nil {
		return nil, 0, constant.ErrVoucherNotFound
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	tx := d.DB.Where("deleted_at IS NULL AND expired_at > ?", time.Now()).
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&dataVoucher)

	if tx.Error != nil {
		return nil, 0, constant.ErrVoucherNotFound
	}
	return dataVoucher, totalPages, nil
}

func (d *VoucherData) GetByIDVoucher(id string) (voucher.Voucher, error) {
	var dataVoucher voucher.Voucher
	if err := d.DB.Where("id = ?", id).Where("deleted_at IS NULL AND expired_at > ?", time.Now()).First(&dataVoucher).Error; err != nil {
		return voucher.Voucher{}, constant.ErrVoucherIDNotFound
	}
	return dataVoucher, nil
}

func (d *VoucherData) CreateVoucher(data voucher.Voucher) error {
	if err := d.DB.Create(&data).Error; err != nil {
		return constant.ErrVoucherFailedCreate
	}
	return nil
}

func (d *VoucherData) UpdateVoucher(data voucher.Voucher) error {
	var exisistingVoucher voucher.Voucher
	if err := d.DB.Where("id = ?", data.ID).First(&exisistingVoucher).Error; err != nil {
		return err
	}

	// ingat update satu data, updates multiple
	if err := d.DB.Model(&exisistingVoucher).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

func (d *VoucherData) DeleteVoucher(id string) error {
	res := d.DB.Begin()

	if err := res.Where("id = ?", id).Delete(&Voucher{}); err.Error != nil {
		res.Rollback()
		return constant.ErrVoucherIDNotFound
	} else if err.RowsAffected == 0 {
		res.Rollback()
		return constant.ErrVoucherIDNotFound
	}

	return res.Commit().Error
}

func (d *VoucherData) ValidateVoucher(userID string) ([]voucher.Voucher, error) {
	var dataVoucher []voucher.Voucher
	err := d.DB.Model(voucher.Voucher{}).Where("deleted_at IS NULL AND expired_at > ?", time.Now()).Not("id IN (?)", d.DB.Model(&voucher.VoucherUsed{}).
		Select("voucher_id").
		Where("user_id = ?", userID)).Find(&dataVoucher).Error
	if err != nil {
		return nil, err
	}

	return dataVoucher, nil
}
