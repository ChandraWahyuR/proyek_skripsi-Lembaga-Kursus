package service

import (
	"context"
	"errors"
	"fmt"
	"skripsi/constant"
	"skripsi/features/users"
	"skripsi/helper"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// go test ./... -v
// go test ./... -v -cover
// go test ./... -coverprofile=coverage.out
// go tool cover -html=coverage.out

type MockUserData struct {
	mock.Mock
}

func (m *MockUserData) Register(user users.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserData) Login(user users.User) (users.User, error) {
	args := m.Called(user)
	return args.Get(0).(users.User), args.Error(1)
}

// Forgot password
func (m *MockUserData) ForgotPassword(forgot users.ForgotPassword) error {
	args := m.Called(forgot)
	return args.Error(0)
}

func (m *MockUserData) VerifyOTP(verifikasi users.VerifyOtp) error {
	args := m.Called(verifikasi)
	return args.Error(0)
}

func (m *MockUserData) ResetPassword(reset users.ResetPassword) error {
	args := m.Called(reset)
	return args.Error(0)
}

// Validator
func (m *MockUserData) IsEmailExist(email string) bool {
	args := m.Called(email)
	return args.Bool(0)
}

func (m *MockUserData) IsUsernameExist(username string) bool {
	args := m.Called(username)
	return args.Bool(0)
}

func (m *MockUserData) GetByEmail(email string) (users.User, error) {
	args := m.Called(email)
	return args.Get(0).(users.User), args.Error(1)
}

func (m *MockUserData) VerifyEmail(email string, isValid bool) error {
	args := m.Called(email, isValid)
	return args.Error(0)
}

// Pagination
func (m *MockUserData) GetAllUserPagination(page, limit int) ([]users.User, int, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]users.User), args.Int(1), args.Error(2)
}

func (m *MockUserData) GetUserByID(id string) (users.User, error) {
	args := m.Called(id)
	return args.Get(0).(users.User), args.Error(1)
}

func (m *MockUserData) UpdateUser(id users.EditUser) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserData) DeleteUser(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *MockUserData) SearchUserByUsernameEmail(search string, page, limit int) ([]users.User, int, error) {
	args := m.Called(search, page, limit)
	return args.Get(0).([]users.User), args.Int(1), args.Error(2)
}

// Mock Helper
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

type MockEmail struct {
	mock.Mock
}

func (m *MockEmail) SendEmail(to, subject, body string) error {
	args := m.Called(to, subject, body)
	return args.Error(0)
}

// =============================================================================

