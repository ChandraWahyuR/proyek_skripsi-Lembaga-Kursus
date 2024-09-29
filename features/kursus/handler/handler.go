package handler

import (
	"skripsi/features/kursus"
	"skripsi/helper"
)

type KursusHandler struct {
	s kursus.KursusServiceInterface
	j helper.JWTInterface
}

// func New(u kursus.KursusServiceInterface, j helper.JWTInterface) kursus.KursusHandlerInterface {
// 	return &KursusHandler{
// 		s: u,
// 		j: j,
// 	}
// }

// func(h *KursusHandler)
