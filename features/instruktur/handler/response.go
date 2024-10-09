package handler

type MetadataResponse struct {
	TotalPage int `json:"total_page"`
	Page      int `json:"current_page"`
}

type DataInsrukturResponseID struct {
	ID     string `json:"id"`
	NIK    string `json:"nik"`
	NIP    string `json:"nip"`
	Image  string `json:"image"`
	Name   string `json:"name"`
	Gender string `json:"gender"`
	Email  string `json:"email"`
	Alamat string `json:"alamat"`
	NoHp   string `json:"no_hp"`
}

type DataInsrukturResponseAll struct {
	ID    string `json:"id"`
	NIK   string `json:"nik"`
	NIP   string `json:"nip"`
	Name  string `json:"name"`
	Email string `json:"email"`
	NoHp  string `json:"no_hp"`
}
