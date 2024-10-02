package handler

import (
	"fmt"
	"log"
	"net/http"
	"skripsi/constant"
	"skripsi/features/kursus"
	"skripsi/helper"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type KursusHandler struct {
	s kursus.KursusServiceInterface
	j helper.JWTInterface
}

func New(u kursus.KursusServiceInterface, j helper.JWTInterface) kursus.KursusHandlerInterface {
	return &KursusHandler{
		s: u,
		j: j,
	}
}

func (h *KursusHandler) GetAllKursus() echo.HandlerFunc {
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

		tokenData := h.j.ExtractUserToken(token)
		role, ok := tokenData[constant.JWT_ROLE]
		// ini intinya jika token tidak ada role atau user tanpa role return error atau jika role bukan admin dan juga bukan user return error
		if !ok || (role != constant.RoleAdmin && role != constant.RoleUser) {
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
		data, totalPages, err := h.s.GetKursusPagination(page, limit)
		metadata := MetadataResponse{
			TotalPage: totalPages,
			Page:      page,
		}

		var response []ResponseGetAllKursus
		for _, responseData := range data {
			var imageResponse []ImageKursus
			for _, img := range responseData.Image {
				imageResponse = append(imageResponse, ImageKursus{
					Name:     img.Name,
					Url:      img.Url,
					Position: img.Position,
				})
			}
			response = append(response, ResponseGetAllKursus{
				ID:        responseData.ID,
				Nama:      responseData.Nama,
				Deskripsi: responseData.Deskripsi,
				Image:     imageResponse,
				Harga:     responseData.Harga,
				Jadwal:    responseData.Jadwal,
			})
		}
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		return c.JSON(http.StatusOK, helper.MetadataFormatResponse(true, "Berhasil", metadata, response))
	}
}
func (h *KursusHandler) GetAllKursusById() echo.HandlerFunc {
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

		tokenData := h.j.ExtractUserToken(token)
		role, ok := tokenData[constant.JWT_ROLE]
		// ini intinya jika token tidak ada role atau user tanpa role return error atau jika role bukan admin dan juga bukan user return error
		if !ok || (role != constant.RoleAdmin && role != constant.RoleUser) {
			return helper.UnauthorizedError(c)
		}

		id := c.Param("id")
		dataKursus, err := h.s.GetAllKursusById(id)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		kursusResponse := ResponseGetKursus{
			ID:        dataKursus.ID,
			Nama:      dataKursus.Nama,
			Deskripsi: dataKursus.Deskripsi,
			Jadwal:    dataKursus.Jadwal,
			Harga:     dataKursus.Harga,
			Instruktur: Instruktur{
				Name: dataKursus.InstruktorID,
			},
			Image:              mapImages(dataKursus.Image),
			Kategori:           mapKategori(dataKursus.Kategori),
			MateriPembelajaran: mapMateri(dataKursus.MateriPembelajaran),
		}
		return c.JSON(http.StatusOK, helper.ObjectFormatResponse(true, "Berhasil", kursusResponse))
	}
}

func (h *KursusHandler) AddKursus() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Token
		tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
		if tokenString == "" {
			return helper.UnauthorizedError(c)
		}
		token, err := h.j.ValidateToken(tokenString)
		if err != nil {
			return helper.UnauthorizedError(c)
		}

		adminData := h.j.ExtractAdminToken(token)
		role, ok := adminData[constant.JWT_ROLE]
		if !ok || role != constant.RoleAdmin {
			return helper.UnauthorizedError(c)
		}

		var kursusRequest KursusRequest
		if err := c.Bind(&kursusRequest); err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Invalid request", nil))
		}

		form, err := c.MultipartForm()
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Failed to get form-data", nil))
		}

		// Inisialisasi kursusData
		kursusData := kursus.Kursus{
			ID:           uuid.New().String(), // UUID unik untuk kursus
			Nama:         kursusRequest.Nama,
			Deskripsi:    kursusRequest.Deskripsi,
			Jadwal:       kursusRequest.Jadwal,
			Harga:        kursusRequest.Harga,
			InstruktorID: kursusRequest.InstruktorID,
		}

		// Upload multiple image files
		files := form.File["image"] // "image" adalah key untuk semua file yang diupload
		for i, file := range files {
			src, err := file.Open()
			if err != nil {
				return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, "Error opening file", nil))
			}
			defer src.Close()

			// Generate unique filename
			objectName := fmt.Sprintf("%s_%s", uuid.New().String(), file.Filename)

			// Upload file to GCS
			err = helper.Uploader.UploadFileGambarKursus(src, objectName)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, "Failed to upload file to GCS", nil))
			}

			// Generate image URL
			imageUrl := fmt.Sprintf("https://storage.googleapis.com/%s/%s%s", helper.Uploader.BucketName, helper.UploadPathKursus, objectName)
			kursusData.Image = append(kursusData.Image, kursus.ImageKursus{
				ID:       uuid.New().String(),
				Name:     kursusData.Nama,
				Url:      imageUrl,
				Position: i + 1,
				KursusID: kursusData.ID,
			})
		}

		// Menambahkan kategori
		for _, kategoriData := range kursusRequest.Kategori {
			kursusData.Kategori = append(kursusData.Kategori, kursus.KategoriKursus{
				ID:         uuid.New().String(),
				KursusID:   kursusData.ID,
				KategoriID: kategoriData,
			})
		}

		// Menambahkan materi pembelajaran
		for i, materiData := range kursusRequest.MateriPembelajaran {
			kursusData.MateriPembelajaran = append(kursusData.MateriPembelajaran, kursus.MateriPembelajaran{
				ID:       uuid.New().String(),
				KursusID: kursusData.ID,
				MateriID: materiData,
				Position: i + 1,
			})
		}

		// Simpan kursus ke database
		err = h.s.AddKursus(kursusData)
		if err != nil {
			// Tambahkan log error untuk tracking
			log.Println("Error inserting kursus data: ", err)
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, err.Error(), nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse(true, "Kursus added successfully", nil))
	}
}

// ============================================

func mapImages(images []kursus.ImageKursus) []ImageKursus {
	var result []ImageKursus
	for _, img := range images {
		result = append(result, ImageKursus{
			Name:     img.Name,
			Url:      img.Url,
			Position: img.Position,
		})
	}
	return result
}

func mapKategori(kategori []kursus.KategoriKursus) []KategoriKursus {
	var result []KategoriKursus
	for _, kat := range kategori {
		result = append(result, KategoriKursus{
			Nama:      kat.Kategori.Nama,
			Deskripsi: kat.Kategori.Deskripsi,
			ImageUrl:  kat.Kategori.ImageUrl,
		})
	}
	return result
}

func mapMateri(materi []kursus.MateriPembelajaran) []MateriPembelajaran {
	var result []MateriPembelajaran
	for _, mat := range materi {
		result = append(result, MateriPembelajaran{
			Name:      mat.Materi.Name,
			Deskripsi: mat.Materi.Deskripsi,
		})
	}
	return result
}
