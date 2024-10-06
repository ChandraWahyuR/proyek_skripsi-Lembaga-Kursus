package handler

type MetadataResponse struct {
	TotalPage int `json:"total_page"`
	Page      int `json:"current_page"`
}

// Get All
type ResponseGetKursus struct {
	ID                 string               `json:"id"`
	Nama               string               `json:"nama"`
	Image              []ImageKursus        `json:"image"`
	Deskripsi          string               `json:"deskripsi"`
	Kategori           []KategoriKursus     `json:"kategori"`
	Jadwal             []JadwalKursus       `json:"jadwal"`
	Harga              int                  `json:"harga"`
	Instruktur         Instruktur           `json:"instruktur"`
	MateriPembelajaran []MateriPembelajaran `json:"materi_pembelajaran"`
}

type ResponseGetAllKursus struct {
	ID        string         `json:"id"`
	Nama      string         `json:"nama"`
	Deskripsi string         `json:"deskripsi"`
	Image     []ImageKursus  `json:"image"`
	Harga     int            `json:"harga"`
	Jadwal    []JadwalKursus `json:"jadwal"`
}
type ImageKursus struct {
	Name     string `json:"nama"`
	Url      string `json:"url"`
	Position int    `json:"position"`
}

type KategoriKursus struct {
	Nama      string `json:"nama"`
	Deskripsi string `json:"deskripsi"`
	ImageUrl  string `json:"image_url"`
}

type Instruktur struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type MateriPembelajaran struct {
	Deskripsi string `json:"deskripsi"`
}

type JadwalKursus struct {
	Hari       string
	JamMulai   string
	JamSelesai string
}

//
