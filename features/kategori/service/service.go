package service

import (
	"skripsi/constant"
	"skripsi/features/kategori"
	"skripsi/helper"
)

type KategoriService struct {
	d kategori.KategoriDataInterface
	j helper.JWTInterface
}

func New(u kategori.KategoriDataInterface, j helper.JWTInterface) kategori.KategoriServiceInterface {
	return &KategoriService{
		d: u,
		j: j,
	}
}

func (s KategoriService) GetAllKategori() ([]kategori.Kategori, error) {
	return s.d.GetAllKategori()

}
func (s KategoriService) GetKategoriById(id string) (kategori.Kategori, error) {
	return s.d.GetKategoriById(id)
}
func (s KategoriService) CreateKategori(data kategori.Kategori) error {
	switch {
	case data.Nama == "":
		return constant.ErrEmptyNamaKategori
	case data.Deskripsi == "":
		return constant.ErrEmptyDeskripsiKategori
	case data.ImageUrl == "":
		return constant.ErrEmptyImageUrlKategori
	}

	return s.d.CreateKategori(data)
}
func (s KategoriService) UpdateKategori(data kategori.Kategori) error {
	if data.ID == "" {
		return constant.ErrEmptyId
	}

	if data.Nama == "" && data.Deskripsi == "" && data.ImageUrl == "" {
		return constant.ErrUpdate
	}

	return s.d.UpdateKategori(data)
}
func (s KategoriService) DeleteKategori(id string) error {
	if id == "" {
		return constant.ErrEmptyId
	}

	return s.d.DeleteKategori(id)
}
func (s KategoriService) GetKategoriWithPagination(page int, limit int) ([]kategori.Kategori, int, error) {
	return s.d.GetKategoriWithPagination(page, limit)
}
