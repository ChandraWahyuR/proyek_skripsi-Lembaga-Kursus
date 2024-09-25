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
