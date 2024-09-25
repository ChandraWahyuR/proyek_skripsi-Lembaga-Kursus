package helper

import (
	"net/http"
	"skripsi/constant"

	"github.com/labstack/echo/v4"
)

func ConverResponse(err error) int {
	switch err {
	case constant.ErrBadRequest:
		return http.StatusBadRequest
	case constant.ErrUnauthorized:
		return http.StatusUnauthorized
	case constant.ErrEmptyOtp:
		return http.StatusBadRequest
	case constant.ErrInvalidPhone:
		return http.StatusBadRequest
	case constant.ErrEmptyEmailRegister:
		return http.StatusBadRequest
	case constant.ErrEmptyNameRegister:
		return http.StatusBadRequest
	case constant.ErrEmptyPasswordRegister:
		return http.StatusBadRequest
	case constant.ErrPasswordNotMatch:
		return http.StatusBadRequest
	case constant.ErrInvalidEmail:
		return http.StatusBadRequest
	case constant.ErrInvalidUsername:
		return http.StatusBadRequest
	case constant.ErrInvalidPassword:
		return http.StatusBadRequest
	case constant.ErrEmptyPasswordLogin:
		return http.StatusBadRequest

	default:
		return http.StatusInternalServerError
	}
}

func HandleEchoError(err error) (int, string) {
	if _, ok := err.(*echo.HTTPError); ok {
		return http.StatusBadRequest, constant.BadInput
	}
	return http.StatusBadRequest, constant.BadInput
}

func UnauthorizedError(c echo.Context) error {
	return c.JSON(http.StatusUnauthorized, FormatResponse(false, constant.Unauthorized, nil))
}
func InternalServerError(c echo.Context) error {
	return c.JSON(http.StatusInternalServerError, FormatResponse(false, constant.InternalServerError, nil))
}
func JWTErrorHandler(c echo.Context, err error) error {
	return c.JSON(http.StatusUnauthorized, FormatResponse(false, constant.Unauthorized, nil))
}
