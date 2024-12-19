package service

// import (
// 	"errors"
// 	"skripsi/constant"
// 	"skripsi/features/kursus"
// 	"skripsi/features/transaksi"
// 	"skripsi/features/voucher"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// 	"github.com/veritrans/go-midtrans"
// )

// type MockTransaksiData struct {
// 	mock.Mock
// }

// func (m *MockTransaksiData) ValidateUserDokumentation(userID string) bool {
// 	args := m.Called(userID)
// 	return args.Bool(0)
// }

// func (m *MockTransaksiData) GetKursusByID(kursusID string) (kursus.Kursus, error) {
// 	args := m.Called(kursusID)
// 	return args.Get(0).(kursus.Kursus), args.Error(1)
// }

// func (m *MockTransaksiData) GetByIDVoucher(voucherID string) (voucher.Voucher, error) {
// 	args := m.Called(voucherID)
// 	return args.Get(0).(voucher.Voucher), args.Error(1)
// }

// func (m *MockTransaksiData) UsedVoucherCheck(userID, voucherID string) bool {
// 	args := m.Called(userID, voucherID)
// 	return args.Bool(0)
// }

// func (m *MockTransaksiData) CreateTransaksi(data transaksi.Transaksi) (transaksi.Transaksi, error) {
// 	args := m.Called(data)
// 	return args.Get(0).(transaksi.Transaksi), args.Error(1)
// }

// // Tambahkan fungsi mock lain jika diperlukan.
// type MockMidtransClient struct{}

// func (m *MockMidtransClient) GetToken(req *midtrans.SnapReq) (*midtrans.SnapResponse, error) {
// 	// Simulasi respons dari Midtrans Snap API
// 	if req.TransactionDetails.OrderID == "error_case" {
// 		return nil, errors.New("midtrans error")
// 	}
// 	return &midtrans.SnapResponse{RedirectURL: "https://mock.snap.url"}, nil
// }

// func TestCreateTransaksi(t *testing.T) {
// 	mockData := new(MockTransaksiData)
// 	mockMidtrans := &MockMidtransClient{}

// 	service := New(mockData, nil, mockMidtrans)

// 	t.Run("success - valid transaction with discount", func(t *testing.T) {
// 		transaksiInput := transaksi.Transaksi{
// 			UserID:    "user123",
// 			KursusID:  "kursus123",
// 			VoucherID: "voucher123",
// 		}
// 		kursusData := kursus.Kursus{
// 			ID:    "kursus123",
// 			Harga: 100000,
// 		}
// 		voucherData := voucher.Voucher{
// 			ID:       "voucher123",
// 			Discount: 50, // 50% discount
// 		}
// 		expectedTransaction := transaksi.Transaksi{
// 			ID:         "trans123",
// 			UserID:     "user123",
// 			KursusID:   "kursus123",
// 			TotalHarga: 50000, // 50% of 100000
// 			Status:     "Pending",
// 			SnapURL:    "https://mock.snap.url",
// 		}

// 		mockData.On("ValidateUserDokumentation", transaksiInput.UserID).Return(true)
// 		mockData.On("GetKursusByID", transaksiInput.KursusID).Return(kursusData, nil)
// 		mockData.On("GetByIDVoucher", transaksiInput.VoucherID).Return(voucherData, nil)
// 		mockData.On("UsedVoucherCheck", transaksiInput.UserID, transaksiInput.VoucherID).Return(false)
// 		mockData.On("CreateTransaksi", mock.Anything).Return(expectedTransaction, nil)

// 		result, err := service.CreateTransaksi(transaksiInput)

// 		assert.Nil(t, err)
// 		assert.Equal(t, expectedTransaction.TotalHarga, result.TotalHarga)
// 		assert.Equal(t, expectedTransaction.SnapURL, result.SnapURL)

// 		mockData.AssertExpectations(t)
// 	})

// 	t.Run("fail - user documentation invalid", func(t *testing.T) {
// 		transaksiInput := transaksi.Transaksi{
// 			UserID:   "user123",
// 			KursusID: "kursus123",
// 		}

// 		mockData.On("ValidateUserDokumentation", transaksiInput.UserID).Return(false)

// 		_, err := service.CreateTransaksi(transaksiInput)

// 		assert.NotNil(t, err)
// 		assert.Equal(t, constant.ErrValidateDokumenUser, err)

// 		mockData.AssertExpectations(t)
// 	})

// 	t.Run("fail - voucher already used", func(t *testing.T) {
// 		transaksiInput := transaksi.Transaksi{
// 			UserID:    "user123",
// 			KursusID:  "kursus123",
// 			VoucherID: "voucher123",
// 		}
// 		kursusData := kursus.Kursus{
// 			ID:    "kursus123",
// 			Harga: 100000,
// 		}

// 		mockData.On("ValidateUserDokumentation", transaksiInput.UserID).Return(true)
// 		mockData.On("GetKursusByID", transaksiInput.KursusID).Return(kursusData, nil)
// 		mockData.On("UsedVoucherCheck", transaksiInput.UserID, transaksiInput.VoucherID).Return(true)

// 		_, err := service.CreateTransaksi(transaksiInput)

// 		assert.NotNil(t, err)
// 		assert.Equal(t, constant.ErrVoucherUsed, err)

// 		mockData.AssertExpectations(t)
// 	})

// 	t.Run("fail - midtrans payment error", func(t *testing.T) {
// 		transaksiInput := transaksi.Transaksi{
// 			UserID:    "user123",
// 			KursusID:  "kursus123",
// 			VoucherID: "voucher123",
// 		}
// 		kursusData := kursus.Kursus{
// 			ID:    "kursus123",
// 			Harga: 100000,
// 		}

// 		mockData.On("ValidateUserDokumentation", transaksiInput.UserID).Return(true)
// 		mockData.On("GetKursusByID", transaksiInput.KursusID).Return(kursusData, nil)
// 		mockData.On("UsedVoucherCheck", transaksiInput.UserID, transaksiInput.VoucherID).Return(false)

// 		_, err := service.CreateTransaksi(transaksi.Transaksi{
// 			ID: "error_case",
// 		})

// 		assert.NotNil(t, err)
// 		assert.Contains(t, err.Error(), "failed to create midtrans payment")

// 		mockData.AssertExpectations(t)
// 	})
// }
