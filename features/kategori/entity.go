package kategori

import "github.com/labstack/echo/v4"

type Kategori struct {
	ID        string
	Nama      string
	Deskripsi string
	ImageUrl  string
}

type KategoriHandlerInterface interface {
	GetAllKategori() echo.HandlerFunc
	GetKategoriById() echo.HandlerFunc

	CreateKategori() echo.HandlerFunc
	UpdateKategori() echo.HandlerFunc
	DeleteKategori() echo.HandlerFunc
}

type KategoriDataInterface interface {
	GetAllKategori() ([]Kategori, error)
	GetKategoriById(id string) (Kategori, error)

	CreateKategori(Kategori) error
	UpdateKategori(Kategori) error
	DeleteKategori(id string) error

	GetKategoriWithPagination(page int, limit int) ([]Kategori, int, error)
}
type KategoriServiceInterface interface {
	GetAllKategori() ([]Kategori, error)
	GetKategoriById(id string) (Kategori, error)

	CreateKategori(Kategori) error
	UpdateKategori(Kategori) error
	DeleteKategori(id string) error

	GetKategoriWithPagination(page int, limit int) ([]Kategori, int, error)
}
