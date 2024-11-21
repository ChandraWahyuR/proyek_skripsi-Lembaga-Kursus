package data

import (
	"skripsi/constant"
	"skripsi/features/voucher"
	"time"

	"gorm.io/gorm"
)

type VoucherData struct {
	DB *gorm.DB
}

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
	// Count ambil data ditable sesuai where nya apa
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
