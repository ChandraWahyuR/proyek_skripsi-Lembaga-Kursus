package service

import (
	"context"
	"errors"
	"skripsi/constant"
	"skripsi/features/voucher"
	"skripsi/helper"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Data Layer
type MockVoucherData struct {
	mock.Mock
}

func (m *MockVoucherData) GetAllVoucher() ([]voucher.Voucher, error) {
	args := m.Called()
	return args.Get(0).([]voucher.Voucher), args.Error(1)
}

func (m *MockVoucherData) GetAllVoucherPagination(page, limit int) ([]voucher.Voucher, int, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]voucher.Voucher), args.Int(1), args.Error(2)
}

func (m *MockVoucherData) GetByIDVoucher(id string) (voucher.Voucher, error) {
	args := m.Called(id)
	return args.Get(0).(voucher.Voucher), args.Error(1)
}

func (m *MockVoucherData) CreateVoucher(data voucher.Voucher) error {
	args := m.Called(data)
	return args.Error(0)
}

func (m *MockVoucherData) UpdateVoucher(data voucher.Voucher) error {
	args := m.Called(data)
	return args.Error(0)
}

func (m *MockVoucherData) DeleteVoucher(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockHelper struct {
	mock.Mock
}

func (m *MockHelper) ValidateEmail(email string) bool {
	args := m.Called(email)
	return args.Bool(0)
}

func (m *MockHelper) ValidateUsername(username string) (string, error) {
	args := m.Called(username)
	return args.String(0), args.Error(1)
}

func (m *MockHelper) ValidatePassword(password string) (string, error) {
	if len(password) < 6 {
		return "", errors.New("password too short")
	}
	return password, nil
}

func (m *MockHelper) SendEmail(to string, subject string, body string) error {
	return nil
}

func (m *MockHelper) HashPassword(password string) (string, error) {
	return "hashed-" + password, nil
}

func (m *MockHelper) TelephoneValidator(number string) (string, error) {
	if len(number) < 11 || len(number) > 12 {
		return "", errors.New("invalid phone number")
	}
	return number, nil
}

// =========================================================================================
type MockJWT struct {
	mock.Mock
}

// User JWT
func (m *MockJWT) GenerateUserToken(user helper.UserJWT) string {
	args := m.Called(user)
	return args.String(0)
}

func (m *MockJWT) GenerateUserJWT(user helper.UserJWT) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockJWT) ExtractUserToken(token *jwt.Token) map[string]interface{} {
	args := m.Called(token)
	return args.Get(0).(map[string]interface{})
}

// Admin JWT
func (m *MockJWT) GenerateAdminToken(admin helper.AdminJWT) string {
	args := m.Called(admin)
	return args.String(0)
}

func (m *MockJWT) GenerateAdminJWT(admin helper.AdminJWT) (string, error) {
	args := m.Called(admin)
	return args.String(0), args.Error(1)
}

func (m *MockJWT) ExtractAdminToken(token *jwt.Token) map[string]interface{} {
	args := m.Called(token)
	return args.Get(0).(map[string]interface{})
}

// Forgot JWT
func (m *MockJWT) GenerateForgotPassToken(user helper.ForgotPassJWT) string {
	args := m.Called(user)
	return args.String(0)
}

func (m *MockJWT) GenerateForgotPassJWT(user helper.ForgotPassJWT) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

// Valid
func (m *MockJWT) ValidateToken(ctx context.Context, token string) (*jwt.Token, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(*jwt.Token), args.Error(1)
}

// Email
func (m *MockJWT) GenerateVerifikasiEmailToken(user helper.UserJWT) string {
	args := m.Called(user)
	return args.String(0)
}
func (m *MockJWT) GenerateVerifikasiEmailJWT(user helper.UserJWT) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockJWT) ValidateEmailToken(tokenString string) (*jwt.Token, error) {
	args := m.Called(tokenString)
	return args.Get(0).(*jwt.Token), args.Error(1)
}

