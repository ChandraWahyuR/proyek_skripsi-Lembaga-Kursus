package handler

import (
	"fmt"
	"net/http"
	"skripsi/constant"
	"skripsi/features/kursus"
	"skripsi/helper"
	"strconv"
	"time"

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
		ctx := c.Request().Context()
		token, err := h.j.ValidateToken(ctx, tokenString)
		if err != nil {
			return helper.UnauthorizedError(c)
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
				Jadwal:    mapJadwal(responseData.Jadwal),
			})
		}
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		return c.JSON(http.StatusOK, helper.MetadataFormatResponse(true, constant.GetAllKursus, metadata, response))
	}
}
func (h *KursusHandler) GetAllKursusById() echo.HandlerFunc {
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
			Jadwal:    mapJadwal(dataKursus.Jadwal),
			Harga:     dataKursus.Harga,
			Instruktur: Instruktur{
				ID:   dataKursus.Instruktur.ID,
				Name: dataKursus.Instruktur.Name,
			},
			Image:              mapImages(dataKursus.Image),
			Kategori:           mapKategori(dataKursus.Kategori),
			MateriPembelajaran: mapMateri(dataKursus.MateriPembelajaran),
		}
		return c.JSON(http.StatusOK, helper.ObjectFormatResponse(true, constant.GetAllKursus, kursusResponse))
	}
}

func (h *KursusHandler) AddKursus() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Token
		tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
		if tokenString == "" {
			return helper.UnauthorizedError(c)
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

		var kursusRequest KursusRequest
		if err := c.Bind(&kursusRequest); err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
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
			Harga:        kursusRequest.Harga,
			InstrukturID: kursusRequest.InstrukturID,
		}

		// Upload multiple image files
		files := form.File["image"]
		for i, file := range files {
			src, err := file.Open()
			if err != nil {
				return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, constant.ErrOpeningFile.Error(), nil))
			}
			defer src.Close()

			// Generate unique filename
			objectName := fmt.Sprintf("%s_%s", uuid.New().String(), file.Filename)

			// Upload file to GCS
			err = helper.Uploader.UploadFileGambarKursus(src, objectName)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, constant.ErrUploadGCS.Error(), nil))
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
				ID:        uuid.New().String(),
				KursusID:  kursusData.ID,
				Position:  i + 1,
				Deskripsi: materiData,
			})
		}

		//
		i := 0
		for {
			hari := c.FormValue(fmt.Sprintf("jadwal[%d][hari]", i))
			jamMulai := c.FormValue(fmt.Sprintf("jadwal[%d][jam_mulai]", i))
			jamSelesai := c.FormValue(fmt.Sprintf("jadwal[%d][jam_selesai]", i))

			if hari == "" || jamMulai == "" || jamSelesai == "" {
				break
			}

			// Conver waktu
			jamMulaiParsed, err := helper.ValidateTime(jamMulai)
			if err != nil {
				return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Invalid jam_mulai format", nil))
			}

			jamSelesaiParsed, err := helper.ValidateTime(jamSelesai)
			if err != nil {
				return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Invalid jam_selesai format", nil))
			}

			kursusData.Jadwal = append(kursusData.Jadwal, kursus.JadwalKursus{
				ID:         uuid.New().String(),
				KursusID:   kursusData.ID,
				Hari:       hari,
				JamMulai:   jamMulaiParsed,
				JamSelesai: jamSelesaiParsed,
			})

			i++
		}

		// Simpan kursus ke database
		err = h.s.AddKursus(kursusData)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, err.Error(), nil))
		}

		return c.JSON(http.StatusCreated, helper.FormatResponse(true, constant.PostKursus, nil))
	}
}