func TestRegister(t *testing.T) {
	mockUserData := new(MockUserData)
	mockHelper := new(MockHelper)
	service := New(mockUserData, nil, mockHelper)

	t.Run("success register", func(t *testing.T) {
		password := "Password123!"
		mockUser := users.User{
			Username:        "validuser",
			Email:           "user@example.com",
			Password:        password,
			ConfirmPassword: password,
			NomorHP:         "081234567890",
		}

		mockUserData.On("Register", mock.MatchedBy(func(user users.User) bool {
			return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil
		})).Return(nil).Once()

		err := service.Register(mockUser)
		assert.Nil(t, err)
		mockUserData.AssertExpectations(t)
	})

	t.Run("fail because empty username", func(t *testing.T) {
		mockUser := users.User{
			Username:        "",
			Email:           "user@example.com",
			Password:        "password123",
			ConfirmPassword: "password123",
			NomorHP:         "081234567890",
		}

		err := service.Register(mockUser)
		assert.NotNil(t, err)
		assert.Equal(t, constant.ErrEmptyNameRegister, err)
	})

	t.Run("fail because empty email", func(t *testing.T) {
		mockUser := users.User{
			Username:        "validuser",
			Email:           "",
			Password:        "password123",
			ConfirmPassword: "password123",
			NomorHP:         "081234567890",
		}

		err := service.Register(mockUser)
		assert.NotNil(t, err)
		assert.Equal(t, constant.ErrEmptyEmailRegister, err)
	})

	t.Run("fail because confirm password email", func(t *testing.T) {
		mockUser := users.User{
			Username:        "validuser",
			Email:           "user@example.com",
			Password:        "Password123!",
			ConfirmPassword: "",
			NomorHP:         "081234567890",
		}

		err := service.Register(mockUser)
		assert.NotNil(t, err)
		assert.Equal(t, errors.New("confirm password cannot be empty"), err)
	})

	t.Run("fail because email invalid", func(t *testing.T) {
		mockUser := users.User{
			Username:        "validuser",
			Email:           "invalidemail",
			Password:        "password123",
			ConfirmPassword: "password123",
			NomorHP:         "081234567890",
		}

		err := service.Register(mockUser)
		assert.NotNil(t, err)
		assert.Equal(t, constant.ErrInvalidEmail, err)
	})

	t.Run("fail because password mismatch", func(t *testing.T) {
		mockUser := users.User{
			Username:        "validuser",
			Email:           "user@example.com",
			Password:        "password123",
			ConfirmPassword: "wrongpassword",
			NomorHP:         "081234567890",
		}

		err := service.Register(mockUser)
		assert.NotNil(t, err)
		assert.Equal(t, constant.ErrPasswordNotMatch, err)
	})

	t.Run("fail because NomorHP empty", func(t *testing.T) {
		mockUser := users.User{
			Username:        "validuser",
			Email:           "user@example.com",
			Password:        "Password123!",
			ConfirmPassword: "Password123!",
		}

		err := service.Register(mockUser)
		assert.NotNil(t, err)
		assert.Equal(t, errors.New("telephone numbber cannot be empty and must be 11 or 12 digit"), err)
	})

	t.Run("fail because username not valid", func(t *testing.T) {
		mockUser := users.User{
			Username:        "a1",
			Email:           "user@example.com",
			Password:        "Password123!",
			ConfirmPassword: "Password123!",
			NomorHP:         "012345678912",
		}

		err := service.Register(mockUser)
		assert.NotNil(t, err)
		assert.Equal(t, constant.ErrInvalidUsername, err)
	})

	t.Run("fail because NomorHP not valid", func(t *testing.T) {
		mockUser := users.User{
			Username:        "validuser",
			Email:           "user@example.com",
			Password:        "Password123!",
			ConfirmPassword: "Password123!",
			NomorHP:         "01234",
		}

		err := service.Register(mockUser)
		assert.NotNil(t, err)
		assert.Equal(t, errors.New("invalid phone number format"), err)
	})
}

func TestLogin(t *testing.T) {
	mockUserData := new(MockUserData)
	mockJWT := new(MockJWT)
	mockHelper := new(MockHelper)
	service := New(mockUserData, mockJWT, mockHelper)

	t.Run("success login", func(t *testing.T) {
		mockUser := users.User{
			Email:    "user@example.com",
			Password: "password123",
		}

		mockUserDataResponse := users.User{
			ID:    "123",
			Email: "user@example.com",
		}

		mockJWTResponse := "mockToken123"

		// Mocking UserDataInterface.Login
		mockUserData.On("Login", mockUser).Return(mockUserDataResponse, nil).Once()

		// Mocking JWTInterface.GenerateUserJWT
		mockJWT.On("GenerateUserJWT", mock.Anything).Return(mockJWTResponse, nil).Once()

		result, err := service.Login(mockUser)

		assert.Nil(t, err)
		assert.Equal(t, mockJWTResponse, result.Token)
		mockUserData.AssertExpectations(t)
		mockJWT.AssertExpectations(t)
	})

	t.Run("fail login - empty email or password", func(t *testing.T) {
		mockUser := users.User{
			Email:    "",
			Password: "",
		}

		result, err := service.Login(mockUser)

		assert.NotNil(t, err)
		assert.Equal(t, constant.ErrEmptyLogin, err)
		assert.Equal(t, users.Login{}, result)
	})

	t.Run("fail login - invalid email", func(t *testing.T) {
		mockUser := users.User{
			Email:    "invalid-email",
			Password: "password123",
		}

		result, err := service.Login(mockUser)

		assert.NotNil(t, err)
		assert.Equal(t, constant.ErrInvalidEmail, err)
		assert.Equal(t, users.Login{}, result)
	})

	t.Run("fail login - user not found", func(t *testing.T) {
		mockUser := users.User{
			Email:    "user@example.com",
			Password: "password123",
		}

		mockUserData.On("Login", mockUser).Return(users.User{}, errors.New("user not found")).Once()

		result, err := service.Login(mockUser)

		assert.NotNil(t, err)
		assert.Equal(t, "user not found", err.Error())
		assert.Equal(t, users.Login{}, result)
		mockUserData.AssertExpectations(t)
	})

	t.Run("fail login - JWT generation error", func(t *testing.T) {
		mockUser := users.User{
			Email:    "user@example.com",
			Password: "password123",
		}

		mockUserDataResponse := users.User{
			ID:    "123",
			Email: "user@example.com",
		}

		mockUserData.On("Login", mockUser).Return(mockUserDataResponse, nil).Once()
		mockJWT.On("GenerateUserJWT", mock.Anything).Return("", errors.New("failed to generate token")).Once()

		result, err := service.Login(mockUser)

		assert.NotNil(t, err)
		assert.Equal(t, "failed to generate token", err.Error())
		assert.Equal(t, users.Login{}, result)
		mockUserData.AssertExpectations(t)
		mockJWT.AssertExpectations(t)
	})
}

