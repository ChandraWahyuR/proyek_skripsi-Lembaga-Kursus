package handler

import (
	"fmt"
	"log"
	"net/http"
	"skripsi/constant"
	"skripsi/features/kategori"
	"skripsi/helper"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type KategoriHandler struct {
	s kategori.KategoriServiceInterface
	j helper.JWTInterface
}

func New(u kategori.KategoriServiceInterface, j helper.JWTInterface) kategori.KategoriHandlerInterface {
	return &KategoriHandler{
		s: u,
		j: j,
	}
}

func (h *KategoriHandler) GetAllKategori() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Token
		tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
		if tokenString == "" {
			helper.UnauthorizedError(c)
		}
		token, err := h.j.ValidateToken(tokenString)
		if err != nil {
			helper.UnauthorizedError(c)
		}

		adminData := h.j.ExtractAdminToken(token)
		role, ok := adminData[constant.JWT_ROLE]
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
		data, totalPages, err := h.s.GetKategoriWithPagination(page, limit)
		metadata := MetadataResponse{
			TotalPage: totalPages,
			Page:      page,
		}

		var response []KategoriResponse
		for _, v := range data {
			response = append(response, KategoriResponse{
				ID:        v.ID,
				Nama:      v.Nama,
				Deskripsi: v.Deskripsi,
				ImageUrl:  v.ImageUrl,
			})
		}
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		return c.JSON(http.StatusOK, helper.MetadataFormatResponse(true, "Berhasil", metadata, response))
	}
}

func (h *KategoriHandler) GetKategoriById() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Token
		tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
		if tokenString == "" {
			helper.UnauthorizedError(c)
		}
		token, err := h.j.ValidateToken(tokenString)
		if err != nil {
			helper.UnauthorizedError(c)
		}

		adminData := h.j.ExtractAdminToken(token)
		role, ok := adminData[constant.JWT_ROLE]
		if !ok || role != constant.RoleAdmin {
			return helper.UnauthorizedError(c)
		}

		id := c.Param("id")
		dataKategori, err := h.s.GetKategoriById(id)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		dataResponse := KategoriResponse{
			ID:        dataKategori.ID,
			Nama:      dataKategori.Nama,
			Deskripsi: dataKategori.Deskripsi,
			ImageUrl:  dataKategori.ImageUrl,
		}

		return c.JSON(http.StatusOK, helper.ObjectFormatResponse(true, "Berhasil", dataResponse))

	}
}

func (h *KategoriHandler) CreateKategori() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Token
		tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
		if tokenString == "" {
			helper.UnauthorizedError(c)
		}
		token, err := h.j.ValidateToken(tokenString)
		if err != nil {
			helper.UnauthorizedError(c)
		}

		adminData := h.j.ExtractAdminToken(token)
		role, ok := adminData[constant.JWT_ROLE]
		if !ok || role != constant.RoleAdmin {
			return helper.UnauthorizedError(c)
		}

		// Logic
		var dataRequest RequestKategori
		err = c.Bind(&dataRequest)
		if err != nil {
			code, message := helper.HandleEchoError(err)
			return c.JSON(code, helper.FormatResponse(false, message, nil))
		}

		if dataRequest.Nama == "" {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, constant.ErrEmptyNamaKategori.Error(), nil))
		}
		if dataRequest.Deskripsi == "" {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, constant.ErrEmptyDeskripsiKategori.Error(), nil))
		}
		// Validate the existence of image file without uploading it
		file, err := c.FormFile("image")
		if err != nil {
			log.Println("Error getting file:", err)
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Failed to get image", nil))
		}

		// All validations passed, now proceed to upload the image
		src, err := file.Open()
		if err != nil {
			log.Println("Error opening file:", err)
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, "Error opening file", nil))
		}
		defer src.Close()

		// Generate unique filename
		objectName := fmt.Sprintf("%s_%s", uuid.New().String(), file.Filename)

		// Upload file to GCS
		err = helper.Uploader.UploadFileGambarKategori(src, objectName)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, "Failed to upload file to GCS", nil))
		}

		// Generate image URL
		imageUrl := fmt.Sprintf("https://storage.googleapis.com/%s/%s%s", helper.Uploader.BucketName, helper.UploadPathKategori, objectName)

		dataResponse := kategori.Kategori{
			ID:        uuid.New().String(),
			Nama:      dataRequest.Nama,
			Deskripsi: dataRequest.Deskripsi,
			ImageUrl:  imageUrl,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err = h.s.CreateKategori(dataResponse)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		return c.JSON(http.StatusCreated, helper.FormatResponse(true, "Success", nil))
	}
}
func (h *KategoriHandler) UpdateKategori() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Token
		tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
		if tokenString == "" {
			helper.UnauthorizedError(c)
		}
		token, err := h.j.ValidateToken(tokenString)
		if err != nil {
			helper.UnauthorizedError(c)
		}

		adminData := h.j.ExtractAdminToken(token)
		role, ok := adminData[constant.JWT_ROLE]
		if !ok || role != constant.RoleAdmin {
			return helper.UnauthorizedError(c)
		}

		id := c.Param("id")
		dataId, err := h.s.GetKategoriById(id)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		var dataRequest RequestKategori
		err = c.Bind(&dataRequest)
		if err != nil {
			code, message := helper.HandleEchoError(err)
			return c.JSON(code, helper.FormatResponse(false, message, nil))
		}

		file, err := c.FormFile("image")
		var imageUrl string
		if err == nil {
			// Gambar Baru
			src, err := file.Open()
			if err != nil {
				log.Println("Error opening file:", err)
				return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, "Error opening file", nil))
			}
			defer src.Close()

			// Generate unique filename and upload path
			objectName := fmt.Sprintf("%s_%s", uuid.New().String(), file.Filename)
			err = helper.Uploader.UploadFileGambarKategori(src, objectName)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, "Failed to upload file to GCS", nil))
			}
			imageUrl = fmt.Sprintf("https://storage.googleapis.com/%s/%s%s", helper.Uploader.BucketName, helper.UploadPathKategori, objectName)
		} else {
			// Data lama
			imageUrl = dataId.ImageUrl
		}
		dataResponse := kategori.Kategori{
			ID:        dataId.ID,
			Nama:      dataRequest.Nama,
			Deskripsi: dataRequest.Deskripsi,
			ImageUrl:  imageUrl,
		}

		err = h.s.UpdateKategori(dataResponse)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		return c.JSON(http.StatusOK, helper.ObjectFormatResponse(true, "Success", nil))
	}
}
func (h *KategoriHandler) DeleteKategori() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Token
		tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
		if tokenString == "" {
			helper.UnauthorizedError(c)
		}
		token, err := h.j.ValidateToken(tokenString)
		if err != nil {
			helper.UnauthorizedError(c)
		}

		adminData := h.j.ExtractAdminToken(token)
		role, ok := adminData[constant.JWT_ROLE]
		if !ok || role != constant.RoleAdmin {
			return helper.UnauthorizedError(c)
		}

		id := c.Param("id")
		err = h.s.DeleteKategori(id)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse(true, "Success", nil))
	}
}
