package handler

type AdminRegisterRequest struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type AdminLoginRequest struct {
	Username string `json:"username"  `
	Password string `json:"password"`
}
