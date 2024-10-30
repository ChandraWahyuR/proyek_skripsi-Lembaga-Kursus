package seeders

import (
	"fmt"
	"skripsi/features/instruktur/data"
)

func (s *Seeder) SeedInstruktur() {
	instrukturs := []data.Instruktur{
		{
			ID:                   "231de840-1008-4b24-ba95-8586e5cac20b",
			Name:                 "Adit Setiawan",
			NomorIndukPendidikan: "198502212009031005",
			NIK:                  "3451234114131",
			Gender:               "laki-laki",
			Email:                "aditsetiawan0121@gmail.com",
			Alamat:               "Cilacap, Jeruklegi",
			NoHp:                 "081234453231",
			UrlImage:             "https://storage.googleapis.com/image_skripsi/gambar/users/default/Profile_User.png",
		},
		{
			ID:                   "d825743a-2698-4227-8964-bbbfe91e2d93",
			Name:                 "Budi Setiawan",
			NomorIndukPendidikan: "197505012002071021",
			NIK:                  "3451234114131",
			Gender:               "laki-laki",
			Email:                "budisetiawan0121@gmail.com",
			Alamat:               "Banyumas, Purwokerto",
			NoHp:                 "081232153258",
			UrlImage:             "https://storage.googleapis.com/image_skripsi/gambar/users/default/Profile_User.png",
		},
		{
			ID:                   "c8282a4e-1c2e-4a55-86ca-643089272922",
			Name:                 "Dini Setianingsih",
			NomorIndukPendidikan: "199501152010022003",
			NIK:                  "3451234114131",
			Gender:               "perempuan",
			Email:                "dinisetianingsih0121@gmail.com",
			Alamat:               "Cilacap, Bantasari",
			NoHp:                 "081232153258",
			UrlImage:             "https://storage.googleapis.com/image_skripsi/gambar/users/default/Profile_User.png",
		},
	}
	for _, instruktur := range instrukturs {
		result := s.db.FirstOrCreate(&instruktur, data.Instruktur{ID: instruktur.ID})
		if result.Error != nil {
			fmt.Printf("Failed to seed voucher %v: %v\n", instruktur.ID, result.Error)
		}
	}
}
