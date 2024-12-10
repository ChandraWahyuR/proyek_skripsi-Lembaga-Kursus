package service

import (
	"skripsi/features/gmaps"
)

type GmapsService struct {
	d gmaps.GmapsDataInterface
}

func New(u gmaps.GmapsDataInterface) gmaps.GmapsServiceInterface {
	return &GmapsService{
		d: u,
	}
}

func (s *GmapsService) GetDirections(req gmaps.DirectionsRequest) (gmaps.DirectionsResponse, error) {
	return s.d.GetDirections(req)
}
