package handler

import (
	"net/http"
	"skripsi/constant"
	"skripsi/features/users"
	"skripsi/helper"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	s users.UserServiceInterface
	j helper.JWTInterface
}

func New(u users.UserServiceInterface, j helper.JWTInterface) users.UserHandlerInterface {
	return &UserHandler{
		s: u,
		j: j,
	}
}

func (h *UserHandler) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		var reqRegister UserRegisterRequest

		err := c.Bind(&reqRegister)
		if err != nil {
			err, message := helper.HandleEchoError(err)
			return c.JSON(err, helper.FormatResponse(false, message, nil))
		}
		user := users.User{
			Username:        reqRegister.Username,
			Email:           reqRegister.Email,
			NomorHP:         reqRegister.NomorHP,
			Password:        reqRegister.Password,
			ConfirmPassword: reqRegister.ConfirmPassword,
		}
		err = h.s.Register(user)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}
		return c.JSON(http.StatusCreated, helper.FormatResponse(true, "Success", nil))
	}
}

func (h *UserHandler) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		var reqLogin UserLoginRequest

		err := c.Bind(&reqLogin)
		if err != nil {
			err, message := helper.HandleEchoError(err)
			return c.JSON(err, helper.FormatResponse(false, message, nil))
		}

		user := users.User{
			Email:    reqLogin.Email,
			Password: reqLogin.Password,
		}

		userData, err := h.s.Login(user)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
		}

		var response UserLoginResponse
		response.Token = userData.Token

		return c.JSON(http.StatusOK, helper.FormatResponse(true, "Login Success", response))
	}
}

func (h *UserHandler) ForgotPassword() echo.HandlerFunc {
	return func(c echo.Context) error {
		var reqForgotPassword ForgotPasswordRequest

		err := c.Bind(&reqForgotPassword)
		if err != nil {
			err, message := helper.HandleEchoError(err)
			return c.JSON(err, helper.FormatResponse(false, message, nil))
		}

		user := users.User{
			Email: reqForgotPassword.Email,
		}

		token, err := h.s.ForgotPassword(user)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse(true, "Success", map[string]string{"token": token}))
	}
}

func (h *UserHandler) VerifyOTP() echo.HandlerFunc {
	return func(c echo.Context) error {
		var reqVerifyOTP VerifyOtpRequest

		err := c.Bind(&reqVerifyOTP)
		if err != nil {
			err, message := helper.HandleEchoError(err)
			return c.JSON(err, helper.FormatResponse(false, message, nil))
		}

		tokenString := c.Request().Header.Get("Authorization")
		token, err := h.j.ValidateToken(tokenString)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		extract := h.j.ExtractUserToken(token)
		email, ok := extract[constant.JWT_EMAIL].(string)
		if !ok || email == "" {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Invalid token: Email not found", nil))
		}

		user := users.VerifyOtp{
			Email:  email,
			Otp:    reqVerifyOTP.Otp,
			Status: "Success",
		}

		err = h.s.VerifyOTP(user)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}
		return c.JSON(http.StatusOK, helper.FormatResponse(true, "Success", nil))
	}
}

func (h *UserHandler) ResetPassword() echo.HandlerFunc {
	return func(c echo.Context) error {
		var reqResetPassword ResetPasswordRequest

		err := c.Bind(&reqResetPassword)
		if err != nil {
			err, message := helper.HandleEchoError(err)
			return c.JSON(err, helper.FormatResponse(false, message, nil))
		}

		// Pastikan data di-bind dengan benar
		if reqResetPassword.Password == "" || reqResetPassword.ConfirmationPassword == "" {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "password, confirm password cannot be empty", nil))
		}

		tokenString := c.Request().Header.Get("Authorization")
		token, err := h.j.ValidateToken(tokenString)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		extract := h.j.ExtractUserToken(token)
		email, ok := extract[constant.JWT_EMAIL].(string)
		if !ok || email == "" {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Invalid token: Email not found", nil))
		}

		user := users.ResetPassword{
			Email:                email,
			Password:             reqResetPassword.Password,
			ConfirmationPassword: reqResetPassword.ConfirmationPassword,
		}

		err = h.s.ResetPassword(user)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse(true, "Success", nil))
	}
}
