package handler

import (
	"net/http"
	"skripsi/constant"
	"skripsi/features/transaksi"
	"skripsi/helper"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TransaksiHanlder struct {
	s transaksi.TransaksiServiceInterface
	j helper.JWTInterface
}

func New(u transaksi.TransaksiServiceInterface, j helper.JWTInterface) transaksi.TransaksiHandlerInterface {
	return &TransaksiHanlder{
		s: u,
		j: j,
	}
}

func (h *TransaksiHanlder) CreateTransaksi() echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
		if tokenString == "" {
			return helper.UnauthorizedError(c)
		}
		token, err := h.j.ValidateToken(tokenString)
		if err != nil {
			helper.UnauthorizedError(c)
		}

		tokenData := h.j.ExtractUserToken(token)
		role, ok := tokenData[constant.JWT_ROLE]
		userId := tokenData[constant.JWT_ID]
		if !ok || (role != constant.RoleAdmin && role != constant.RoleUser) {
			return helper.UnauthorizedError(c)
		}

		// Ambil user id
		var dataTransaksi TransaksiRequest
		if err := c.Bind(&dataTransaksi); err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		response := transaksi.Transaksi{
			ID:       uuid.New().String(),
			KursusID: dataTransaksi.Kursus,
			UserID:   userId.(string),
		}
		if dataTransaksi.VoucherID != "" {
			response.VoucherID = dataTransaksi.VoucherID
		}
		transaksiResponse, err := h.s.CreateTransaksi(response)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		paymentResponse := PaymentResponse{
			Amount:  int(transaksiResponse.TotalHarga),
			SnapURL: transaksiResponse.SnapURL,
		}
		return c.JSON(http.StatusCreated, helper.FormatResponse(true, "Kursus added successfully", paymentResponse))
	}
}
