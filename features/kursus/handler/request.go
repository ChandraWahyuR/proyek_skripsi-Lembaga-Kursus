package handler

import (
	"time"
)

type KursusRequest struct {
	Nama               string    `json:"nama" form:"nama"`
	Image              []string  `json:"image" form:"image"`
	Deskripsi          string    `json:"deskripsi" form:"deskripsi"`
	Kategori           []string  `json:"kategori" form:"kategori"`
	Jadwal             time.Time `json:"jadwal" form:"jadwal"`
	Harga              int       `json:"harga" form:"harga"`
	InstruktorID       string    `json:"instruktur_id" form:"instruktor_id"`
	MateriPembelajaran []string  `json:"materi_pembelajaran" form:"materi"`
}
