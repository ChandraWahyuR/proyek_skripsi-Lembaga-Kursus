package handler

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"skripsi/constant"
	"skripsi/features/users"
	"skripsi/helper"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	s users.UserServiceInterface
	j helper.JWTInterface
	// redis *redis.Client
}

func New(u users.UserServiceInterface, j helper.JWTInterface) users.UserHandlerInterface {
	return &UserHandler{
		s: u,
		j: j,
		// redis: redis,
	}
}

func (h *UserHandler) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		var reqRegister UserRegisterRequest

		err := c.Bind(&reqRegister)
		if err != nil {
			err, message := helper.HandleEchoError(err)
			return c.JSON(err, helper.FormatResponse(false, message, nil))
		}
		user := users.User{
			Username:        reqRegister.Username,
			Email:           reqRegister.Email,
			NomorHP:         reqRegister.NomorHP,
			Password:        reqRegister.Password,
			ConfirmPassword: reqRegister.ConfirmPassword,
		}
		err = h.s.Register(user)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		token, err := h.j.GenerateVerifikasiEmailJWT(helper.UserJWT{
			Email: user.Email,
		})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, "Failed to generate verification token", nil))
		}

		// Deploy nya dihapus
		link := fmt.Sprintf("%s/verify?token=%s", os.Getenv("CLOUD_RUN_ENDPOINT"), token)
		err = h.s.SendVerificationEmail(user.Email, link)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, "Failed to send verification email", nil))
		}

		return c.JSON(http.StatusCreated, helper.FormatResponse(true, constant.RegisterBerhasil, nil))
	}
}

func (h *UserHandler) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		var reqLogin UserLoginRequest

		err := c.Bind(&reqLogin)
		if err != nil {
			err, message := helper.HandleEchoError(err)
			return c.JSON(err, helper.FormatResponse(false, message, nil))
		}

		user := users.User{
			Email:    reqLogin.Email,
			Password: reqLogin.Password,
		}

		userData, err := h.s.Login(user)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
		}

		var response UserLoginResponse
		response.Token = userData.Token

		return c.JSON(http.StatusOK, helper.FormatResponse(true, constant.LoginBerhasil, response))
	}
}

func (h *UserHandler) ForgotPassword() echo.HandlerFunc {
	return func(c echo.Context) error {
		var reqForgotPassword ForgotPasswordRequest

		err := c.Bind(&reqForgotPassword)
		if err != nil {
			err, message := helper.HandleEchoError(err)
			return c.JSON(err, helper.FormatResponse(false, message, nil))
		}

		user := users.User{
			Email: reqForgotPassword.Email,
		}

		token, err := h.s.ForgotPassword(user)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse(true, constant.KirimLupaPassword, map[string]string{"token": token}))
	}
}

func (h *UserHandler) VerifyOTP() echo.HandlerFunc {
	return func(c echo.Context) error {
		var reqVerifyOTP VerifyOtpRequest

		err := c.Bind(&reqVerifyOTP)
		if err != nil {
			err, message := helper.HandleEchoError(err)
			return c.JSON(err, helper.FormatResponse(false, message, nil))
		}

		tokenString := c.Request().Header.Get("Authorization")
		ctx := c.Request().Context()
		token, err := h.j.ValidateToken(ctx, tokenString)
		if err != nil {
			helper.UnauthorizedError(c)
		}

		extract := h.j.ExtractUserToken(token)
		email, ok := extract[constant.JWT_EMAIL].(string)
		if !ok || email == "" {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Invalid token: Email not found", nil))
		}

		user := users.VerifyOtp{
			Email:  email,
			Otp:    reqVerifyOTP.Otp,
			Status: "Success",
		}

		err = h.s.VerifyOTP(user)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}
		return c.JSON(http.StatusOK, helper.FormatResponse(true, constant.OtpStatus, nil))
	}
}

func (h *UserHandler) ResetPassword() echo.HandlerFunc {
	return func(c echo.Context) error {
		var reqResetPassword ResetPasswordRequest

		err := c.Bind(&reqResetPassword)
		if err != nil {
			err, message := helper.HandleEchoError(err)
			return c.JSON(err, helper.FormatResponse(false, message, nil))
		}

		// Pastikan data di-bind dengan benar
		if reqResetPassword.Password == "" || reqResetPassword.ConfirmationPassword == "" {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "password, confirm password cannot be empty", nil))
		}

		tokenString := c.Request().Header.Get("Authorization")
		ctx := c.Request().Context()
		token, err := h.j.ValidateToken(ctx, tokenString)
		if err != nil {
			helper.UnauthorizedError(c)
		}

		extract := h.j.ExtractUserToken(token)
		email, ok := extract[constant.JWT_EMAIL].(string)
		if !ok || email == "" {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Invalid token: Email not found", nil))
		}

		user := users.ResetPassword{
			Email:                email,
			Password:             reqResetPassword.Password,
			ConfirmationPassword: reqResetPassword.ConfirmationPassword,
		}

		err = h.s.ResetPassword(user)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse(true, constant.ResetPassword, nil))
	}
}

