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
	Deskripsi          string
	Harga              int
	InstrukturID       string
	Instruktur         instruktur.Instruktur
	Jadwal             []JadwalKursus
	Kategori           []KategoriKursus
	Image              []ImageKursus
	MateriPembelajaran []MateriPembelajaran
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
type JadwalKursus struct {
	ID         string
	KursusID   string
	Kursus     Kursus
	Hari       string
	JamMulai   time.Time
	JamSelesai time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
type KategoriKursus struct {
	ID         string
	KursusID   string
	KategoriID string
	Kategori   kategori.Kategori
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type MateriPembelajaran struct {
	ID        string
	KursusID  string
	Position  int
	Deskripsi string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ImageKursus struct {
	ID        string
	Name      string
	Url       string
	Position  int
	KursusID  string
	CreatedAt time.Time
	UpdatedAt time.Time
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
	GetAllKursusByName(name string, page int, limit int) ([]Kursus, int, error)

	AddKursus(Kursus) error
	UpdateKursus(Kursus) error
	DeleteKursus(id string) error

	GetKursusPagination(page, limit int) ([]Kursus, int, error)
	//
	DeleteImageKursus(id string) error
	DeleteMateriKursus(id string) error
	DeleteKategoriKursus(id string) error
}

type KursusServiceInterface interface {
	GetAllKursus() ([]Kursus, error)
	GetAllKursusById(id string) (Kursus, error)
	GetAllKursusByName(name string, page int, limit int) ([]Kursus, int, error)

	AddKursus(Kursus) error
	UpdateKursus(Kursus) error
	DeleteKursus(id string) error

	GetKursusPagination(page, limit int) ([]Kursus, int, error)
	//
	DeleteImageKursus(id string) error
	DeleteMateriKursus(id string) error
	DeleteKategoriKursus(id string) error
}
