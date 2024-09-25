package admin

import "github.com/labstack/echo/v4"

type Admin struct {
	ID              string
	Username        string
	Email           string
	Password        string
	ConfirmPassword string
}

type Login struct {
	Username string
	Password string
	Token    string
}
type AdminHandlerInterface interface {
	RegisterAdmin() echo.HandlerFunc
	LoginAdmin() echo.HandlerFunc
}

type AdminServiceInterface interface {
	RegisterAdmin(admin Admin) error
	LoginAdmin(admin Admin) (Login, error)
}

type AdminDataInterface interface {
	RegisterAdmin(admin Admin) error
	LoginAdmin(admin Admin) (Admin, error)
	IsEmailExist(email string) bool
	IsUsernameExist(username string) bool
}
