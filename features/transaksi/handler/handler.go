package handler

import (
	"net/http"
	"skripsi/constant"
	"skripsi/features/transaksi"
	"skripsi/helper"
	"strconv"

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
		ctx := c.Request().Context()
		token, err := h.j.ValidateToken(ctx, tokenString)
		if err != nil {
			return helper.UnauthorizedError(c)
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

		// Transaksi History
		dataHistory := transaksi.TransaksiHistory{
			ID:          uuid.New().String(),
			KursusID:    dataTransaksi.Kursus,
			UserID:      userId.(string),
			TransaksiID: transaksiResponse.ID,
			Status:      "Not Active",
		}
		err = h.s.CreateTransaksiHistory(dataHistory)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}
		return c.JSON(http.StatusCreated, helper.FormatResponse(true, "Transaksi has been sent", paymentResponse))
	}
}

func (h *TransaksiHanlder) GetAllStatusTransaksi() echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
		if tokenString == "" {
			return helper.UnauthorizedError(c)
		}
		ctx := c.Request().Context()
		token, err := h.j.ValidateToken(ctx, tokenString)
		if err != nil {
			return helper.UnauthorizedError(c)
		}

		tokenData := h.j.ExtractAdminToken(token)
		role, ok := tokenData[constant.JWT_ROLE]
		if !ok || role != constant.RoleAdmin {
			return helper.UnauthorizedError(c)
		}

		// Pagination
		pageStr := c.QueryParam("page")
		page, err := strconv.Atoi(pageStr)
		if err != nil || page <= 0 {
			page = 1
		}

		limitStr := c.QueryParam("limit")
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			limit = 10
		}

		data, totalPages, err := h.s.GetAllTransaksiPagination(page, limit)
		metadata := MetadataResponse{
			TotalPage: totalPages,
			Page:      page,
		}

		var response []GetAllTransaksiAdminResponse
		for _, dataResponse := range data {
			response = append(response, GetAllTransaksiAdminResponse{
				ID: dataResponse.ID,
				User: User{
					ID:       dataResponse.UserID,
					Username: dataResponse.User.Username,
					Email:    dataResponse.User.Email,
				},
				Kursus: Kursus{
					ID:   dataResponse.KursusID,
					Nama: dataResponse.Kursus.Nama,
				},
				TotalHarga: dataResponse.TotalHarga,
				Status:     dataResponse.Status,
			})
		}
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}
		return c.JSON(http.StatusOK, helper.MetadataFormatResponse(true, "Succses", metadata, response))
	}
}

func (h *TransaksiHanlder) GetStatusTransaksiForUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
		if tokenString == "" {
			return helper.UnauthorizedError(c)
		}
		ctx := c.Request().Context()
		token, err := h.j.ValidateToken(ctx, tokenString)
		if err != nil {
			return helper.UnauthorizedError(c)
		}
		tokenData := h.j.ExtractUserToken(token)
		role, ok := tokenData[constant.JWT_ROLE]
		userId := tokenData[constant.JWT_ID].(string)
		if !ok || role != constant.RoleUser {
			return helper.UnauthorizedError(c)
		}

		// Pagination
		pageStr := c.QueryParam("page")
		page, err := strconv.Atoi(pageStr)
		if err != nil || page <= 0 {
			page = 1
		}

		limitStr := c.QueryParam("limit")
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			limit = 10
		}

		data, totalPages, err := h.s.GetStatusTransaksiForUser(userId, page, limit)
		metadata := MetadataResponse{
			TotalPage: totalPages,
			Page:      page,
		}

		var response []GetAllTransaksiUserResponse
		for _, dataResponse := range data {
			response = append(response, GetAllTransaksiUserResponse{
				ID: dataResponse.ID,
				Kursus: Kursus{
					ID:   dataResponse.KursusID,
					Nama: dataResponse.Kursus.Nama,
				},
				TotalHarga: dataResponse.TotalHarga,
				Status:     dataResponse.Status,
			})
		}
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}
		return c.JSON(http.StatusOK, helper.MetadataFormatResponse(true, "Succses", metadata, response))
	}
}

func (h *TransaksiHanlder) GetStatusTransaksiByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
		if tokenString == "" {
			return helper.UnauthorizedError(c)
		}
		ctx := c.Request().Context()
		token, err := h.j.ValidateToken(ctx, tokenString)
		if err != nil {
			return helper.UnauthorizedError(c)
		}
		tokenData := h.j.ExtractAdminToken(token)
		role, ok := tokenData[constant.JWT_ROLE]
		if !ok || (role != constant.RoleAdmin && role != constant.RoleUser) {
			return helper.UnauthorizedError(c)
		}

		id := c.Param("id")
		dataTransaksi, err := h.s.GetStatusTransaksiByID(id)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}
		responseData := GetAllTransaksiByIDAdminResponse{
			ID: dataTransaksi.ID,
			User: User{
				ID:       dataTransaksi.UserID,
				Username: dataTransaksi.User.Username,
				Email:    dataTransaksi.User.Email,
			},
			Kursus: Kursus{
				ID:   dataTransaksi.KursusID,
				Nama: dataTransaksi.Kursus.Nama,
			},
			VoucherID:  dataTransaksi.VoucherID,
			TotalHarga: dataTransaksi.TotalHarga,
			Status:     dataTransaksi.Status,
		}

		return c.JSON(http.StatusOK, helper.FormatResponse(true, "Success", responseData))
	}
}

