package helper

import (
	"log"
	"net/http"
	"skripsi/constant"

	"github.com/labstack/echo/v4"
)

func ConverResponse(err error) int {
	log.Printf("Received error: %v", err)
	switch err {
	// General errors
	case constant.ErrBadRequest:
		return http.StatusBadRequest
	case constant.ErrUnauthorized:
		return http.StatusUnauthorized
	case constant.ErrEmptyOtp:
		return http.StatusBadRequest
	case constant.ErrDataNotfound, constant.ErrKursusNotfound, constant.ErrInstrukturNotFound, constant.ErrKategoriNotFound, constant.ErrUserNotFound:
		return http.StatusNotFound
	case constant.ErrGetData, constant.ErrGetInstruktur, constant.ErrGetID:
		return http.StatusNotFound
	case constant.ErrEmptyId:
		return http.StatusBadRequest

	// JWT errors
	case constant.ErrGenerateJWT, constant.ErrValidateJWT:
		return http.StatusUnauthorized

	// Update and validator errors
	case constant.ErrUpdate, constant.ErrHashPassword:
		return http.StatusInternalServerError

	// Register errors
	case constant.ErrEmptyEmailRegister, constant.ErrEmptyNameRegister, constant.ErrEmptyPasswordRegister, constant.ErrPasswordNotMatch, constant.ErrInvalidEmail, constant.ErrInvalidUsername, constant.ErrInvalidPhone:
		return http.StatusBadRequest

	// Login errors
	case constant.ErrEmptyLogin, constant.ErrEmptyPasswordLogin, constant.ErrInvalidPassword, constant.ErrLenPassword:
		return http.StatusBadRequest

	// Admin errors
	case constant.ErrAdminNotFound, constant.ErrAdminUserNameEmpty, constant.ErrAdminPasswordEmpty, constant.ErrEmptyGender, constant.ErrGenderChoice:
		return http.StatusBadRequest

	// Instruktur errors
	case constant.ErrInstrukturNotFound, constant.ErrInstrukturID, constant.ErrEmptyNameInstuktor, constant.ErrEmptyEmailInstuktor, constant.ErrEmptyAlamatInstuktor, constant.ErrEmptyNumbertelponInstuktor, constant.ErrEmptyDescriptionInstuktor:
		return http.StatusBadRequest

	// Kategori errors
	case constant.ErrKategoriNotFound, constant.ErrEmptyNamaKategori, constant.ErrEmptyImageUrlKategori, constant.ErrEmptyDeskripsiKategori:
		return http.StatusBadRequest

	// Kursus errors
	case constant.ErrKursusNotFound, constant.ErrJadwal, constant.ErrJadwalFormat, constant.ErrGambarKursus, constant.ErrKategoriKursus, constant.ErrMateriPembelajaran:
		return http.StatusBadRequest

	// GCS errors
	case constant.ErrOpeningFile, constant.ErrUploadGCS:
		return http.StatusInternalServerError

	// Voucher errors
	case constant.ErrVoucherNotFound, constant.ErrVoucherFailedCreate, constant.ErrVoucherUsed:
		return http.StatusBadRequest
	case constant.ErrVoucherIDNotFound:
		return http.StatusNotFound

	// Transaksi errors
	case constant.ErrTransaksiNotFound, constant.ErrValidateDokumenUser:
		return http.StatusBadRequest
	// Default case for internal server errors
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
