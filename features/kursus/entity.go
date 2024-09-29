package kursus

import (
	"skripsi/features/instruktur"
	"skripsi/features/kategori"
	"time"

	"github.com/labstack/echo/v4"
)

type Kursus struct {
	ID                 string
	Nama               string
	Image              []ImageKursus
	Deskripsi          string
	Kategori           []KategoriKursus
	Jadwal             time.Time
	Harga              int
	Instruktur         Instruktur
	MateriPembelajaran []MateriPembelajaran
}

type KategoriKursus struct {
	ID         string
	KursusID   string
	KategoriID string
	Kategori   kategori.Kategori
}

type Instruktur struct {
	ID           string
	KursusID     string
	InstrukturID string
	Instruktur   instruktur.Instruktur
}

type MateriPembelajaran struct {
	ID       string
	KursusID string
	MateriID string
	Materi   Materi
}

type Materi struct {
	ID        string
	Name      string
	Deskripsi string
}

type ImageKursus struct {
	ID       string
	Name     string
	Url      string
	Position int
	KursusID string
}

type KursusHandlerInterface interface {
	GetAllKursus() echo.HandlerFunc
	GetAllKursusById() echo.HandlerFunc
	GetAllKursusByName() echo.HandlerFunc

	AddKursus() echo.HandlerFunc
	UpdateKursus() echo.HandlerFunc
	DeleteKursus() echo.HandlerFunc
}
type KursusDataInterface interface {
	GetAllKursus() ([]Kursus, error)
	GetAllKursusById(id string) (Kursus, error)
	GetAllKursusByName(name string) ([]Kursus, error)

	AddKursus(Kursus) error
	UpdateKursus(Kursus) error
	DeleteKursus(id string) error

	GetKursusPagination(page, limit int) ([]Kursus, int, error)
}

type KursusServiceInterface interface {
	GetAllKursus() ([]MateriPembelajaran, error)
	GetAllKursusById(id string) (MateriPembelajaran, error)
	GetAllKursusByName(name string) ([]Kursus, error)

	AddKursus(Kursus) error
	UpdateKursus(Kursus) error
	DeleteKursus(id string) error

	GetKursusPagination(page, limit int) ([]Kursus, int, error)
}