func (h *UserHandler) VerifyAccount() echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.QueryParam("token")
		if tokenString == "" {
			fmt.Println("Token is empty")
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Token is missing", nil))
		}

		fmt.Println("Received token:", tokenString)

		token, err := h.j.ValidateEmailToken(tokenString)
		if err != nil {
			fmt.Println("Failed to validate token:", err)
			return helper.UnauthorizedError(c)
		}

		fmt.Println("Token validated successfully")

		extract := h.j.ExtractUserToken(token)
		email, ok := extract[constant.JWT_EMAIL].(string)
		if !ok || email == "" {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Invalid token: Email not found", nil))
		}
		fmt.Println("Email extracted from token:", email)
		err = h.s.ActivateAccount(email)
		if err != nil {
			fmt.Println("Failed to activate account:", err)
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, "Failed to activate account", nil))
		}
		fmt.Println("Redirecting to verification success page...")
		link := os.Getenv("CLOUD_RUN_ENDPOINT")
		res := link + "/assets/verifikasi_berhasil.html"
		return c.Redirect(http.StatusTemporaryRedirect, res)

	}
}

func (h *UserHandler) GetAllUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
		if tokenString == "" {
			return helper.UnauthorizedError(c)
		}
		ctx := c.Request().Context()
		token, err := h.j.ValidateToken(ctx, tokenString)
		if err != nil {
			helper.UnauthorizedError(c)
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

		data, totalPages, err := h.s.GetAllUserPagination(page, limit)
		metadata := MetadataResponse{
			TotalPage: totalPages,
			Page:      page,
		}

		var dataResponse []GetAllUserResponse
		for _, value := range data {
			formattedTanggalLahir := ""
			if !value.TanggalLahir.IsZero() {
				formattedTanggalLahir = value.TanggalLahir.Format(constant.TTLFormat)
			}
			dataResponse = append(dataResponse, GetAllUserResponse{
				ID:            value.ID,
				NIS:           value.NIS,
				Username:      value.Username,
				Nama:          value.Nama,
				Email:         value.Email,
				NomorHP:       value.NomorHP,
				ProfileUrl:    value.ProfileUrl,
				IsActive:      value.IsActive,
				Agama:         value.Agama,
				Gender:        value.Gender,
				TempatLahir:   value.TempatLahir,
				TanggalLahir:  formattedTanggalLahir,
				OrangTua:      value.OrangTua,
				Profesi:       value.Profesi,
				Ijazah:        value.Ijazah,
				KTP:           value.KTP,
				KartuKeluarga: value.KartuKeluarga,
			})
		}
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}
		return c.JSON(http.StatusOK, helper.MetadataFormatResponse(true, constant.GetAllUser, metadata, dataResponse))
	}
}

func (h *UserHandler) GetUserByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Token
		tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
		if tokenString == "" {
			helper.UnauthorizedError(c)
		}
		ctx := c.Request().Context()
		token, err := h.j.ValidateToken(ctx, tokenString)
		if err != nil {
			helper.UnauthorizedError(c)
		}

		adminData := h.j.ExtractAdminToken(token)
		role, ok := adminData[constant.JWT_ROLE]
		if !ok || role != constant.RoleAdmin {
			return helper.UnauthorizedError(c)
		}

		id := c.Param("id")
		dataUser, err := h.s.GetUserByID(id)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		formattedTanggalLahir := ""
		if !dataUser.TanggalLahir.IsZero() {
			formattedTanggalLahir = dataUser.TanggalLahir.Format(constant.TTLFormat)
		}

		responseData := GetUserIDResponse{
			ID:            dataUser.ID,
			NIS:           dataUser.NIS,
			Username:      dataUser.Username,
			Nama:          dataUser.Nama,
			Email:         dataUser.Email,
			NomorHP:       dataUser.NomorHP,
			Agama:         dataUser.Agama,
			Gender:        dataUser.Gender,
			TempatLahir:   dataUser.TempatLahir,
			TanggalLahir:  formattedTanggalLahir,
			OrangTua:      dataUser.OrangTua,
			Profesi:       dataUser.Profesi,
			Ijazah:        dataUser.Ijazah,
			KTP:           dataUser.KTP,
			KartuKeluarga: dataUser.KartuKeluarga,
			ProfileUrl:    dataUser.ProfileUrl,
		}
		return c.JSON(http.StatusOK, helper.FormatResponse(true, constant.GetAllUser, responseData))
	}
}

