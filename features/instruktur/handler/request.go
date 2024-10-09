package handler

type PostInstrukturRequest struct {
	Name   string `json:"name" form:"name"`
	Gender string `json:"gender" form:"gender"`
	Email  string `json:"email" form:"email"`
	Alamat string `json:"alamat" form:"alamat"`
	NIP    string `json:"nip" form:"nip"`
	NIK    string `json:"nik" form:"nik"`
	Image  string `json:"image" form:"image"`
	NoHp   string `json:"no_hp" form:"no_hp"`
}

type UpdateInstrukturRequest struct {
	Name   string `json:"name" form:"name"`
	Gender string `json:"gender" form:"gender"`
	Email  string `json:"email" form:"email"`
	Alamat string `json:"alamat" form:"alamat"`
	NIP    string `json:"nip" form:"nip"`
	NIK    string `json:"nik" form:"nik"`
	Image  string `json:"image" form:"image"`
	NoHp   string `json:"no_hp" form:"no_hp"`
}
