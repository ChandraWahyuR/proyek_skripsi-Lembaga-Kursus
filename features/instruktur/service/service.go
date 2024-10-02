package service

import (
	"skripsi/constant"
	"skripsi/features/instruktur"
	"skripsi/helper"
)

type InstrukturService struct {
	d instruktur.InstrukturDataInterface
	j helper.JWTInterface
}

func New(u instruktur.InstrukturDataInterface, j helper.JWTInterface) instruktur.InstrukturServiceInterface {
	return &InstrukturService{
		d: u,
		j: j,
	}
}

func (s *InstrukturService) GetInstrukturWithPagination(page int, limit int) ([]instruktur.Instruktur, int, error) {
	return s.d.GetInstrukturWithPagination(page, limit)
}

func (s *InstrukturService) GetAllInstruktur() ([]instruktur.Instruktur, error) {
	return s.d.GetAllInstruktur()
}

func (s *InstrukturService) GetAllInstrukturByID(id string) (instruktur.Instruktur, error) {
	return s.d.GetAllInstrukturByID(id)
}

func (s *InstrukturService) PostInstruktur(data instruktur.Instruktur) error {
	switch {
	case data.Name == "":
		return constant.ErrEmptyNameInstuktor
	case data.Email == "":
		return constant.ErrEmptyEmailInstuktor
	case data.Alamat == "":
		return constant.ErrEmptyAlamatInstuktor
	case data.NoHp == "":
		return constant.ErrEmptyNumbertelponInstuktor
	case data.Gender == "":
		return constant.ErrEmptyDescriptionInstuktor
	}

	return s.d.PostInstruktur(data)
}

func (s *InstrukturService) UpdateInstruktur(data instruktur.UpdateInstruktur) error {
	if data.ID == "" {
		return constant.ErrEmptyId
	}
	// 	&& memastikan bahwa semua field kosong sebelum mengembalikan error.
	// || memastikan bahwa hanya salah satu field yang kosong untuk mengembalikan error.
	// && and kan wajib terpenuhi semua jika return err. jika tidak, aman
	if data.Name == "" && data.Email == "" && data.Alamat == "" && data.Gender == "" && data.NoHp == "" {
		return constant.ErrUpdate
	}

	nomorHp, err := helper.TelephoneValidator(data.NoHp)
	if err != nil {
		return err
	}
	data.NoHp = nomorHp

	return s.d.UpdateInstruktur(data)
}

func (s *InstrukturService) DeleteInstruktur(id string) error {
	if id == "" {
		return constant.ErrEmptyId
	}

	return s.d.DeleteInstruktur(id)
}

func (s *InstrukturService) GetInstruktorByName(name string, page int, limit int) ([]instruktur.Instruktur, int, error) {
	data, total, err := s.d.GetInstruktorByName(name, page, limit)

	if err != nil {
		return nil, 0, err
	}

	return data, total, nil
}
