package handler

import (
	"net/http"
	"skripsi/constant"
	"skripsi/features/instruktur"
	"skripsi/helper"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type InstrukturHandler struct {
	s instruktur.InstrukturServiceInterface
	j helper.JWTInterface
}

func New(u instruktur.InstrukturServiceInterface, j helper.JWTInterface) InstrukturHandler {
	return InstrukturHandler{
		s: u,
		j: j,
	}
}

func (h InstrukturHandler) GetAllInstruktur() echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
		if tokenString == "" {
			helper.UnauthorizedError(c)
		}
		token, err := h.j.ValidateToken(tokenString)
		if err != nil {
			helper.UnauthorizedError(c)
		}

		adminData := h.j.ExtractAdminToken(token)
		role := adminData[constant.JWT_ROLE]
		if role != constant.RoleAdmin {
			helper.UnauthorizedError(c)
		}

		pageStr := c.QueryParam("page")
		page, err := strconv.Atoi(pageStr)
		if err != nil || page <= 0 {
			page = 1
		}

		limitStr := c.QueryParam("limit")
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			limit = 10 // Limit 10 walaupun fe minta 100
		}

		data, totalPages, err := h.s.GetInstrukturWithPagination(page, limit)
		metadata := MetadataResponse{
			TotalPage: totalPages,
			Page:      page,
		}
		var response []DataInsrukturResponse
		for _, f := range data {
			response = append(response, DataInsrukturResponse{
				ID:     f.ID,
				Name:   f.Name,
				Gender: f.Gender,
				Email:  f.Email,
				Alamat: f.Alamat,
				NoHp:   f.NoHp,
			})
		}

		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		return c.JSON(http.StatusOK, helper.MetadataFormatResponse(true, "Berhasil", metadata, response))
	}
}

func (h InstrukturHandler) GetAllInstrukturByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
		if tokenString == "" {
			helper.UnauthorizedError(c)
		}
		token, err := h.j.ValidateToken(tokenString)
		if err != nil {
			helper.UnauthorizedError(c)
		}

		adminData := h.j.ExtractAdminToken(token)
		role := adminData[constant.JWT_ROLE]
		if role != constant.RoleAdmin {
			helper.UnauthorizedError(c)
		}

		id := c.Param("id")
		dataInstruktur, err := h.s.GetAllInstrukturByID(id)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		dataResponse := DataInsrukturResponse{
			ID:     dataInstruktur.ID,
			Email:  dataInstruktur.Email,
			Name:   dataInstruktur.Name,
			Gender: dataInstruktur.Gender,
			Alamat: dataInstruktur.Alamat,
			NoHp:   dataInstruktur.NoHp,
		}
		return c.JSON(http.StatusOK, helper.ObjectFormatResponse(true, "sukses", dataResponse))
	}
}

func (h InstrukturHandler) PostInstruktur() echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
		if tokenString == "" {
			helper.UnauthorizedError(c)
		}
		token, err := h.j.ValidateToken(tokenString)
		if err != nil {
			helper.UnauthorizedError(c)
		}

		adminData := h.j.ExtractAdminToken(token)
		role := adminData[constant.JWT_ROLE]
		if role != constant.RoleAdmin {
			helper.UnauthorizedError(c)
		}

		var dataRequest PostInstrukturRequest
		err = c.Bind(&dataRequest)
		if err != nil {
			code, message := helper.HandleEchoError(err)
			return c.JSON(code, helper.FormatResponse(false, message, nil))
		}

		validGender := false
		dataRequest.Gender = strings.TrimSpace(strings.ToLower(dataRequest.Gender))
		for _, v := range constant.ValidGenders {
			if v == dataRequest.Gender {
				validGender = true
				break
			}
		}
		if !validGender {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, constant.ErrGenderChoice.Error(), nil))
		}

		dataInstruktur := instruktur.Instruktur{
			ID:     uuid.New().String(),
			Email:  dataRequest.Email,
			Name:   dataRequest.Name,
			Gender: dataRequest.Gender,
			Alamat: dataRequest.Alamat,
			NoHp:   dataRequest.NoHp,
		}
		err = h.s.PostInstruktur(dataInstruktur)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		return c.JSON(http.StatusCreated, helper.FormatResponse(true, "Success", nil))
	}
}

func (h InstrukturHandler) UpdateInstruktur() echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
		if tokenString == "" {
			helper.UnauthorizedError(c)
		}
		token, err := h.j.ValidateToken(tokenString)
		if err != nil {
			helper.UnauthorizedError(c)
		}

		adminData := h.j.ExtractAdminToken(token)
		role := adminData[constant.JWT_ROLE]
		if role != constant.RoleAdmin {
			helper.UnauthorizedError(c)
		}
		id := c.Param("id")
		dataId, err := h.s.GetAllInstrukturByID(id)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		var updataRequest UpdateInstrukturRequest
		err = c.Bind(&updataRequest)
		if err != nil {
			code, message := helper.HandleEchoError(err)
			return c.JSON(code, helper.FormatResponse(false, message, nil))
		}

		dataResponse := instruktur.UpdateInstruktur{
			ID:     dataId.ID,
			Name:   updataRequest.Name,
			Gender: updataRequest.Gender,
			Email:  updataRequest.Email,
			Alamat: updataRequest.Alamat,
			NoHp:   updataRequest.NoHp,
		}

		err = h.s.UpdateInstruktur(dataResponse)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}
		return c.JSON(http.StatusOK, helper.ObjectFormatResponse(true, "Success", nil))
	}
}

func (h InstrukturHandler) DeleteInstruktur() echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
		if tokenString == "" {
			helper.UnauthorizedError(c)
		}
		token, err := h.j.ValidateToken(tokenString)
		if err != nil {
			helper.UnauthorizedError(c)
		}

		adminData := h.j.ExtractAdminToken(token)
		role := adminData[constant.JWT_ROLE]
		if role != constant.RoleAdmin {
			helper.UnauthorizedError(c)
		}
		id := c.Param("id")
		err = h.s.DeleteInstruktur(id)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse(true, "Success", nil))
	}
}

func (h InstrukturHandler) GetInstruktorByName() echo.HandlerFunc {
	return func(c echo.Context) error {
		name := c.QueryParam("name") // Ambil parameter nama dari query URL
		tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
		if tokenString == "" {
			helper.UnauthorizedError(c)
		}
		token, err := h.j.ValidateToken(tokenString)
		if err != nil {
			helper.UnauthorizedError(c)
		}

		adminData := h.j.ExtractAdminToken(token)
		role := adminData[constant.JWT_ROLE]
		if role != constant.RoleAdmin {
			helper.UnauthorizedError(c)
		}

		pageStr := c.QueryParam("page")
		page, err := strconv.Atoi(pageStr)
		if err != nil || page <= 0 {
			page = 1
		}

		limitStr := c.QueryParam("limit")
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			limit = 10 // Limit 10 walaupun fe minta 100
		}

		result, totalPages, err := h.s.GetInstruktorByName(name, page, limit)
		metadata := MetadataResponse{
			TotalPage: totalPages,
			Page:      page,
		}

		var response []DataInsrukturResponse
		for _, f := range result {
			response = append(response, DataInsrukturResponse{
				ID:     f.ID,
				Name:   f.Name,
				Gender: f.Gender,
				Email:  f.Email,
				Alamat: f.Alamat,
				NoHp:   f.NoHp,
			})
		}
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		return c.JSON(http.StatusOK, helper.MetadataFormatResponse(true, "Berhasil", metadata, response))
	}
}
