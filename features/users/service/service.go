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
	case users.Name == "":
		return constant.ErrEmptyNameRegister
	case users.Email == "":
		return constant.ErrEmptyEmailRegister
	case users.Password == "":
		return constant.ErrEmptyPasswordRegister
	case users.ConfirmPassword == "":
		return errors.New("confirm password cannot be empty")
	}
	users.Email = strings.ToLower(users.Email)
	isEmailValid := helper.ValidateEmail(users.Email)
	if !isEmailValid {
		return constant.ErrInvalidEmail
	}
	if users.Password != users.ConfirmPassword {
		return constant.ErrPasswordNotMatch
	}
	hashedPassword, err := helper.HashPassword(users.Password)
	if err != nil {
		return err
	}

	users.Password = hashedPassword

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

	if err := s.d.ForgotPassword(dataOtp); err != nil {
		return "", err
	}

	if err := s.m.SendEmail(forgot.Email, dataOtp.Otp); err != nil {
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
		return errors.New("otp cannot be empty")
	}

	isEmailValid := helper.ValidateEmail(verifikasi.Email)
	if !isEmailValid {
		return constant.ErrInvalidEmail
	}

	err := s.d.VerifyOTP(verifikasi)
	if err != nil {
		return err
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

	hashPass, err := helper.HashPassword(reset.Password)
	if err != nil {
		return err
	}

	reset.Password = hashPass

	if err := s.d.ResetPassword(reset); err != nil {
		return err
	}

	return nil
}

//
