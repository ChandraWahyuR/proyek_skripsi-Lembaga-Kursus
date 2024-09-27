package instruktur

import "github.com/labstack/echo/v4"

type Instruktur struct {
	ID     string
	Name   string
	Gender string
	Email  string
	Alamat string
	NoHp   string
}

type UpdateInstruktur struct {
	ID     string
	Name   string
	Gender string
	Email  string
	Alamat string
	NoHp   string
}

type InstrukturHandlerInterface interface {
	GetAllInstruktur() echo.HandlerFunc
	GetAllInstrukturByID() echo.HandlerFunc

	PostInstruktur() echo.HandlerFunc
	UpdateInstruktur() echo.HandlerFunc
	DeleteInstruktur() echo.HandlerFunc

	GetInstruktorByName() echo.HandlerFunc
}

type InstrukturDataInterface interface {
	GetAllInstruktur() ([]Instruktur, error)
	GetAllInstrukturByID(id string) (Instruktur, error)

	PostInstruktur(data Instruktur) error
	UpdateInstruktur(data UpdateInstruktur) error
	DeleteInstruktur(id string) error

	GetInstruktorByName(name string, page int, limit int) ([]Instruktur, int, error)
	GetInstrukturWithPagination(page int, limit int) ([]Instruktur, int, error)
}

type InstrukturServiceInterface interface {
	GetAllInstruktur() ([]Instruktur, error)
	GetAllInstrukturByID(id string) (Instruktur, error)

	PostInstruktur(data Instruktur) error
	UpdateInstruktur(data UpdateInstruktur) error
	DeleteInstruktur(id string) error

	GetInstruktorByName(name string, page int, limit int) ([]Instruktur, int, error)
	GetInstrukturWithPagination(page int, limit int) ([]Instruktur, int, error)
}
