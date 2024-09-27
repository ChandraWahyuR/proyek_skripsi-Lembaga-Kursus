package handler

type PostInstrukturRequest struct {
	Name   string `json:"name"`
	Gender string `json:"gender"`
	Email  string `json:"email"`
	Alamat string `json:"alamat"`
	NoHp   string `json:"no_hp"`
}

type UpdateInstrukturRequest struct {
	Name   string `json:"name"`
	Gender string `json:"gender"`
	Email  string `json:"email"`
	Alamat string `json:"alamat"`
	NoHp   string `json:"no_hp"`
}
