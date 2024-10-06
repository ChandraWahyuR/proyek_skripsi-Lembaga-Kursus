package handler

import "time"

type KursusRequest struct {
	Nama               string          `json:"nama" form:"nama"`
	Image              []string        `json:"image" form:"image"`
	Deskripsi          string          `json:"deskripsi" form:"deskripsi"`
	Kategori           []string        `json:"kategori" form:"kategori"`
	Jadwal             []JadwalRequest `json:"jadwal"`
	Harga              int             `json:"harga" form:"harga"`
	InstrukturID       string          `json:"instruktur_id" form:"instruktur_id"`
	MateriPembelajaran []string        `json:"materi_pembelajaran" form:"materi"`
}
type JadwalRequest struct {
	Hari       string    `json:"hari" form:"hari"`
	JamMulai   time.Time `json:"jam_mulai" form:"jam_mulai"`
	JamSelesai time.Time `json:"jam_selesai" form:"jam_selesai"`
}