func (h *UserHandler) GetUserByUser() echo.HandlerFunc {
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

		userData := h.j.ExtractUserToken(token)
		role, ok := userData[constant.JWT_ROLE]
		userId := userData[constant.JWT_ID].(string)
		if !ok || role != constant.RoleUser {
			return helper.UnauthorizedError(c)
		}

		data, err := h.s.GetUserByID(userId)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		formattedTanggalLahir := ""
		if !data.TanggalLahir.IsZero() {
			formattedTanggalLahir = data.TanggalLahir.Format(constant.TTLFormat)
		}

		responseData := GetUserIDResponse{
			ID:            data.ID,
			NIS:           data.NIS,
			Username:      data.Username,
			Nama:          data.Nama,
			Email:         data.Email,
			NomorHP:       data.NomorHP,
			Agama:         data.Agama,
			Gender:        data.Gender,
			TempatLahir:   data.TempatLahir,
			TanggalLahir:  formattedTanggalLahir,
			OrangTua:      data.OrangTua,
			Profesi:       data.Profesi,
			Ijazah:        data.Ijazah,
			KTP:           data.KTP,
			KartuKeluarga: data.KartuKeluarga,
			ProfileUrl:    data.ProfileUrl,
		}
		return c.JSON(http.StatusOK, helper.FormatResponse(true, constant.GetProfile, responseData))
	}
}

func (h *UserHandler) UpdateUser() echo.HandlerFunc {
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

		userData := h.j.ExtractUserToken(token)
		role, ok := userData[constant.JWT_ROLE]
		userId := userData[constant.JWT_ID].(string)
		if !ok || role != constant.RoleUser {
			return helper.UnauthorizedError(c)
		}

		data, err := h.s.GetUserByID(userId)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		var dataRequest EditUserRequest
		err = c.Bind(&dataRequest)
		if err != nil {
			code, message := helper.HandleEchoError(err)
			return c.JSON(code, helper.FormatResponse(false, message, nil))
		}

		var profileUrl, ktpUrl, kkUrl, ijazahUrl string

		// Handle Profile Image
		if file, err := c.FormFile("profile_url"); err == nil {
			profileUrl, err = handleImageUpload(file, helper.UploadPathUser)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, err.Error(), nil))
			}
		} else {
			profileUrl = data.ProfileUrl
		}

		// KTP
		if file, err := c.FormFile("ktp"); err == nil {
			ktpUrl, err = handleImageUpload(file, helper.UploadPathUser+"ktp/")
			if err != nil {
				return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, err.Error(), nil))
			}
		} else {
			ktpUrl = data.KTP
		}

		// KK
		if file, err := c.FormFile("kartu_keluarga"); err == nil {
			kkUrl, err = handleImageUpload(file, helper.UploadPathUser+"kk/")
			if err != nil {
				return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, err.Error(), nil))
			}
		} else {
			kkUrl = data.KartuKeluarga
		}

		// Ijazah
		if file, err := c.FormFile("ijazah"); err == nil {
			ijazahUrl, err = handleImageUpload(file, helper.UploadPathUser+"ijazah/")
			if err != nil {
				return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, err.Error(), nil))
			}
		} else {
			ijazahUrl = data.Ijazah
		}

		tanggalLahir := data.TanggalLahir // Default value dari database
		if dataRequest.TanggalLahir != "" {
			parsedDate, err := time.Parse(constant.TTLFormat, dataRequest.TanggalLahir)
			if err != nil {
				return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Format tanggal tidak valid. Gunakan format: YYYY-MM-DD.", nil))
			}
			tanggalLahir = parsedDate
		}

		responseData := users.EditUser{
			ID:            data.ID,
			Nama:          dataRequest.Nama,
			NomorHP:       dataRequest.NomorHP,
			Agama:         dataRequest.Agama,
			Gender:        dataRequest.Gender,
			TempatLahir:   dataRequest.TempatLahir,
			TanggalLahir:  tanggalLahir,
			OrangTua:      dataRequest.OrangTua,
			Profesi:       dataRequest.Profesi,
			Ijazah:        ijazahUrl,
			KTP:           ktpUrl,
			KartuKeluarga: kkUrl,
			ProfileUrl:    profileUrl,
			Password:      dataRequest.Password,
		}
		err = h.s.UpdateUser(responseData)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		return c.JSON(http.StatusOK, helper.ObjectFormatResponse(true, constant.EditProfile, nil))
	}
}

func (h *UserHandler) DeleteUser() echo.HandlerFunc {
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

		userData := h.j.ExtractUserToken(token)
		role, ok := userData[constant.JWT_ROLE]
		userId := userData[constant.JWT_ID].(string)
		if !ok || role != constant.RoleUser {
			return helper.UnauthorizedError(c)
		}

		err = h.s.DeleteUser(userId)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse(true, constant.DeleteUser, nil))
	}
}