func (h *KursusHandler) UpdateKursus() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Token
		tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
		if tokenString == "" {
			return helper.UnauthorizedError(c)
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

		id := c.Param("id")
		dataKursus, err := h.s.GetAllKursusById(id)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		var kursusRequest KursusRequest
		err = c.Bind(&kursusRequest)
		if err != nil {
			code, message := helper.HandleEchoError(err)
			return c.JSON(code, helper.FormatResponse(false, message, nil))
		}

		form, err := c.MultipartForm()
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Failed to get form-data", nil))
		}

		// Data kursus yang akan di-update
		responseData := kursus.Kursus{
			ID:           dataKursus.ID,
			Nama:         kursusRequest.Nama,
			Deskripsi:    kursusRequest.Deskripsi,
			InstrukturID: kursusRequest.InstrukturID,
			Harga:        kursusRequest.Harga,
			UpdatedAt:    time.Now(),
		}

		// Add new materials
		for i, materi := range kursusRequest.MateriPembelajaran {
			responseData.MateriPembelajaran = append(responseData.MateriPembelajaran, kursus.MateriPembelajaran{
				ID:        uuid.New().String(),
				Deskripsi: materi,
				KursusID:  dataKursus.ID,
				Position:  i + 1,
			})
		}

		i := 0
		for {
			hari := c.FormValue(fmt.Sprintf("jadwal[%d][hari]", i))
			jamMulai := c.FormValue(fmt.Sprintf("jadwal[%d][jam_mulai]", i))
			jamSelesai := c.FormValue(fmt.Sprintf("jadwal[%d][jam_selesai]", i))

			if hari == "" || jamMulai == "" || jamSelesai == "" {
				break
			}

			// Conver waktu
			jamMulaiParsed, err := helper.ValidateTime(jamMulai)
			if err != nil {
				return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Invalid jam_mulai format", nil))
			}

			jamSelesaiParsed, err := helper.ValidateTime(jamSelesai)
			if err != nil {
				return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Invalid jam_selesai format", nil))
			}

			responseData.Jadwal = append(responseData.Jadwal, kursus.JadwalKursus{
				ID:         uuid.New().String(),
				KursusID:   responseData.ID,
				Hari:       hari,
				JamMulai:   jamMulaiParsed,
				JamSelesai: jamSelesaiParsed,
			})

			i++
		}

		// Add new categories
		for _, dataKategori := range kursusRequest.Kategori {
			responseData.Kategori = append(responseData.Kategori, kursus.KategoriKursus{
				ID:         uuid.New().String(),
				KategoriID: dataKategori,
				KursusID:   dataKursus.ID,
			})
		}

		// Add new images
		files := form.File["image"]
		for i, file := range files {
			src, err := file.Open()
			if err != nil {
				return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, constant.ErrOpeningFile.Error(), nil))
			}
			defer src.Close()

			objectName := fmt.Sprintf("%s_%s", uuid.New().String(), file.Filename)
			err = helper.Uploader.UploadFileGambarKursus(src, objectName)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, constant.ErrUploadGCS.Error(), nil))
			}

			imageUrl := fmt.Sprintf("https://storage.googleapis.com/%s/%s%s", helper.Uploader.BucketName, helper.UploadPathKursus, objectName)
			responseData.Image = append(responseData.Image, kursus.ImageKursus{
				ID:       uuid.New().String(),
				Name:     kursusRequest.Nama,
				Url:      imageUrl,
				Position: i + 1,
				KursusID: dataKursus.ID,
			})
		}

		// Update kursus data
		err = h.s.UpdateKursus(responseData)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}
		return c.JSON(http.StatusOK, helper.FormatResponse(true, constant.EditKursus, nil))
	}
}

func (h *KursusHandler) DeleteKursus() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Token
		tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
		if tokenString == "" {
			return helper.UnauthorizedError(c)
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
		id := c.Param("id")
		err = h.s.DeleteKursus(id)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse(true, constant.DeleteKursus, nil))
	}
}

func (h *KursusHandler) GetAllKursusByName() echo.HandlerFunc {
	return func(c echo.Context) error {
		namaKursus := c.QueryParam("name")
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
		result, totalpage, err := h.s.GetAllKursusByName(namaKursus, page, limit)
		metaData := MetadataResponse{
			Page:      page,
			TotalPage: totalpage,
		}
		var response []ResponseGetAllKursus
		for _, data := range result {
			response = append(response, ResponseGetAllKursus{
				ID:        data.ID,
				Nama:      data.Nama,
				Deskripsi: data.Deskripsi,
				Image:     mapImages(data.Image),
				Jadwal:    mapJadwal(data.Jadwal),
				Harga:     data.Harga,
			})
		}
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		return c.JSON(http.StatusOK, helper.MetadataFormatResponse(true, constant.GetAllKursus, metaData, response))
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
			Deskripsi: mat.Deskripsi,
		})
	}
	return result
}

func mapJadwal(jadwal []kursus.JadwalKursus) []JadwalKursus {
	var result []JadwalKursus
	for _, j := range jadwal {
		result = append(result, JadwalKursus{
			Hari:       j.Hari,
			JamMulai:   j.JamMulai.Format("15:04"),
			JamSelesai: j.JamSelesai.Format("15:04"),
		})
	}
	return result
}
