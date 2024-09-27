package handler

type MetadataResponse struct {
	TotalPage int `json:"total_page"`
	Page      int `json:"current_page"`
}

type DataInsrukturResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Gender string `json:"gender"`
	Email  string `json:"email"`
	Alamat string `json:"alamat"`
	NoHp   string `json:"no_hp"`
}
