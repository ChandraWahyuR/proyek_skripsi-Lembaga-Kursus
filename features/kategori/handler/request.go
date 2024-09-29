package handler

type RequestKategori struct {
	Nama      string `json:"nama" form:"nama"`
	Deskripsi string `json:"deskripsi" form:"deskripsi"`
	Image     string `json:"image" form:"file"`
}
