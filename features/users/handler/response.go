package handler

import "time"

type UserLoginResponse struct {
	Token string `json:"token"`
}

type MetadataResponse struct {
	TotalPage int `json:"total_page"`
	Page      int `json:"current_page"`
}

type GetAllUserResponse struct {
	ID       string `json:"id"`
	NIS      string `json:"nis"`
	Username string `json:"username"`
	Nama     string `json:"nama"`
	Email    string `json:"email"`
	NomorHP  string `json:"nomor_hp"`
	//
	ProfileUrl    string    `json:"profile_url"`
	IsActive      bool      `json:"is_active"`
	Agama         string    `json:"agama"`
	Gender        string    `json:"gender"`
	TempatLahir   string    `json:"tempat_lahir"`
	TanggalLahir  time.Time `json:"tanggal_lahir"`
	OrangTua      string    `json:"orangtua"`
	Profesi       string    `json:"profesi"`
	Ijazah        string    `json:"ijazah"`
	KTP           string    `json:"ktp"`
	KartuKeluarga string    `json:"kartu_keluarga"`
}

type GetUserIDResponse struct {
	ID            string    `json:"id"`
	NIS           string    `json:"nis"`
	Username      string    `json:"username"`
	Nama          string    `json:"nama"`
	Email         string    `json:"email"`
	NomorHP       string    `json:"nomor_hp"`
	Agama         string    `json:"agama"`
	Gender        string    `json:"gender"`
	TempatLahir   string    `json:"tempat_lahir"`
	TanggalLahir  time.Time `json:"tanggal_lahir"`
	OrangTua      string    `json:"orang_tua"`
	Profesi       string    `json:"profesi"`
	Ijazah        string    `json:"ijazah"`
	KTP           string    `json:"ktp"`
	KartuKeluarga string    `json:"kartu_keluarga"`
	ProfileUrl    string    `json:"profile_url"`
}
