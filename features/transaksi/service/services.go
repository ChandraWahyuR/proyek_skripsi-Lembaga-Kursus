package service

import (
	"fmt"
	"log"
	"skripsi/constant"
	"skripsi/features/transaksi"
	"skripsi/helper"

	"github.com/veritrans/go-midtrans"
)

type TransaksiService struct {
	d              transaksi.TransaksiDataInterface
	j              helper.JWTInterface
	midtransClient midtrans.Client
}

func New(u transaksi.TransaksiDataInterface, j helper.JWTInterface, midtransClient midtrans.Client) transaksi.TransaksiServiceInterface {
	return &TransaksiService{
		d:              u,
		j:              j,
		midtransClient: midtransClient,
	}
}

func (s *TransaksiService) CreateTransaksi(transaksiData transaksi.Transaksi) (transaksi.Transaksi, error) {
	// Validator
	if !s.d.ValidateUserDokumentation(transaksiData.UserID) {
		return transaksi.Transaksi{}, constant.ErrValidateDokumenUser
	}

	var finalPrice float64
	kursusData, err := s.d.GetKursusByID(transaksiData.KursusID)
	if err != nil {
		// log.Printf("Error fetching Kursus for ID: %s, Error: %v", transaksiData.KursusID, err)
		return transaksi.Transaksi{}, fmt.Errorf("failed to get Kursus data: %v", err)
	}

	finalPrice = float64(kursusData.Harga)
	// log.Printf("Kursus Price: %.2f", finalPrice)

	// Kalkulasi diskon
	if transaksiData.VoucherID != "" {
		voucher, err := s.d.GetByIDVoucher(transaksiData.VoucherID)
		if err != nil {
			// log.Printf("Error fetching voucher for ID: %s, Error: %v", transaksiData.VoucherID, err)
			return transaksi.Transaksi{}, err
		}

		discount := voucher.Discount / 100
		// log.Printf("Voucher Discount: %.2f%%", voucher.Discount)

		finalPrice -= (finalPrice * discount)
		if finalPrice < 0.01 {
			finalPrice = 0.01
		}
	}

	// log.Printf("Calculated Final Price for transaction: %.2f", finalPrice)

	transaksiData.Status = "Pending"
	transaksiData.TotalHarga = finalPrice

	snapURL, err := s.createMidtransPayment(transaksiData)
	if err != nil {
		log.Printf("Error creating Midtrans payment for transaction ID: %s, Error: %v", transaksiData.ID, err)
		return transaksi.Transaksi{}, fmt.Errorf("failed to create midtrans payment: %v", err)
	}

	// log.Printf("Midtrans Snap URL for transaction ID: %s: %s", transaksiData.ID, snapURL)

	transaksiData.SnapURL = snapURL
	savedTransaksi, err := s.d.CreateTransaksi(transaksiData)
	if err != nil {
		// log.Printf("Error saving transaction ID: %s, Error: %v", transaksiData.ID, err)
		return transaksi.Transaksi{}, err
	}

	return savedTransaksi, nil
}

func (s *TransaksiService) GetAllTransaksiPagination(page, limit int) ([]transaksi.Transaksi, int, error) {
	return s.d.GetAllTransaksiPagination(page, limit)
}

func (s *TransaksiService) CreateTransaksiHistory(data transaksi.TransaksiHistory) error {
	return s.d.CreateTransaksiHistory(data)
}

func (s *TransaksiService) GetStatusTransaksiForUser(userID string, page int, limit int) ([]transaksi.Transaksi, int, error) {
	if userID == "" {
		return nil, 0, constant.ErrUnauthorized
	}
	return s.d.GetStatusTransaksiForUser(userID, page, limit)
}

func (s *TransaksiService) GetStatusTransaksiByID(id string) (transaksi.Transaksi, error) {
	if id == "" {
		return transaksi.Transaksi{}, constant.ErrGetID
	}
	return s.d.GetStatusTransaksiByID(id)
}

func (s *TransaksiService) GetAllHistoryTransaksiPagination(page, limit int) ([]transaksi.TransaksiHistory, int, error) {
	return s.d.GetAllHistoryTransaksiPagination(page, limit)
}

func (s *TransaksiService) GetAllTransaksiHistoryForUser(userID string, page, limit int) ([]transaksi.TransaksiHistory, int, error) {
	return s.d.GetAllTransaksiHistoryForUser(userID, page, limit)
}

func (s *TransaksiService) GetTransaksiHistoryByID(id string) (transaksi.TransaksiHistory, error) {
	if id == "" {
		return transaksi.TransaksiHistory{}, constant.ErrGetID
	}
	return s.d.GetTransaksiHistoryByID(id)
}

// =============================================================================================
func (s *TransaksiService) createMidtransPayment(transaksi transaksi.Transaksi) (string, error) {
	snapGateway := midtrans.SnapGateway{
		Client: s.midtransClient,
	}
	userData, err := s.d.GetUserByID(transaksi.UserID)
	if err != nil {
		return "", err
	}

	kursusData, err := s.d.GetKursusByID(transaksi.KursusID)
	if err != nil {
		return "", err
	}
	custAddress := &midtrans.CustAddress{
		FName:       userData.Username,
		Phone:       userData.NomorHP,
		CountryCode: "IDN",
	}

	// log.Printf("Kursus Data: ID=%s, Name=%s, Price=%d", kursusData.ID, kursusData.Nama, kursusData.Harga)
	itemDetails := []midtrans.ItemDetail{
		{
			ID:    kursusData.ID,
			Name:  kursusData.Nama,
			Price: int64(transaksi.TotalHarga),
			Qty:   1,
		},
	}

	snapReq := &midtrans.SnapReq{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  transaksi.ID,
			GrossAmt: int64(transaksi.TotalHarga),
		},
		CustomerDetail: &midtrans.CustDetail{
			FName:    userData.Username,
			Email:    userData.Email,
			Phone:    userData.NomorHP,
			BillAddr: custAddress,
			ShipAddr: custAddress,
		},
		Items: &itemDetails,
	}

	// log.Printf("Midtrans Request OrderID: %s, Gross Amount: %.2f", transaksi.ID, transaksi.TotalHarga)

	snapResp, err := snapGateway.GetToken(snapReq)
	if err != nil {
		// log.Printf("Midtrans Payment Error for OrderID: %s, Error: %v", transaksi.ID, err)
		return "", err
	}

	// log.Printf("Midtrans Payment Response URL for OrderID: %s: %s", transaksi.ID, snapResp.RedirectURL)

	return snapResp.RedirectURL, nil
}
