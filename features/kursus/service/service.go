package service

import (
	"skripsi/constant"
	"skripsi/features/kursus"
	"skripsi/helper"
)

type KursusService struct {
	d kursus.KursusDataInterface
	j helper.JWTInterface
}

func New(u kursus.KursusDataInterface, j helper.JWTInterface) kursus.KursusServiceInterface {
	return &KursusService{
		d: u,
		j: j,
	}
}

func (s *KursusService) GetAllKursus() ([]kursus.Kursus, error) {
	return s.d.GetAllKursus()
}

func (s *KursusService) GetAllKursusById(id string) (kursus.Kursus, error) {
	if id == "" {
		return kursus.Kursus{}, constant.ErrGetID
	}
	return s.d.GetAllKursusById(id)
}

func (s *KursusService) AddKursus(data kursus.Kursus) error {
	switch {
	case data.Nama == "":
		return constant.ErrEmptyNameInstuktor
	case data.Deskripsi == "":
		return constant.ErrDekripsiKursus
	case data.Harga == 0:
		return constant.ErrHargaKursus
	case data.InstrukturID == "":
		return constant.ErrInstrukturID
	case len(data.Jadwal) == 0:
		return constant.ErrJadwal
	case len(data.Image) == 0:
		return constant.ErrGambarKursus
	case len(data.Kategori) == 0:
		return constant.ErrKategoriKursus
	case len(data.MateriPembelajaran) == 0:
		return constant.ErrMateriPembelajaran
	}
	return s.d.AddKursus(data)
}

func (s *KursusService) GetKursusPagination(page int, limit int) ([]kursus.Kursus, int, error) {
	return s.d.GetKursusPagination(page, limit)
}

func (s *KursusService) UpdateKursus(data kursus.Kursus) error {
	if data.ID == "" {
		return constant.ErrEmptyId
	}
	if data.Nama == "" && data.Deskripsi == "" && data.Harga == 0 && data.InstrukturID == "" && len(data.Jadwal) == 0 && len(data.Image) == 0 && len(data.Kategori) == 0 && len(data.MateriPembelajaran) == 0 {
		return constant.ErrUpdate
	}

	return s.d.UpdateKursus(data)
}

func (s *KursusService) DeleteKursus(id string) error {
	if id == "" {
		return constant.ErrEmptyId
	}

	return s.d.DeleteKursus(id)
}

func (s *KursusService) DeleteImageKursus(id string) error {
	return s.d.DeleteImageKursus(id)
}

func (s *KursusService) DeleteMateriKursus(id string) error {
	return s.d.DeleteMateriKursus(id)
}

func (s *KursusService) DeleteKategoriKursus(id string) error {
	return s.d.DeleteKategoriKursus(id)
}

func (s *KursusService) GetAllKursusByName(name string, page int, limit int) ([]kursus.Kursus, int, error) {
	data, total, err := s.d.GetAllKursusByName(name, page, limit)

	if err != nil {
		return nil, 0, err
	}

	// Pastikan data tidak nil
	if data == nil {
		data = []kursus.Kursus{}
	}

	return data, total, nil
}
