package handler

type Request struct {
	InstrukturId string `json:"instruktur_id"`
	KursusId     string `json:"kursus_id"`
	Tanggal      string `json:"tanggal"`
	JamMulai     string `json:"jam_mulai"`
	JamAkhir     string `json:"jam_akhir"`
	Status       bool   `json:"status"`
}