func (h *UserHandler) SearchUserByUsernameEmail() echo.HandlerFunc {
	return func(c echo.Context) error {
		searchData := c.QueryParam("user")
		tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
		if tokenString == "" {
			return helper.UnauthorizedError(c)
		}

		ctx := c.Request().Context()
		token, err := h.j.ValidateToken(ctx, tokenString)
		if err != nil {
			return helper.UnauthorizedError(c)
		}

		userData := h.j.ExtractUserToken(token)
		role, ok := userData[constant.JWT_ROLE]
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
		dataUser, totalPages, err := h.s.SearchUserByUsernameEmail(searchData, page, limit)
		metadata := MetadataResponse{
			TotalPage: totalPages,
			Page:      page,
		}
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		var dataResponse []GetAllUserResponse
		for _, value := range dataUser {
			formattedTanggalLahir := ""
			if !value.TanggalLahir.IsZero() {
				formattedTanggalLahir = value.TanggalLahir.Format(constant.TTLFormat)
			}
			dataResponse = append(dataResponse, GetAllUserResponse{
				ID:            value.ID,
				NIS:           value.NIS,
				Username:      value.Username,
				Nama:          value.Nama,
				Email:         value.Email,
				NomorHP:       value.NomorHP,
				ProfileUrl:    value.ProfileUrl,
				IsActive:      value.IsActive,
				Agama:         value.Agama,
				Gender:        value.Gender,
				TempatLahir:   value.TempatLahir,
				TanggalLahir:  formattedTanggalLahir,
				OrangTua:      value.OrangTua,
				Profesi:       value.Profesi,
				Ijazah:        value.Ijazah,
				KTP:           value.KTP,
				KartuKeluarga: value.KartuKeluarga,
			})
		}
		return c.JSON(http.StatusOK, helper.MetadataFormatResponse(true, constant.GetAllHistoryTransaki, metadata, dataResponse))
	}
}

// redis ini bayar tambahan, buat ada ada aja sih
// func (h *UserHandler) Logout() echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
// 		if tokenString == "" {
// 			return helper.UnauthorizedError(c)
// 		}

// 		// Validasi token
// 		ctx := c.Request().Context()
// 		token, err := h.j.ValidateToken(ctx, tokenString)
// 		if err != nil {
// 			return helper.UnauthorizedError(c)
// 		}

// 		userData := h.j.ExtractUserToken(token)
// 		role, ok := userData[constant.JWT_ROLE]
// 		if !ok || role != constant.RoleUser {
// 			return helper.UnauthorizedError(c)
// 		}

// 		claims, ok := token.Claims.(jwt.MapClaims)
// 		if !ok || !token.Valid {
// 			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, "Token tidak valid", nil))
// 		}

// 		exp, ok := claims["exp"].(float64)
// 		if !ok {
// 			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, "Token tidak memiliki waktu kedaluwarsa", nil))
// 		}

// 		expiresAt := time.Unix(int64(exp), 0)

// 		err = h.redis.Set(ctx, tokenString, "blacklisted", time.Until(expiresAt)).Err()
// 		if err != nil {
// 			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, "Gagal logout", nil))
// 		}

// 		// Berhasil logout
// 		return c.JSON(http.StatusOK, helper.FormatResponse(true, "Berhasil logout", nil))
// 	}
// }

// Tanpa redis hemat, kalau lokal tah aman
func (h *UserHandler) Logout() echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get(constant.HeaderAuthorization)
		if tokenString == "" {
			return helper.UnauthorizedError(c)
		}

		// Validasi token
		ctx := c.Request().Context()
		token, err := h.j.ValidateToken(ctx, tokenString)
		if err != nil {
			return helper.UnauthorizedError(c)
		}

		userData := h.j.ExtractUserToken(token)
		role, ok := userData[constant.JWT_ROLE]
		if !ok || role != constant.RoleUser {
			return helper.UnauthorizedError(c)
		}

		return c.JSON(http.StatusOK, helper.FormatResponse(true, "Berhasil logout. Token telah dihapus pada client side.", nil))
	}
}

func handleImageUpload(file *multipart.FileHeader, uploadPath string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("error opening file: %w", err)
	}
	defer src.Close()

	objectName := fmt.Sprintf("%s_%s", uuid.New().String(), file.Filename)
	if err := helper.Uploader.UploadFile(src, objectName, uploadPath); err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	return fmt.Sprintf("https://storage.googleapis.com/%s/%s%s", helper.Uploader.BucketName, uploadPath, objectName), nil
}
