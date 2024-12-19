package service

import (
	"bytes"
	"context"
	"errors"
	"skripsi/constant"
	"skripsi/features/admin"
	"skripsi/helper"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockAdminData struct {
	mock.Mock
}

type MockJWT struct {
	mock.Mock
}

func (m *MockAdminData) RegisterAdmin(user admin.Admin) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockAdminData) LoginAdmin(user admin.Admin) (admin.Admin, error) {
	args := m.Called(user)
	return args.Get(0).(admin.Admin), args.Error(1)
}

func (m *MockAdminData) DownloadLaporanPembelian(startDate, endDate time.Time) ([]map[string]interface{}, error) {
	args := m.Called(startDate, endDate)
	return args.Get(0).([]map[string]interface{}), args.Error(1)
}

// Validator
func (m *MockAdminData) IsEmailExist(email string) bool {
	args := m.Called(email)
	return args.Bool(0)
}

func (m *MockAdminData) IsUsernameExist(username string) bool {
	args := m.Called(username)
	return args.Bool(0)
}

func (m *MockAdminData) GetByEmail(email string) (admin.Admin, error) {
	args := m.Called(email)
	return args.Get(0).(admin.Admin), args.Error(1)
}

func (m *MockAdminData) VerifyEmail(email string, isValid bool) error {
	args := m.Called(email, isValid)
	return args.Error(0)
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

func TestAdminService(t *testing.T) {
	mockData := new(MockAdminData)
	mockJWT := new(MockJWT)
	service := New(mockData, mockJWT)

	t.Run("RegisterAdmin success", func(t *testing.T) {
		password := "Password123!"
		mockAdmin := admin.Admin{
			Email:           "admin@example.com",
			Username:        "adminuser",
			Password:        password,
			ConfirmPassword: password,
		}

		mockData.On("RegisterAdmin", mock.MatchedBy(func(admin admin.Admin) bool {
			err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password))
			return err == nil && admin.Username == mockAdmin.Username && admin.Email == mockAdmin.Email
		})).Return(nil).Once()

		err := service.RegisterAdmin(mockAdmin)

		assert.Nil(t, err)
		mockData.AssertExpectations(t)
	})

	t.Run("RegisterAdmin failure - empty email", func(t *testing.T) {
		mockAdmin := admin.Admin{
			Email:           "",
			Username:        "adminuser",
			Password:        "password123",
			ConfirmPassword: "password123",
		}
		err := service.RegisterAdmin(mockAdmin)
		assert.NotNil(t, err)
		assert.Equal(t, constant.ErrEmptyEmailRegister, err)
	})

	t.Run("RegisterAdmin failure - empty username", func(t *testing.T) {
		mockAdmin := admin.Admin{
			Email:           "admin@gmail.com",
			Username:        "",
			Password:        "password123",
			ConfirmPassword: "password123",
		}
		err := service.RegisterAdmin(mockAdmin)
		assert.NotNil(t, err)
		assert.Equal(t, constant.ErrEmptyNameRegister, err)
	})

	t.Run("RegisterAdmin failure - empty password", func(t *testing.T) {
		mockAdmin := admin.Admin{
			Email:           "admin@gmail.com",
			Username:        "avc123",
			Password:        "",
			ConfirmPassword: "",
		}
		err := service.RegisterAdmin(mockAdmin)
		assert.NotNil(t, err)
		assert.Equal(t, constant.ErrEmptyPasswordRegister, err)
	})

	t.Run("RegisterAdmin failure - not match password", func(t *testing.T) {
		password := "Password123!"
		mockAdmin := admin.Admin{
			Email:           "admin@example.com",
			Username:        "adminuser",
			Password:        password,
			ConfirmPassword: "password",
		}

		mockData.On("RegisterAdmin", mock.MatchedBy(func(admin admin.Admin) bool {
			err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password))
			return err == nil && admin.Username == mockAdmin.Username && admin.Email == mockAdmin.Email
		})).Return(nil).Once()

		err := service.RegisterAdmin(mockAdmin)
		assert.NotNil(t, err)
		assert.Equal(t, constant.ErrPasswordNotMatch, err)
	})

	t.Run("GenerateLaporanCSV success", func(t *testing.T) {
		histories := []map[string]interface{}{
			{
				"id":               "1",
				"transaksi_id":     "T123",
				"kursus_id":        "K123",
				"user_id":          "U123",
				"user_nama":        "John Doe",
				"email":            "john@example.com",
				"nama_kursus":      "Kursus A",
				"status":           "Active",
				"valid_until":      time.Now(),
				"total_harga":      100000.0,
				"transaksi_status": "Success",
			},
		}

		startDate := time.Now().AddDate(0, 0, -7)
		endDate := time.Now()

		buffer := &bytes.Buffer{}
		err := service.GenerateLaporanCSV(buffer, histories, startDate, endDate)

		assert.Nil(t, err)
		csvContent := buffer.String()
		assert.Contains(t, csvContent, "John Doe")
		assert.Contains(t, csvContent, "100000.00")
	})
}

func TestLoginAdmin(t *testing.T) {
	mockUserData := new(MockAdminData)
	mockJWT := new(MockJWT)
	service := New(mockUserData, mockJWT)

	t.Run("success login", func(t *testing.T) {
		mockUser := admin.Admin{
			Username: "abc123",
			Password: "password123",
		}

		mockUserDataResponse := admin.Admin{
			ID:       "123",
			Username: "abc123",
		}

		mockJWTResponse := "mockToken123"

		mockUserData.On("LoginAdmin", mockUser).Return(mockUserDataResponse, nil).Once()
		mockJWT.On("GenerateAdminJWT", mock.Anything).Return(mockJWTResponse, nil).Once()

		result, err := service.LoginAdmin(mockUser)

		assert.Nil(t, err)
		assert.Equal(t, mockJWTResponse, result.Token)
		mockUserData.AssertExpectations(t)
		mockJWT.AssertExpectations(t)
	})

	t.Run("fail - username empty", func(t *testing.T) {
		mockUser := admin.Admin{
			Username: "",
			Password: "password123",
		}

		_, err := service.LoginAdmin(mockUser)
		assert.NotNil(t, err)
		assert.Equal(t, constant.ErrAdminUserNameEmpty, err)
	})

	t.Run("fail - password empty", func(t *testing.T) {
		mockUser := admin.Admin{
			Username: "admin123",
			Password: "",
		}

		_, err := service.LoginAdmin(mockUser)
		assert.NotNil(t, err)
		assert.Equal(t, constant.ErrAdminPasswordEmpty, err)
	})

	t.Run("fail - data not found", func(t *testing.T) {
		mockUser := admin.Admin{
			Username: "abc123",
			Password: "password123",
		}

		mockUserData.On("LoginAdmin", mockUser).Return(admin.Admin{}, errors.New("data not found")).Once()

		_, err := service.LoginAdmin(mockUser)

		assert.NotNil(t, err)
		assert.Equal(t, "data not found", err.Error())

		mockUserData.AssertExpectations(t)
		mockJWT.AssertExpectations(t)
	})

}
