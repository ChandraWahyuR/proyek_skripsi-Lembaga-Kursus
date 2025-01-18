package handler

import (
	"net/http"
	"skripsi/constant"
	"skripsi/features/transaksi"
	"skripsi/features/voucher"
	"skripsi/helper"
	"strconv"
	"time"

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

		// Validasi voucher
		if dataTransaksi.VoucherID != "" {
			voucherExists, err := h.s.CheckVoucherExists(dataTransaksi.VoucherID)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, "Failed to validate voucher", nil))
			}
			if !voucherExists {
				return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Invalid voucher ID", nil))
			}
			response.VoucherID = dataTransaksi.VoucherID
		}

		transaksiResponse, err := h.s.CreateTransaksi(response)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		if dataTransaksi.VoucherID != "" {
			response.VoucherID = dataTransaksi.VoucherID
			usedVoucher := voucher.VoucherUsed{
				ID:        uuid.New().String(),
				VoucherID: response.VoucherID,
				UserID:    userId.(string),
			}
			err = h.s.UsedVoucher(usedVoucher)
			if err != nil {
				return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
			}
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
			ValidUntil:  time.Time{},
		}

		err = h.s.CreateTransaksiHistory(dataHistory)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		return c.JSON(http.StatusCreated, helper.FormatResponse(true, constant.PostTransaksi, paymentResponse))
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
		return c.JSON(http.StatusOK, helper.MetadataFormatResponse(true, constant.GetAllTransaski, metadata, response))
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
		return c.JSON(http.StatusOK, helper.MetadataFormatResponse(true, constant.GetAllTransaski, metadata, response))
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

		return c.JSON(http.StatusOK, helper.FormatResponse(true, constant.GetAllTransaski, responseData))
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
		return c.JSON(http.StatusOK, helper.MetadataFormatResponse(true, constant.GetAllHistoryTransaki, metadata, response))
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
		return c.JSON(http.StatusOK, helper.ObjectFormatResponse(true, constant.GetAllHistoryTransaki, response))
	}
}

func (h *TransaksiHanlder) GetResponseTransaksi() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.QueryParam("order_id")
		dataHistoryT, err := h.s.GetTransaksiHistoryByID(id)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}
		var voucherID string
		if dataHistoryT.Transaksi.VoucherID != "" {
			voucherID = dataHistoryT.Transaksi.VoucherID
		} else {
			voucherID = "Kosong"
		}
		response := GetTransaksiResponse{
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
			VoucherID:  voucherID,
			TotalHarga: dataHistoryT.Transaksi.TotalHarga,
		}
		return c.JSON(http.StatusOK, helper.ObjectFormatResponse(true, constant.ResponTransaksi, response))
	}
}

func (h *TransaksiHanlder) GetStatusTransaksiForUserByID() echo.HandlerFunc {
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

		id := c.Param("id")
		dataTransaksi, err := h.s.GetStatusTransaksiByID(id)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}
		if dataTransaksi.UserID != userId {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, "bukan haknya kamu!", nil))
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

		return c.JSON(http.StatusOK, helper.ObjectFormatResponse(true, constant.GetAllTransaski, responseData))
	}
}

// History Transaksi
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
		return c.JSON(http.StatusOK, helper.MetadataFormatResponse(true, constant.GetAllHistoryTransaki, metadata, response))
	}
}

func (h *TransaksiHanlder) GetAllTransaksiHistoryForUserByID() echo.HandlerFunc {
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

		id := c.Param("id")
		dataHistoryT, err := h.s.GetTransaksiHistoryByID(id)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}
		if dataHistoryT.UserID != userId {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, "bukan haknya kamu!", nil))
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

		return c.JSON(http.StatusOK, helper.ObjectFormatResponse(true, constant.GetAllHistoryTransaki, response))
	}
}

func (h *TransaksiHanlder) GetActiveUsersFromTransaksiHistory() echo.HandlerFunc {
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
		dataUser, totalPages, err := h.s.GetActiveUsersFromTransaksiHistory(page, limit)
		metadata := MetadataResponse{
			TotalPage: totalPages,
			Page:      page,
		}
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		var response []GetAllUserActiveResponse
		for _, dataResponse := range dataUser {
			response = append(response, GetAllUserActiveResponse{
				UserID:   dataResponse.UserID,
				UserName: dataResponse.User.Username,
				Email:    dataResponse.User.Email,
				Kursus: Kursus{
					ID:   dataResponse.KursusID,
					Nama: dataResponse.Kursus.Nama,
				},
				Status:     dataResponse.Status,
				ValidUntil: dataResponse.ValidUntil,
			})
		}
		return c.JSON(http.StatusOK, helper.MetadataFormatResponse(true, constant.GetAllTransaski, metadata, response))
	}
}
