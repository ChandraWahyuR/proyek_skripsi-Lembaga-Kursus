package handler

type UserRegisterRequest struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	NomorHP         string `json:"nomor_hp"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type VerifyOtpRequest struct {
	Email string `json:"email"`
	Otp   string `json:"otp"`
}

type ResetPasswordRequest struct {
	Email                string `json:"email"`
	Password             string `json:"password"`
	ConfirmationPassword string `json:"confirmation_password"`
}

type EditUserRequest struct {
	ProfileUrl    string `json:"profile_url" form:"profile_url"`
	Password      string `json:"password" form:"password"`
	NomorHP       string `json:"nomor_hp" form:"nomor_hp"`
	Nama          string `json:"nama" form:"nama"`
	Agama         string `json:"agama" form:"agama"`
	Gender        string `json:"gender" form:"gender"`
	TempatLahir   string `json:"tempat_lahir" form:"tempat_lahir"`
	TanggalLahir  string `json:"tanggal_lahir" form:"tanggal_lahir"`
	OrangTua      string `json:"orangtua" form:"orang_tua"`
	Profesi       string `json:"profesi" form:"profesi"`
	Ijazah        string `json:"ijazah" form:"ijazah"`
	KTP           string `json:"ktp" form:"ktp"`
	KartuKeluarga string `json:"kartu_keluarga" form:"kartu_keluarga"`
}
