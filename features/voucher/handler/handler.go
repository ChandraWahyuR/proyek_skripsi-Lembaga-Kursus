package handler

import (
	"net/http"
	"skripsi/constant"
	"skripsi/features/voucher"
	"skripsi/helper"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type VoucherHandler struct {
	s voucher.VoucherServiceInterface
	j helper.JWTInterface
}

func New(u voucher.VoucherServiceInterface, j helper.JWTInterface) voucher.VoucherHandlerInterface {
	return &VoucherHandler{
		s: u,
		j: j,
	}
}

func (h *VoucherHandler) GetAllVoucher() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Token
		tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
		if tokenString == "" {
			helper.UnauthorizedError(c)
		}

		ctx := c.Request().Context()
		token, err := h.j.ValidateToken(ctx, tokenString)
		if err != nil {
			return helper.UnauthorizedError(c)
		}

		tokenData := h.j.ExtractUserToken(token)
		role, ok := tokenData[constant.JWT_ROLE]
		if !ok || (role != constant.RoleAdmin && role != constant.RoleUser) {
			return helper.UnauthorizedError(c)
		}
		// Pagination
		// Page
		pageStr := c.QueryParam("page")
		page, err := strconv.Atoi(pageStr)
		if err != nil || page <= 0 {
			page = 1
		}

		// Limit
		limitStr := c.QueryParam("limit")
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			limit = 10
		}
		data, totalPages, err := h.s.GetAllVoucherPagination(page, limit)
		metadata := MetadataResponse{
			TotalPage: totalPages,
			Page:      page,
		}
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		var dataResponse []ResponseGetAllVoucher
		for _, value := range data {
			dataResponse = append(dataResponse, ResponseGetAllVoucher{
				ID:        value.ID,
				Nama:      value.Nama,
				Code:      value.Code,
				Discount:  value.Discount,
				ExpiredAt: value.ExpiredAt,
			})
		}
		return c.JSON(http.StatusOK, helper.MetadataFormatResponse(true, "berhasil", metadata, dataResponse))
	}
}

func (h *VoucherHandler) GetByIDVoucher() echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
		if tokenString == "" {
			helper.UnauthorizedError(c)
		}
		ctx := c.Request().Context()
		token, err := h.j.ValidateToken(ctx, tokenString)
		if err != nil {
			return helper.UnauthorizedError(c)
		}

		tokenData := h.j.ExtractUserToken(token)
		role, ok := tokenData[constant.JWT_ROLE]
		if !ok || (role != constant.RoleAdmin && role != constant.RoleUser) {
			return helper.UnauthorizedError(c)
		}

		id := c.Param("id")
		data, err := h.s.GetByIDVoucher(id)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}
		response := ResponseGetIDVoucher{
			ID:        data.ID,
			Nama:      data.Nama,
			Deskripsi: data.Nama,
			Code:      data.Code,
			Discount:  data.Discount,
			ExpiredAt: data.ExpiredAt,
		}
		return c.JSON(http.StatusOK, helper.ObjectFormatResponse(true, "sukses", response))
	}
}

func (h *VoucherHandler) CreateVoucher() echo.HandlerFunc {
	return func(c echo.Context) error {
		// ini di echo juga punya echo.HeaderAuthorization sama aja lah initnya
		tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
		if tokenString == "" {
			helper.UnauthorizedError(c)
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

		var voucherRequest RequestCreateVoucher
		if err := c.Bind(&voucherRequest); err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		response := voucher.Voucher{
			ID:        uuid.New().String(),
			Nama:      voucherRequest.Nama,
			Deskripsi: voucherRequest.Deskripsi,
			Code:      voucherRequest.Code,
			Discount:  voucherRequest.Discount,
			ExpiredAt: voucherRequest.ExpiredAt,
		}
		err = h.s.CreateVoucher(response)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}
		return c.JSON(http.StatusCreated, helper.ObjectFormatResponse(true, "Success", nil))
	}
}

func (h *VoucherHandler) UpdateVoucher() echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
		if tokenString == "" {
			helper.UnauthorizedError(c)
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

		id := c.Param("id")
		data, err := h.s.GetByIDVoucher(id)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		var voucherRequest RequestCreateVoucher
		if err := c.Bind(&voucherRequest); err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		response := voucher.Voucher{
			Nama:      data.ID,
			Deskripsi: voucherRequest.Deskripsi,
			Code:      voucherRequest.Code,
			Discount:  voucherRequest.Discount,
			ExpiredAt: voucherRequest.ExpiredAt,
		}
		err = h.s.UpdateVoucher(response)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}
		return c.JSON(http.StatusOK, helper.ObjectFormatResponse(true, "Success", nil))
	}
}

func (h *VoucherHandler) DeleteVoucher() echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
		if tokenString == "" {
			helper.UnauthorizedError(c)
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

		id := c.Param("id")
		err = h.s.DeleteVoucher(id)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}
		return c.JSON(http.StatusOK, helper.ObjectFormatResponse(true, "Success", nil))
	}
}
