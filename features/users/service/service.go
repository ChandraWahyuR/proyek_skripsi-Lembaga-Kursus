package service

import (
	"errors"
	"fmt"
	"math/rand"
	"skripsi/constant"
	"skripsi/features/users"
	"skripsi/helper"
	"strings"
)

type UserService struct {
	d users.UserDataInterface
	j helper.JWTInterface
	m helper.Sent
}

func New(u users.UserDataInterface, j helper.JWTInterface, m helper.Sent) users.UserServiceInterface {
	return &UserService{
		d: u,
		j: j,
		m: m,
	}
}

func (s *UserService) Register(users users.User) error {
	switch {
	case users.Username == "":
		return constant.ErrEmptyNameRegister
	case users.Email == "":
		return constant.ErrEmptyEmailRegister
	case users.Password == "":
		return constant.ErrEmptyPasswordRegister
	case users.ConfirmPassword == "":
		return errors.New("confirm password cannot be empty")
	case users.NomorHP == "":
		return errors.New("telephone numbber cannot be empty and must be 11 or 12 digit")
	}
	users.Email = strings.ToLower(users.Email)
	isEmailValid := helper.ValidateEmail(users.Email)
	if !isEmailValid {
		return constant.ErrInvalidEmail
	}
	if users.Password != users.ConfirmPassword {
		return constant.ErrPasswordNotMatch
	}

	// validate username
	isUsernameValid, err := helper.ValidateUsername(users.Username)
	if err != nil {
		return constant.ErrInvalidUsername
	}

	// validate password
	pass, err := helper.ValidatePassword(users.Password)
	if err != nil {
		return err
	}

	// hashing password
	hashedPassword, err := helper.HashPassword(pass)
	if err != nil {
		return err
	}
	users.Password = hashedPassword

	// format nomor hp
	nomorHp, err := helper.TelephoneValidator(users.NomorHP)
	if err != nil {
		return err
	}
	users.Username = isUsernameValid
	users.NomorHP = nomorHp

	err = s.d.Register(users)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) Login(user users.User) (users.Login, error) {
	if user.Email == "" || user.Password == "" {
		return users.Login{}, constant.ErrEmptyLogin
	}
	isEmailValid := helper.ValidateEmail(user.Email)
	if !isEmailValid {
		return users.Login{}, constant.ErrInvalidEmail
	}
	user.Email = strings.ToLower(user.Email)

	userData, err := s.d.Login(user)
	if err != nil {
		return users.Login{}, err
	}

	var UserLogin helper.UserJWT

	UserLogin.ID = userData.ID
	UserLogin.Email = userData.Email
	UserLogin.Role = constant.RoleUser

	token, err := s.j.GenerateUserJWT(UserLogin)
	if err != nil {
		return users.Login{}, err
	}

	var UserLoginData users.Login
	UserLoginData.Token = token

	return UserLoginData, nil
}

func (s *UserService) ForgotPassword(forgot users.User) (string, error) {
	if forgot.Email == "" {
		return "", errors.New("email cannot be empty")
	}

	isEmailValid := helper.ValidateEmail(forgot.Email)
	if !isEmailValid {
		return "", constant.ErrInvalidEmail
	}

	var dataOtp users.ForgotPassword

	dataOtp.Email = forgot.Email
	dataOtp.Otp = fmt.Sprintf("%05d", rand.Intn(100000))
	subject := "Reset Password"

	if err := s.d.ForgotPassword(dataOtp); err != nil {
		return "", err
	}

	if err := s.m.SendEmail(forgot.Email, subject, dataOtp.Otp); err != nil {
		return "", err
	}

	// Buat ambil uuid nya aja, biar token jwt nya jelas aja
	userData, err := s.d.GetByEmail(forgot.Email)
	if err != nil {
		return "", err
	}

	token, err := s.j.GenerateForgotPassJWT(helper.ForgotPassJWT{
		Email: forgot.Email, ID: userData.ID})
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *UserService) VerifyOTP(verifikasi users.VerifyOtp) error {
	if verifikasi.Otp == "" {
		return constant.ErrEmptyOtp
	}

	isEmailValid := helper.ValidateEmail(verifikasi.Email)
	if !isEmailValid {
		return constant.ErrInvalidEmail
	}

	err := s.d.VerifyOTP(verifikasi)
	if err != nil {
		return constant.ErrBadRequest
	}

	return nil
}

func (s *UserService) ResetPassword(reset users.ResetPassword) error {
	if reset.Password == "" || reset.ConfirmationPassword == "" {
		return errors.New("password, confirm password cannot be empty")
	}

	if reset.Password != reset.ConfirmationPassword {
		return constant.ErrPasswordNotMatch
	}

	isPasswordValid, err := helper.ValidatePassword(reset.Password)
	if err != nil {
		return err
	}

	hashPass, err := helper.HashPassword(isPasswordValid)
	if err != nil {
		return err
	}

	reset.Password = hashPass

	if err := s.d.ResetPassword(reset); err != nil {
		return err
	}

	return nil
}

func (s *UserService) ActivateAccount(email string) error {
	err := s.d.VerifyEmail(email, true)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) SendVerificationEmail(email, link string) error {
	subject := "Verifikasi Email"
	body := fmt.Sprintf("<p>Klik email dibawah Untuk memverifikasi email:</p><a href='%s'>Verify</a>", link)
	return s.m.SendEmail(email, subject, body)
}

func (s *UserService) GetAllUserPagination(page, limit int) ([]users.GetUser, int, error) {
	return s.d.GetAllUserPagination(page, limit)
}

func (s *UserService) GetUserByID(id string) (users.User, error) {
	if id == "" {
		return users.User{}, constant.ErrGetID
	}
	return s.d.GetUserByID(id)
}

func (s *UserService) UpdateUser(data users.EditUser) error {
	if data.ID == "" {
		return constant.ErrEmptyId
	}
	if data.Nama == "" && data.Username == "" && data.ProfileUrl == "" && data.Password == "" && data.NomorHP == "" && data.Agama == "" && data.Gender == "" && data.TempatLahir == "" && data.TanggalLahir.IsZero() && data.OrangTua == "" && data.Profesi == "" && data.KTP == "" && data.KartuKeluarga == "" {
		return constant.ErrUpdate
	}

	if data.Password != "" {
		hashedPassword, err := helper.HashPassword(data.Password)
		if err != nil {
			return err
		}
		data.Password = hashedPassword
	} else {
		oldUserData, err := s.d.GetUserByID(data.ID)
		if err != nil {
			return err
		}
		data.Password = oldUserData.Password
	}

	if data.NomorHP != "" {
		nomorHp, err := helper.TelephoneValidator(data.NomorHP)
		if err != nil {
			return err
		}
		data.NomorHP = nomorHp
	} else {
		oldUserData, err := s.d.GetUserByID(data.ID)
		if err != nil {
			return err
		}
		data.NomorHP = oldUserData.NomorHP
	}

	return s.d.UpdateUser(data)
}

func (s *UserService) DeleteUser(userId string) error {
	if userId == "" {
		return constant.ErrEmptyId
	}
	return s.d.DeleteUser(userId)
}
