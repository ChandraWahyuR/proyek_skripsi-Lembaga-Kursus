package users

import (
	"time"

	"github.com/labstack/echo/v4"
)

type User struct {
	ID              string
	Email           string
	Username        string
	Password        string
	ConfirmPassword string
	NomorHP         string
	IsActive        bool
	Nama            string
	NIS             string
	Agama           string
	Gender          string
	TempatLahir     string
	TanggalLahir    time.Time
	OrangTua        string
	Profesi         string
	Ijazah          string
	KTP             string
	KartuKeluarga   string
	ProfileUrl      string
}

type GetUser struct {
	ID            string
	Username      string
	Email         string
	Password      string
	NomorHP       string
	ProfileUrl    string
	IsActive      bool
	Nama          string
	NIS           string
	Agama         string
	Gender        string
	TempatLahir   string
	TanggalLahir  time.Time
	OrangTua      string
	Profesi       string
	Ijazah        string
	KTP           string
	KartuKeluarga string
}

type EditUser struct {
	ID            string
	Username      string
	ProfileUrl    string
	Password      string
	NomorHP       string
	Nama          string
	Agama         string
	Gender        string
	TempatLahir   string
	TanggalLahir  time.Time
	OrangTua      string
	Profesi       string
	Ijazah        string
	KTP           string
	KartuKeluarga string
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

	// User
	GetAllUser() echo.HandlerFunc
	GetUserByID() echo.HandlerFunc
	GetUserByUser() echo.HandlerFunc
	UpdateUser() echo.HandlerFunc
	DeleteUser() echo.HandlerFunc
	Logout() echo.HandlerFunc
	VerifyAccount() echo.HandlerFunc
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
	ActivateAccount(email string) error
	SendVerificationEmail(email, link string) error

	//
	GetAllUserPagination(page, limit int) ([]User, int, error)
	GetUserByID(id string) (User, error)
	UpdateUser(EditUser) error
	DeleteUser(userId string) error
}

type UserDataInterface interface {
	Register(user User) error
	Login(user User) (User, error)

	ForgotPassword(ForgotPassword) error
	VerifyOTP(VerifyOtp) error
	ResetPassword(ResetPassword) error

	IsEmailExist(email string) bool
	IsUsernameExist(username string) bool
	GetByEmail(email string) (User, error)

	// Auth Email
	// VerifyEmail(email string, otp string) error
	VerifyEmail(email string, isValid bool) error

	//
	GetAllUserPagination(page, limit int) ([]User, int, error)
	GetUserByID(id string) (User, error)
	UpdateUser(EditUser) error
	DeleteUser(userId string) error
}