func TestForgotPassword(t *testing.T) {
	mockUserData := new(MockUserData)
	mockMailer := new(MockEmail)
	mockJWT := new(MockJWT)

	service := New(mockUserData, mockJWT, mockMailer)

	t.Run("success forgot password", func(t *testing.T) {
		email := "user@example.com"
		mockUser := users.User{ID: "12345", Email: email}

		mockUserData.On("ForgotPassword", mock.Anything).Return(nil).Once()
		mockMailer.On("SendEmail", email, "Reset Password", mock.AnythingOfType("string")).Return(nil).Once()
		mockUserData.On("GetByEmail", email).Return(mockUser, nil).Once()
		mockJWT.On("GenerateForgotPassJWT", mock.Anything).Return("validToken", nil).Once()

		token, err := service.ForgotPassword(users.User{Email: email})
		assert.Nil(t, err)
		assert.NotEmpty(t, token)

		mockUserData.AssertExpectations(t)
		mockMailer.AssertExpectations(t)
		mockJWT.AssertExpectations(t)
	})

	t.Run("fail because empty email", func(t *testing.T) {
		_, err := service.ForgotPassword(users.User{Email: ""})
		assert.NotNil(t, err)
		assert.Equal(t, "email cannot be empty", err.Error())
	})

	t.Run("fail because invalid email format", func(t *testing.T) {
		email := "invalidemail"
		_, err := service.ForgotPassword(users.User{Email: email})
		assert.NotNil(t, err)
		assert.Equal(t, constant.ErrInvalidEmail, err)
	})

	t.Run("fail to save OTP", func(t *testing.T) {
		email := "user@example.com"
		mockUserData.On("ForgotPassword", mock.Anything).Return(errors.New("failed to save OTP")).Once()

		_, err := service.ForgotPassword(users.User{Email: email})
		assert.NotNil(t, err)
		assert.Equal(t, "failed to save OTP", err.Error())

		mockUserData.AssertExpectations(t)
	})

	t.Run("fail to send email", func(t *testing.T) {
		email := "user@example.com"
		mockUserData.On("ForgotPassword", mock.Anything).Return(nil).Once()
		mockMailer.On("SendEmail", email, "Reset Password", mock.AnythingOfType("string")).Return(errors.New("failed to send email")).Once()

		_, err := service.ForgotPassword(users.User{Email: email})
		assert.NotNil(t, err)
		assert.Equal(t, "failed to send email", err.Error())

		mockUserData.AssertExpectations(t)
		mockMailer.AssertExpectations(t)
	})

	t.Run("fail to get user data", func(t *testing.T) {
		email := "user@example.com"
		mockUserData.On("ForgotPassword", mock.Anything).Return(nil).Once()
		mockMailer.On("SendEmail", email, "Reset Password", mock.AnythingOfType("string")).Return(nil).Once()
		mockUserData.On("GetByEmail", email).Return(users.User{}, errors.New("user not found")).Once()

		_, err := service.ForgotPassword(users.User{Email: email})
		assert.NotNil(t, err)
		assert.Equal(t, "user not found", err.Error())

		mockUserData.AssertExpectations(t)
		mockMailer.AssertExpectations(t)
	})

	t.Run("fail to generate JWT", func(t *testing.T) {
		email := "user@example.com"
		mockUser := users.User{ID: "12345", Email: email}

		mockUserData.On("ForgotPassword", mock.Anything).Return(nil).Once()
		mockMailer.On("SendEmail", email, "Reset Password", mock.AnythingOfType("string")).Return(nil).Once()
		mockUserData.On("GetByEmail", email).Return(mockUser, nil).Once()
		mockJWT.On("GenerateForgotPassJWT", mock.Anything).Return("", errors.New("failed to generate JWT")).Once()

		_, err := service.ForgotPassword(users.User{Email: email})
		assert.NotNil(t, err)
		assert.Equal(t, "failed to generate JWT", err.Error())

		mockUserData.AssertExpectations(t)
		mockMailer.AssertExpectations(t)
		mockJWT.AssertExpectations(t)
	})
}

