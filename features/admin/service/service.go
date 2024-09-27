package service

import (
	"skripsi/constant"
	"skripsi/features/admin"
	"skripsi/helper"
	"strings"
)

type AdminService struct {
	d admin.AdminDataInterface
	j helper.JWTInterface
}

func New(d admin.AdminDataInterface, j helper.JWTInterface) admin.AdminServiceInterface {
	return &AdminService{
		d: d,
		j: j,
	}
}

func (s *AdminService) RegisterAdmin(user admin.Admin) error {
	switch {
	case user.Email == "":
		return constant.ErrEmptyEmailRegister
	case user.Username == "":
		return constant.ErrEmptyNameRegister
	case user.Password == "":
		return constant.ErrEmptyPasswordRegister
	case user.Password != user.ConfirmPassword:
		return constant.ErrPasswordNotMatch
	}
	user.Email = strings.ToLower(user.Email)

	isEmailValid := helper.ValidateEmail(user.Email)
	if !isEmailValid {
		return constant.ErrInvalidEmail
	}

	isUsernameValid, err := helper.ValidateUsername(user.Username)
	if err != nil {
		return constant.ErrInvalidUsername
	}

	pass, err := helper.ValidatePassword(user.Password)
	if err != nil {
		return constant.ErrInvalidPassword
	}

	hashedPassword, err := helper.HashPassword(pass)
	if err != nil {
		return err
	}
	user.Username = isUsernameValid
	user.Password = hashedPassword

	err = s.d.RegisterAdmin(user)
	if err != nil {
		return err
	}
	return nil
}

func (s *AdminService) LoginAdmin(user admin.Admin) (admin.Login, error) {
	switch {
	case user.Username == "":
		return admin.Login{}, constant.ErrAdminUserNameEmpty
	case user.Password == "":
		return admin.Login{}, constant.ErrAdminPasswordEmpty
	}

	adminData, err := s.d.LoginAdmin(user)
	if err != nil {
		return admin.Login{}, err
	}

	var adminLogin helper.AdminJWT
	adminLogin.ID = adminData.ID
	adminLogin.Username = adminData.Username
	adminLogin.Email = adminData.Email

	token, err := s.j.GenerateAdminJWT(adminLogin)
	if err != nil {
		return admin.Login{}, err
	}

	var adminLoginData admin.Login
	adminLoginData.Token = token

	return adminLoginData, nil
}
