package service

import (
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
	return s.d.GetAllKursusById(id)
}

func (s *KursusService) AddKursus(data kursus.Kursus) error {
	return s.d.AddKursus(data)
}

func (s *KursusService) GetKursusPagination(page int, limit int) ([]kursus.Kursus, int, error) {
	return s.d.GetKursusPagination(page, limit)
}