func TestVerifyOTP(t *testing.T) {
	mockData := new(MockUserData)
	service := New(mockData, nil, nil)

	t.Run("success verify otp", func(t *testing.T) {
		input := users.VerifyOtp{
			Email: "test@example.com",
			Otp:   "12345",
		}

		mockData.On("VerifyOTP", input).Return(nil)

		// Run service method
		err := service.VerifyOTP(input)

		// Assertions
		assert.Nil(t, err)
		mockData.AssertExpectations(t)
	})
	t.Run("fail input verify otp", func(t *testing.T) {
		input := users.VerifyOtp{
			Email: "test@example.com",
			Otp:   "otp-wrong",
		}

		mockData.On("VerifyOTP", input).Return(errors.New("bad request")).Once()

		err := service.VerifyOTP(input)

		assert.NotNil(t, err)
		assert.Equal(t, constant.ErrBadRequest, err)

		mockData.AssertExpectations(t)
	})
}

func TestVerifyOTP_InvalidEmail(t *testing.T) {
	service := New(nil, nil, nil)

	input := users.VerifyOtp{
		Email: "invalid-email",
		Otp:   "12345",
	}

	// Run service method
	err := service.VerifyOTP(input)

	// Assertions
	assert.Equal(t, constant.ErrInvalidEmail, err)
}

func TestResetPassword(t *testing.T) {
	mockData := new(MockUserData)
	mockHelper := new(MockHelper)
	service := &UserService{
		d: mockData,
		j: nil,
		m: mockHelper,
	}

	t.Run("success - reset password", func(t *testing.T) {
		input := users.ResetPassword{
			Password:             "Secure@123",
			ConfirmationPassword: "Secure@123",
		}

		hashedPassword := "hashedpassword"

		mockData.On("ResetPassword", mock.Anything).Return(nil)
		helperMock := new(MockHelper)
		helperMock.On("ValidatePassword", input.Password).Return(true, nil)
		helperMock.On("HashPassword", true).Return(hashedPassword, nil)

		// Run service method
		err := service.ResetPassword(input)

		// Assertions
		assert.Nil(t, err)
		mockData.AssertExpectations(t)
	})

	t.Run("fail - reset password", func(t *testing.T) {
		input := users.ResetPassword{
			Password:             "",
			ConfirmationPassword: "",
		}

		err := service.ResetPassword(input)

		assert.NotNil(t, err)
		assert.Equal(t, errors.New("password, confirm password cannot be empty"), err)
	})

	t.Run("fail - mismatch password", func(t *testing.T) {
		input := users.ResetPassword{
			Password:             "Testing123!",
			ConfirmationPassword: "Testing123?",
		}

		err := service.ResetPassword(input)

		assert.NotNil(t, err)
		assert.Equal(t, constant.ErrPasswordNotMatch, err)
	})

}

