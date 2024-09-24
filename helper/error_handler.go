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
