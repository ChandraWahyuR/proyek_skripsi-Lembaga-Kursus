package instruktur

import (
	"time"

	"github.com/labstack/echo/v4"
)

type Instruktur struct {
	ID                   string
	NomorIndukPendidikan string
	NIK                  string
	UrlImage             string
	Name                 string
	Gender               string
	Email                string
	Alamat               string
	NoHp                 string
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type UpdateInstruktur struct {
	ID                   string
	NomorIndukPendidikan string
	NIK                  string
	UrlImage             string
	Name                 string
	Gender               string
	Email                string
	Alamat               string
	NoHp                 string
	UpdatedAt            time.Time
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