func (h *TransaksiHanlder) GetAllTransaksiHistory() echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
		if tokenString == "" {
			return helper.UnauthorizedError(c)
		}
		ctx := c.Request().Context()
		token, err := h.j.ValidateToken(ctx, tokenString)
		if err != nil {
			return helper.UnauthorizedError(c)
		}
		tokenData := h.j.ExtractAdminToken(token)
		role, ok := tokenData[constant.JWT_ROLE]
		if !ok || role != constant.RoleAdmin {
			return helper.UnauthorizedError(c)
		}
		pageStr := c.QueryParam("page")
		page, err := strconv.Atoi(pageStr)
		if err != nil || page <= 0 {
			page = 1
		}

		limitStr := c.QueryParam("limit")
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			limit = 10
		}
		data, totalPages, err := h.s.GetAllHistoryTransaksiPagination(page, limit)
		metadata := MetadataResponse{
			TotalPage: totalPages,
			Page:      page,
		}

		var response []GetAllHistoryAdminResponse
		for _, dataResponse := range data {
			response = append(response, GetAllHistoryAdminResponse{
				ID: dataResponse.ID,
				Transaksi: Transaksi{
					ID:         dataResponse.TransaksiID,
					TotalHarga: dataResponse.Transaksi.TotalHarga,
				},
				KursusID:   dataResponse.KursusID,
				UserID:     dataResponse.UserID,
				Status:     dataResponse.Status,
				ValidUntil: dataResponse.ValidUntil,
			})
		}
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}
		return c.JSON(http.StatusOK, helper.MetadataFormatResponse(true, "Success", metadata, response))
	}
}

func (h *TransaksiHanlder) GetTransaksiHistoryByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
		if tokenString == "" {
			return helper.UnauthorizedError(c)
		}
		ctx := c.Request().Context()
		token, err := h.j.ValidateToken(ctx, tokenString)
		if err != nil {
			return helper.UnauthorizedError(c)
		}
		tokenData := h.j.ExtractAdminToken(token)
		role, ok := tokenData[constant.JWT_ROLE]
		if !ok || (role != constant.RoleAdmin && role != constant.RoleUser) {
			return helper.UnauthorizedError(c)
		}

		id := c.Param("id")
		dataHistoryT, err := h.s.GetTransaksiHistoryByID(id)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		response := GetHistoryAdminByIDResponse{
			ID: dataHistoryT.ID,
			Transaksi: Transaksi{
				ID:         dataHistoryT.Transaksi.ID,
				TotalHarga: dataHistoryT.Transaksi.TotalHarga,
			},
			Kursus: Kursus{
				ID:   dataHistoryT.KursusID,
				Nama: dataHistoryT.Kursus.Nama,
			},
			User: User{
				ID:       dataHistoryT.UserID,
				Username: dataHistoryT.User.Username,
				Email:    dataHistoryT.User.Email,
			},
			VoucherID:  dataHistoryT.Transaksi.VoucherID,
			TotalHarga: dataHistoryT.Transaksi.TotalHarga,
			Status:     dataHistoryT.Status,
			ValidUntil: dataHistoryT.ValidUntil,
		}
		return c.JSON(http.StatusOK, helper.ObjectFormatResponse(true, "Berhasil", response))
	}
}

func (h *TransaksiHanlder) GetAllTransaksiHistoryForUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
		if tokenString == "" {
			return helper.UnauthorizedError(c)
		}
		ctx := c.Request().Context()
		token, err := h.j.ValidateToken(ctx, tokenString)
		if err != nil {
			return helper.UnauthorizedError(c)
		}
		tokenData := h.j.ExtractUserToken(token)
		role, ok := tokenData[constant.JWT_ROLE]
		userId := tokenData[constant.JWT_ID].(string)
		if !ok || role != constant.RoleUser {
			return helper.UnauthorizedError(c)
		}
		pageStr := c.QueryParam("page")
		page, err := strconv.Atoi(pageStr)
		if err != nil || page <= 0 {
			page = 1
		}

		limitStr := c.QueryParam("limit")
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			limit = 10
		}
		data, totalPages, err := h.s.GetAllTransaksiHistoryForUser(userId, page, limit)
		metadata := MetadataResponse{
			TotalPage: totalPages,
			Page:      page,
		}
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		var response []GetAllHistoryUserResponse
		for _, responseData := range data {
			response = append(response, GetAllHistoryUserResponse{
				ID: responseData.ID,
				Transaksi: Transaksi{
					ID:         responseData.TransaksiID,
					TotalHarga: responseData.Transaksi.TotalHarga,
				},
				KursusID:   responseData.KursusID,
				UserID:     responseData.UserID,
				ValidUntil: responseData.ValidUntil,
			})
		}
		return c.JSON(http.StatusOK, helper.MetadataFormatResponse(true, "Success", metadata, response))
	}
}
