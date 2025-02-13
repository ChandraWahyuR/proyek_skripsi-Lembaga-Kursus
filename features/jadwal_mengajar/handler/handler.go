package handler

import (
	"net/http"
	"skripsi/constant"
	jadwal "skripsi/features/jadwal_mengajar"
	"skripsi/helper"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type JadwalHandler struct {
	s jadwal.MengajarServiceInterface
	j helper.JWTInterface
}

func New(s jadwal.MengajarServiceInterface, j helper.JWTInterface) jadwal.MengajarHandlerInterface {
	return &JadwalHandler{
		s: s,
		j: j,
	}
}

func (h *JadwalHandler) GetJadwalMengajar() echo.HandlerFunc {
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

		adminData := h.j.ExtractAdminToken(token)
		role, ok := adminData[constant.JWT_ROLE]
		if !ok || role != constant.RoleAdmin {
			return helper.UnauthorizedError(c)
		}

		data, err := h.s.GetJadwalMengajar()
		var response []GetJadwalMengajar
		for _, f := range data {
			response = append(response, GetJadwalMengajar{
				ID: f.ID,
				Instruktur: Instruktur{
					ID:   f.InstrukturID,
					Name: f.Instruktur.InstrukturNama,
				},
				User: User{
					ID:       f.UserID,
					Username: f.User.UserName,
				},
				Kursus: Kursus{
					ID:   f.KursusID,
					Nama: f.Kursus.KursusNama,
				},
				Tanggal: f.Tanggal.Format("01-01-2006"),
				Status:  f.Status,
			})
		}

		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse(true, constant.GetData, response))
	}
}

func (h *JadwalHandler) GetJadwalMengajarByID() echo.HandlerFunc {
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

		userData := h.j.ExtractUserToken(token)
		role, ok := userData[constant.JWT_ROLE]
		if !ok || (role != constant.RoleUser && role != constant.RoleAdmin) {
			return helper.UnauthorizedError(c)
		}

		id := c.Param("id")
		data, err := h.s.GetJadwalMengajarByID(id)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}
		resp := GetDetailJadwalMengajar{
			ID: data.ID,
			Instruktur: Instruktur{
				ID:   data.InstrukturID,
				Name: data.Instruktur.InstrukturNama,
			},
			User: User{
				ID:       data.UserID,
				Username: data.User.UserName,
			},
			Kursus: Kursus{
				ID:   data.KursusID,
				Nama: data.Kursus.KursusNama,
			},
			Tanggal:  data.Tanggal.Format("01-01-2006"),
			JamMulai: data.JamMulai.Format("15:04"),
			JamAkhir: data.JamAkhir.Format("15:04"),
			Status:   data.Status,
		}

		return c.JSON(http.StatusOK, helper.FormatResponse(true, constant.GetData, resp))
	}
}

func (h *JadwalHandler) GetJadwalMengajarForUser() echo.HandlerFunc {
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

		userData := h.j.ExtractUserToken(token)
		userId := userData[constant.JWT_ID].(string)
		role, ok := userData[constant.JWT_ROLE]
		if !ok || role != constant.RoleUser {
			return helper.UnauthorizedError(c)
		}

		data, err := h.s.GetJadwalMengajarForUser(userId)
		var response []GetJadwalMengajar
		for _, f := range data {
			response = append(response, GetJadwalMengajar{
				ID: f.ID,
				Instruktur: Instruktur{
					ID:   f.InstrukturID,
					Name: f.Instruktur.InstrukturNama,
				},
				User: User{
					ID:       f.UserID,
					Username: f.User.UserName,
				},
				Kursus: Kursus{
					ID:   f.KursusID,
					Nama: f.Kursus.KursusNama,
				},
				Tanggal: f.Tanggal.Format("01-01-2006"),
				Status:  f.Status,
			})
		}

		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse(true, constant.GetData, response))
	}
}

func (h *JadwalHandler) CreateJadwalMengajar() echo.HandlerFunc {
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
		if !ok || (role != constant.RoleAdmin) {
			return helper.UnauthorizedError(c)
		}

		// Bind request
		var req Request
		if err := c.Bind(&req); err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		newJadwal := jadwal.JadwalMengajar{
			ID:           uuid.New().String(),
			InstrukturID: req.InstrukturId,
			KursusID:     req.KursusId,
			Status:       true,
		}

		tanggalParsed, err := time.Parse("2006-01-02", req.Tanggal)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Invalid date format for Tanggal", nil))
		}

		jamMulaiParsed, err := time.Parse("15:04:05", req.JamMulai)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Invalid time format for JamMulai", nil))
		}

		jamAkhirParsed, err := time.Parse("15:04:05", req.JamAkhir)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Invalid time format for JamAkhir", nil))
		}

		newJadwal.JamMulai = time.Date(1970, 1, 1, jamMulaiParsed.Hour(), jamMulaiParsed.Minute(), jamMulaiParsed.Second(), 0, time.UTC)
		newJadwal.JamAkhir = time.Date(1970, 1, 1, jamAkhirParsed.Hour(), jamAkhirParsed.Minute(), jamAkhirParsed.Second(), 0, time.UTC)
		newJadwal.Tanggal = tanggalParsed

		err = h.s.CreateJadwalMengajar(&newJadwal)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		return c.JSON(http.StatusCreated, helper.FormatResponse(true, constant.PostData, nil))
	}
}

func (h *JadwalHandler) EditJadwalMengajar() echo.HandlerFunc {
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
		if !ok || (role != constant.RoleAdmin) {
			return helper.UnauthorizedError(c)
		}

		id := c.Param("id")
		dataId, err := h.s.GetJadwalMengajarByID(id)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		var req Request
		if err := c.Bind(&req); err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		newJadwal := jadwal.JadwalMengajar{
			ID:           dataId.ID,
			InstrukturID: ifEmptyString(req.InstrukturId, dataId.InstrukturID),
			KursusID:     ifEmptyString(req.KursusId, dataId.KursusID),
			Status:       req.Status,
		}

		if req.Tanggal != "" {
			tanggalParsed, err := time.Parse("2006-01-02", req.Tanggal)
			if err != nil {
				return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Invalid date format for Tanggal", nil))
			}
			newJadwal.Tanggal = tanggalParsed
		} else {
			newJadwal.Tanggal = dataId.Tanggal
		}

		if req.JamMulai != "" {
			jamMulaiParsed, err := time.Parse("15:04:05", req.JamMulai)
			if err != nil {
				return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Invalid time format for JamMulai", nil))
			}
			newJadwal.JamMulai = time.Date(1970, 1, 1, jamMulaiParsed.Hour(), jamMulaiParsed.Minute(), jamMulaiParsed.Second(), 0, time.UTC)
		} else {
			newJadwal.JamMulai = dataId.JamMulai
		}

		if req.JamAkhir != "" {
			jamAkhirParsed, err := time.Parse("15:04:05", req.JamAkhir)
			if err != nil {
				return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Invalid time format for JamAkhir", nil))
			}
			newJadwal.JamAkhir = time.Date(1970, 1, 1, jamAkhirParsed.Hour(), jamAkhirParsed.Minute(), jamAkhirParsed.Second(), 0, time.UTC)
		} else {
			newJadwal.JamAkhir = dataId.JamAkhir
		}

		err = h.s.EditJadwalMengajar(&newJadwal)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		return c.JSON(http.StatusCreated, helper.FormatResponse(true, constant.PostData, nil))
	}
}

func ifEmptyString(new, old string) string {
	if new == "" {
		return old
	}
	return new
}
