package users

import (
	"time"

	"github.com/labstack/echo/v4"
)

type User struct {
	ID              string
	Email           string
	Name            string
	Password        string
	ConfirmPassword string
	NomorHP         string
	IsActive        bool
}

type Login struct {
	Email    string
	Password string
	Token    string
}

type ForgotPassword struct {
	ID        string
	Email     string
	Otp       string
	ExpiredAt time.Time
}

type VerifyOtp struct {
	Email     string
	Otp       string
	Status    string
	ExpiredAt time.Time
}

type ResetPassword struct {
	Email                string
	Password             string
	ConfirmationPassword string
	OTP                  string
}

type UserHandlerInterface interface {
	Register() echo.HandlerFunc
	Login() echo.HandlerFunc

	// Forgot Password
	ForgotPassword() echo.HandlerFunc
	VerifyOTP() echo.HandlerFunc
	ResetPassword() echo.HandlerFunc

	// Auth Email
	// AuthEmail() echo.HandlerFunc
	// ConfirmValidatateUser() echo.HandlerFunc
}

type UserServiceInterface interface {
	Register(user User) error
	Login(user User) (Login, error)

	ForgotPassword(User) (string, error)
	VerifyOTP(VerifyOtp) error
	ResetPassword(ResetPassword) error

	// Auth Email
	// VerifyEmail(User) error
}

type UserDataInterface interface {
	Register(user User) error
	Login(user User) (User, error)

	ForgotPassword(ForgotPassword) error
	VerifyOTP(VerifyOtp) error
	ResetPassword(ResetPassword) error

	IsEmailExist(email string) bool
	GetByEmail(email string) (User, error)

	// Auth Email
	// VerifyEmail(email string, otp string) error
}
