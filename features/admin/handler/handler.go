package handler

import (
	"fmt"
	"log"
	"net/http"
	"os/user"
	"path"
	"path/filepath"
	"skripsi/constant"
	"skripsi/features/admin"
	"skripsi/helper"
	"time"

	"github.com/labstack/echo/v4"
)

type AdminHandler struct {
	s admin.AdminServiceInterface
	j helper.JWTInterface
}

func New(u admin.AdminServiceInterface, j helper.JWTInterface) admin.AdminHandlerInterface {
	return &AdminHandler{
		s: u,
		j: j,
	}
}
func (h *AdminHandler) RegisterAdmin() echo.HandlerFunc {
	return func(c echo.Context) error {
		var reqRegister AdminRegisterRequest
		err := c.Bind(&reqRegister)
		if err != nil {
			err, message := helper.HandleEchoError(err)
			return c.JSON(err, helper.FormatResponse(false, message, nil))
		}

		admin := admin.Admin{
			Username:        reqRegister.Username,
			Email:           reqRegister.Email,
			Password:        reqRegister.Password,
			ConfirmPassword: reqRegister.ConfirmPassword,
		}
		err = h.s.RegisterAdmin(admin)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}
		return c.JSON(http.StatusCreated, helper.FormatResponse(true, constant.RegisterBerhasil, nil))
	}
}
func (h *AdminHandler) LoginAdmin() echo.HandlerFunc {
	return func(c echo.Context) error {
		var reqLogin AdminLoginRequest
		err := c.Bind(&reqLogin)
		if err != nil {
			err, message := helper.HandleEchoError(err)
			return c.JSON(err, helper.FormatResponse(false, message, nil))
		}

		admin := admin.Admin{
			Username: reqLogin.Username,
			Password: reqLogin.Password,
		}
		adminData, err := h.s.LoginAdmin(admin)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}
		var response AdminLoginResponse
		response.Token = adminData.Token
		return c.JSON(http.StatusOK, helper.FormatResponse(true, constant.LoginBerhasil, response))
	}
}

// Contoh handler untuk laporan dengan filter tanggal
func (h *AdminHandler) DownloadLaporanPembelian() echo.HandlerFunc {
	return func(c echo.Context) error {
		startDate := c.QueryParam("start_date")
		endDate := c.QueryParam("end_date")

		// Validasi input
		if startDate == "" || endDate == "" {
			return c.JSON(http.StatusBadRequest, "Start date and end date are required")
		}

		tglMulai, err := time.Parse(constant.LayoutFormat, startDate)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "Invalid start date format")
		}

		tglAkhir, err := time.Parse(constant.LayoutFormat, endDate)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "Invalid end date format")
		}

		// Tentukan folder Downloads
		downloadsFolder := GetDownloadsFolder()

		// Generate laporan dan simpan di folder Downloads
		filename, err := h.s.DownloadLaporanPembelian(tglMulai, tglAkhir, downloadsFolder)
		if err != nil {
			log.Println("Error generating report:", err)
			return c.JSON(http.StatusInternalServerError, "Error generating report")
		}

		log.Println("File generated:", filename)

		// Set header untuk mengunduh file
		c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", path.Base(filename)))
		c.Response().Header().Set("Content-Type", "text/csv")

		log.Println("Sending file to client:", filename)
		return c.File(filename)
	}
}

func GetDownloadsFolder() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(usr.HomeDir, "Downloads")
}
