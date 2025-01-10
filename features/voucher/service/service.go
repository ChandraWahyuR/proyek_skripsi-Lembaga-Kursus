package service

import (
	"errors"
	"skripsi/constant"
	"skripsi/features/voucher"
	"skripsi/helper"
	"strings"
	"time"
)

type VoucherService struct {
	d voucher.VoucherDataInterface
	j helper.JWTInterface
}

func New(u voucher.VoucherDataInterface, j helper.JWTInterface) voucher.VoucherServiceInterface {
	return &VoucherService{
		d: u,
		j: j,
	}
}

func (s *VoucherService) GetAllVoucher() ([]voucher.Voucher, error) {
	return s.d.GetAllVoucher()
}

func (s *VoucherService) ValidateVoucher(userId string) ([]voucher.Voucher, error) {
	return s.d.ValidateVoucher(userId)
}

func (s *VoucherService) GetAllVoucherPagination(page, limit int) ([]voucher.Voucher, int, error) {
	return s.d.GetAllVoucherPagination(page, limit)
}

func (s *VoucherService) GetByIDVoucher(id string) (voucher.Voucher, error) {
	if id == "" {
		return voucher.Voucher{}, constant.ErrVoucherIDNotFound
	}
	return s.d.GetByIDVoucher(id)
}

func (s *VoucherService) CreateVoucher(data voucher.Voucher) error {
	data.Code = strings.ToUpper(data.Code)
	switch {
	case data.Nama == "":
		return constant.ErrNameVoucher
	case data.Deskripsi == "":
		return constant.ErrDekripsiVoucher
	case data.Discount == 0:
		return constant.ErrDiscountVoucher
	case data.ExpiredAt == (time.Time{}):
		return constant.ErrExpriedAtVoucher
	}
	if data.Code != "" {
		if len(data.Code) != 10 {
			return errors.New("voucher code must be exactly 10 characters")
		}
	} else {
		data.Code = helper.GenerateCode()
	}
	return s.d.CreateVoucher(data)
}

func (s *VoucherService) UpdateVoucher(data voucher.Voucher) error {
	if data.ID == "" {
		return constant.ErrEmptyId
	}

	if data.Nama == "" && data.Deskripsi == "" && data.Code == "" && data.Discount == 0 && data.ExpiredAt == (time.Time{}) {
		return constant.ErrUpdate
	}
	data.Code = strings.ToUpper(data.Code)

	return s.d.UpdateVoucher(data)
}

func (s *VoucherService) DeleteVoucher(id string) error {
	if id == "" {
		return constant.ErrEmptyId
	}

	return s.d.DeleteVoucher(id)
}
