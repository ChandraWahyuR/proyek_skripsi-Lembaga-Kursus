package handler

type MetadataResponse struct {
	TotalPage int `json:"total_page"`
	Page      int `json:"current_page"`
}

type KategoriResponse struct {
	ID        string `json:"id"`
	Nama      string `json:"nama"`
	Deskripsi string `json:"deskripsi"`
	ImageUrl  string `json:"image_url"`
}
