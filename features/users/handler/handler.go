package handler

import (
	"fmt"
	"net/http"
	"skripsi/constant"
	"skripsi/features/users"
	"skripsi/helper"
	"strconv"

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
		link := "https://skripsi-245802795341.asia-southeast2.run.app/verify?token=" + token
		err = h.s.SendVerificationEmail(user.Email, link)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, "Failed to send verification email", nil))
		}

		return c.JSON(http.StatusCreated, helper.FormatResponse(true, "Success", nil))
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

		return c.JSON(http.StatusOK, helper.FormatResponse(true, "Login Success", response))
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

		return c.JSON(http.StatusOK, helper.FormatResponse(true, "Success", map[string]string{"token": token}))
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
		return c.JSON(http.StatusOK, helper.FormatResponse(true, "Success", nil))
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

		return c.JSON(http.StatusOK, helper.FormatResponse(true, "Success", nil))
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
		// return c.Redirect(http.StatusTemporaryRedirect, "assets/verification-success")
		return c.Redirect(http.StatusTemporaryRedirect, "https://skripsi-245802795341.asia-southeast2.run.app/assets/verifikasi_berhasil.html")

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
			dataResponse = append(dataResponse, GetAllUserResponse{
				ID:       value.ID,
				NIS:      value.NIS,
				Username: value.Username,
				Nama:     value.Nama,
				Email:    value.Email,
				NomorHP:  value.NomorHP,
			})
		}
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}
		return c.JSON(http.StatusOK, helper.MetadataFormatResponse(true, "Success", metadata, dataResponse))
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
			TanggalLahir:  dataUser.TanggalLahir,
			OrangTua:      dataUser.OrangTua,
			Profesi:       dataUser.Profesi,
			Ijazah:        dataUser.Ijazah,
			KTP:           dataUser.KTP,
			KartuKeluarga: dataUser.KartuKeluarga,
			ProfileUrl:    dataUser.ProfileUrl,
		}
		return c.JSON(http.StatusOK, helper.FormatResponse(true, "Success", responseData))
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
			TanggalLahir:  data.TanggalLahir,
			OrangTua:      data.OrangTua,
			Profesi:       data.Profesi,
			Ijazah:        data.Ijazah,
			KTP:           data.KTP,
			KartuKeluarga: data.KartuKeluarga,
			ProfileUrl:    data.ProfileUrl,
		}
		return c.JSON(http.StatusOK, helper.FormatResponse(true, "Success", responseData))
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

		file, err := c.FormFile("image")
		var imageUrl string
		if err == nil {
			// Gambar Baru
			src, err := file.Open()
			if err != nil {
				return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, "Error opening file", nil))
			}
			defer src.Close()

			objectName := fmt.Sprintf("%s_%s", uuid.New().String(), file.Filename)
			err = helper.Uploader.UploadFileGambarUser(src, objectName)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, "Failed to upload file to GCS", nil))
			}
			imageUrl = fmt.Sprintf("https://storage.googleapis.com/%s/%s%s", helper.Uploader.BucketName, helper.UploadPathKategori, objectName)
		} else {
			// Data lama
			imageUrl = data.ProfileUrl
		}
		responseData := users.EditUser{
			ID:            data.ID,
			Nama:          dataRequest.Nama,
			NomorHP:       dataRequest.NomorHP,
			Agama:         dataRequest.Agama,
			Gender:        dataRequest.Gender,
			TempatLahir:   dataRequest.TempatLahir,
			TanggalLahir:  dataRequest.TanggalLahir,
			OrangTua:      dataRequest.OrangTua,
			Profesi:       dataRequest.Profesi,
			Ijazah:        dataRequest.Ijazah,
			KTP:           dataRequest.KTP,
			KartuKeluarga: dataRequest.KartuKeluarga,
			ProfileUrl:    imageUrl,
		}
		err = h.s.UpdateUser(responseData)
		if err != nil {
			return c.JSON(helper.ConverResponse(err), helper.FormatResponse(false, err.Error(), nil))
		}

		return c.JSON(http.StatusOK, helper.ObjectFormatResponse(true, "Success", nil))
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

		return c.JSON(http.StatusOK, helper.FormatResponse(true, "Success", nil))
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
