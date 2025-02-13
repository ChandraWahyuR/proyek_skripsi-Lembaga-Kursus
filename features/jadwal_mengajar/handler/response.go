package handler

type GetJadwalMengajar struct {
	ID         string     `json:"id"`
	Instruktur Instruktur `json:"instruktur"`
	User       User       `json:"user"`
	Kursus     Kursus     `json:"kursus"`
	Tanggal    string     `json:"tanggal"`
	Status     bool       `json:"status"`
}

type GetDetailJadwalMengajar struct {
	ID         string     `json:"id"`
	Instruktur Instruktur `json:"instruktur"`
	User       User       `json:"user"`
	Kursus     Kursus     `json:"kursus"`
	Tanggal    string     `json:"tanggal"`
	JamMulai   string     `json:"jam_mulai"`
	JamAkhir   string     `json:"jam_akhir"`
	Status     bool       `json:"status"`
}

type Instruktur struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}
type Kursus struct {
	ID   string `json:"id"`
	Nama string `json:"nama"`
}