func TestActivateAccount(t *testing.T) {
	mockData := new(MockUserData)
	service := New(mockData, nil, nil)

	email := "user@example.com"

	mockData.On("VerifyEmail", email, true).Return(nil)

	// Run service method
	err := service.ActivateAccount(email)

	// Assertions
	assert.Nil(t, err)
	mockData.AssertExpectations(t)
}

func TestSendVerificationEmail(t *testing.T) {
	mockEmail := new(MockEmail)
	service := New(nil, nil, mockEmail)

	email := "user@example.com"
	link := "http://example.com/verify"
	subject := "Verifikasi Email"
	body := fmt.Sprintf("<p>Klik email dibawah Untuk memverifikasi email:</p><a href='%s'>Verify</a>", link)

	mockEmail.On("SendEmail", email, subject, body).Return(nil)

	// Run service method
	err := service.SendVerificationEmail(email, link)

	// Assertions
	assert.Nil(t, err)
	mockEmail.AssertExpectations(t)
}

func TestGetAllUserPagination(t *testing.T) {
	mockData := new(MockUserData)
	svc := UserService{d: mockData}

	mockUsers := []users.User{
		{ID: "1", Nama: "John Doe"},
		{ID: "2", Nama: "Jane Doe"},
	}

	mockData.On("GetAllUserPagination", 1, 10).Return(mockUsers, 2, nil)

	users, total, err := svc.GetAllUserPagination(1, 10)

	assert.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Equal(t, mockUsers, users)
	mockData.AssertExpectations(t)
}

func TestGetUserByID(t *testing.T) {
	mockData := new(MockUserData)
	svc := UserService{d: mockData}

	t.Run("success - get id", func(t *testing.T) {
		mockUser := users.User{ID: "1", Nama: "John Doe"}

		// Success case
		mockData.On("GetUserByID", "1").Return(mockUser, nil)
		user, err := svc.GetUserByID("1")

		assert.NoError(t, err)
		assert.Equal(t, mockUser, user)
	})

	t.Run("fail - get id", func(t *testing.T) {
		_, err := svc.GetUserByID("")

		assert.Error(t, err)
		assert.Equal(t, constant.ErrGetID, err)
	})
}

func TestUpdateUser(t *testing.T) {
	mockData := new(MockUserData)
	svc := UserService{d: mockData}

	t.Run("success - Update User", func(t *testing.T) {
		mockEdit := users.EditUser{
			ID:       "1",
			Nama:     "Updated Name",
			Password: "newpassword123",
		}

		// helper.HashPassword = func(password string) (string, error) {
		// 	if password == "newpassword123" {
		// 		return "hashedpassword", nil
		// 	}
		// 	return "", errors.New("hashing failed")
		// }

		mockData.On("GetUserByID", "1").Return(users.User{Password: "oldpassword", NomorHP: "08123456789"}, nil)
		mockData.On("UpdateUser", mock.Anything).Return(nil)

		err := svc.UpdateUser(mockEdit)

		assert.NoError(t, err)
		mockData.AssertExpectations(t)
	})

	t.Run("fail login - empty id", func(t *testing.T) {
		mockEdit := users.EditUser{
			ID:       "",
			Nama:     "Updated Name",
			Password: "newpassword123",
		}

		err := svc.UpdateUser(mockEdit)

		assert.Error(t, err)
		assert.Equal(t, constant.ErrEmptyId, err)
	})
	t.Run("fail login - empty field", func(t *testing.T) {
		mockEdit := users.EditUser{
			ID:   "1",
			Nama: "",
		}

		err := svc.UpdateUser(mockEdit)

		assert.Error(t, err)
		assert.Equal(t, constant.ErrUpdate, err)
	})

}

func TestDeleteUser(t *testing.T) {
	mockData := new(MockUserData)
	service := UserService{d: mockData}

	t.Run("success", func(t *testing.T) {
		mockData.On("DeleteUser", "1").Return(nil)

		err := service.DeleteUser("1")

		assert.NoError(t, err)
		mockData.AssertCalled(t, "DeleteUser", "1")
	})

	t.Run("error empty ID", func(t *testing.T) {
		err := service.DeleteUser("")

		assert.Error(t, err)
		assert.Equal(t, constant.ErrEmptyId, err)
	})
}
