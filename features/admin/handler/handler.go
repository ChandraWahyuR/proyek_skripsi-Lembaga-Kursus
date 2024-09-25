package handler

import (
	"net/http"
	"skripsi/features/admin"
	"skripsi/helper"

	"github.com/labstack/echo/v4"
)

type AdminHandler struct {
	s admin.AdminServiceInterface
	j helper.JWTInterface
}

func New(u admin.AdminServiceInterface, j helper.JWTInterface) AdminHandler {
	return AdminHandler{
		s: u,
		j: j,
	}
}
func (h AdminHandler) RegisterAdmin() echo.HandlerFunc {
	return func(c echo.Context) error {
		var reqRegister AdminRegisterRequest
		err := c.Bind(&reqRegister)
		if err != nil {
			err, message := helper.HandleEchoError(err)
			return c.JSON(err, helper.FormatResponse(false, message, nil))
		}

		admin := admin.Admin{
			Username:        reqRegister.Username,
			Email:           reqRegister.Email,
			Password:        reqRegister.Password,
			ConfirmPassword: reqRegister.ConfirmPassword,
		}
		err = h.s.RegisterAdmin(admin)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}
		return c.JSON(http.StatusCreated, helper.FormatResponse(true, "Success", nil))
	}
}
func (h AdminHandler) LoginAdmin() echo.HandlerFunc {
	return func(c echo.Context) error {
		var reqLogin AdminLoginRequest
		err := c.Bind(&reqLogin)
		if err != nil {
			err, message := helper.HandleEchoError(err)
			return c.JSON(err, helper.FormatResponse(false, message, nil))
		}

		admin := admin.Admin{
			Username: reqLogin.Username,
			Password: reqLogin.Password,
		}
		adminData, err := h.s.LoginAdmin(admin)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}
		var response AdminLoginResponse
		response.Token = adminData.Token
		return c.JSON(http.StatusOK, helper.FormatResponse(true, "Login Success", response))
	}
}