func TestVoucherService(t *testing.T) {
	mockData := new(MockVoucherData)
	mockJWT := new(MockJWT)
	service := New(mockData, mockJWT)

	t.Run("GetAllVoucher success", func(t *testing.T) {
		mockResponse := []voucher.Voucher{{ID: "1", Nama: "Voucher A"}}
		mockData.On("GetAllVoucher").Return(mockResponse, nil).Once()

		result, err := service.GetAllVoucher()

		assert.Nil(t, err)
		assert.Equal(t, mockResponse, result)
		mockData.AssertExpectations(t)
	})

	t.Run("GetAllVoucherPagination success", func(t *testing.T) {
		mockResponse := []voucher.Voucher{{ID: "1", Nama: "Voucher A"}}
		mockTotal := 10
		mockData.On("GetAllVoucherPagination", 1, 5).Return(mockResponse, mockTotal, nil).Once()

		result, total, err := service.GetAllVoucherPagination(1, 5)

		assert.Nil(t, err)
		assert.Equal(t, mockResponse, result)
		assert.Equal(t, mockTotal, total)
		mockData.AssertExpectations(t)
	})

	t.Run("GetByIDVoucher success", func(t *testing.T) {
		mockResponse := voucher.Voucher{ID: "1", Nama: "Voucher A"}
		mockData.On("GetByIDVoucher", "1").Return(mockResponse, nil).Once()

		result, err := service.GetByIDVoucher("1")

		assert.Nil(t, err)
		assert.Equal(t, mockResponse, result)
		mockData.AssertExpectations(t)
	})

	t.Run("GetByIDVoucher fail - empty ID", func(t *testing.T) {
		_, err := service.GetByIDVoucher("")

		assert.NotNil(t, err)
		assert.Equal(t, constant.ErrVoucherIDNotFound, err)
	})

	t.Run("CreateVoucher success", func(t *testing.T) {
		mockVoucher := voucher.Voucher{Nama: "Voucher A", Deskripsi: "Description", Discount: 10, ExpiredAt: time.Now(), Code: "TEST123456"}
		mockData.On("CreateVoucher", mockVoucher).Return(nil).Once()

		err := service.CreateVoucher(mockVoucher)

		assert.Nil(t, err)
		mockData.AssertExpectations(t)
	})

	t.Run("CreateVoucher fail - name empty", func(t *testing.T) {
		mockVoucher := voucher.Voucher{Nama: "", Deskripsi: "", Discount: 0, ExpiredAt: time.Time{}, Code: ""}
		err := service.CreateVoucher(mockVoucher)

		assert.NotNil(t, err)
		assert.Equal(t, constant.ErrNameVoucher, err)
	})

	t.Run("CreateVoucher fail - deskripsi empty", func(t *testing.T) {
		mockVoucher := voucher.Voucher{Nama: "Kursus 1", Deskripsi: "", Discount: 0, ExpiredAt: time.Time{}, Code: ""}
		err := service.CreateVoucher(mockVoucher)

		assert.NotNil(t, err)
		assert.Equal(t, constant.ErrDekripsiVoucher, err)
	})

	t.Run("CreateVoucher fail - discount empty", func(t *testing.T) {
		mockVoucher := voucher.Voucher{Nama: "Kursus 1", Deskripsi: "apa coba", Discount: 0, ExpiredAt: time.Time{}, Code: ""}
		err := service.CreateVoucher(mockVoucher)

		assert.NotNil(t, err)
		assert.Equal(t, constant.ErrDiscountVoucher, err)
	})
	t.Run("CreateVoucher fail - expired at empty", func(t *testing.T) {
		mockVoucher := voucher.Voucher{Nama: "Kursus 1", Deskripsi: "apa coba", Discount: 10, ExpiredAt: time.Time{}, Code: ""}
		err := service.CreateVoucher(mockVoucher)

		assert.NotNil(t, err)
		assert.Equal(t, constant.ErrExpriedAtVoucher, err)
	})

	t.Run("CreateVoucher fail - code not valid", func(t *testing.T) {
		validJamSelesai := time.Date(2024, time.December, 18, 13, 0, 0, 0, time.UTC)
		mockVoucher := voucher.Voucher{Nama: "Kursus 1", Deskripsi: "apa coba", Discount: 10, ExpiredAt: validJamSelesai, Code: "1"}
		err := service.CreateVoucher(mockVoucher)

		assert.NotNil(t, err)
		assert.Equal(t, errors.New("voucher code must be exactly 10 characters"), err)
	})

	t.Run("UpdateVoucher success", func(t *testing.T) {
		mockVoucher := voucher.Voucher{ID: "1", Nama: "Updated Voucher"}
		mockData.On("UpdateVoucher", mockVoucher).Return(nil).Once()

		err := service.UpdateVoucher(mockVoucher)

		assert.Nil(t, err)
		mockData.AssertExpectations(t)
	})

	t.Run("UpdateVoucher fail - empty ID", func(t *testing.T) {
		mockVoucher := voucher.Voucher{}
		err := service.UpdateVoucher(mockVoucher)

		assert.NotNil(t, err)
		assert.Equal(t, constant.ErrEmptyId, err)
	})

	t.Run("DeleteVoucher success", func(t *testing.T) {
		mockData.On("DeleteVoucher", "1").Return(nil).Once()

		err := service.DeleteVoucher("1")

		assert.Nil(t, err)
		mockData.AssertExpectations(t)
	})

	t.Run("DeleteVoucher fail - empty ID", func(t *testing.T) {
		err := service.DeleteVoucher("")

		assert.NotNil(t, err)
		assert.Equal(t, constant.ErrEmptyId, err)
	})
}
