package seeders

import (
	"fmt"
	"skripsi/features/kursus/data"
	"time"
)

func (s *Seeder) SeedKursus() {
	kursus := []data.Kursus{
		{
			ID:           "92ca1f3a-3a58-4d63-91c0-30db2642e0af",
			Nama:         "Corel",
			Deskripsi:    "Kursus belajar corel mulai dari hal basic hingga diakhir akan diberikan tugas akhir dalam membuat proyek",
			Harga:        1500000,
			InstrukturID: "231de840-1008-4b24-ba95-8586e5cac20b",
			Image: []data.ImageKursus{
				{
					ID:       "67cf5efb-3fea-46f0-a1f8-18c23cb6586e",
					Name:     "Corel Gambar 1",
					Url:      "https://storage.googleapis.com/image_skripsi/gambar/kursus/default/Corel.png",
					Position: 1,
					KursusID: "92ca1f3a-3a58-4d63-91c0-30db2642e0af",
				},
			},
			Kategori: []data.KategoriKursus{
				{
					ID:         "c19d699e-fb28-486a-baea-84604c0d9b31",
					KursusID:   "92ca1f3a-3a58-4d63-91c0-30db2642e0af",
					KategoriID: "18cc2761-6b57-4fd7-970b-d123a22cf287",
				},
			},
			Jadwal: []data.JadwalKursus{
				{
					ID:         "7930aae2-4995-40b0-9e17-6c5e2a3de13d",
					KursusID:   "92ca1f3a-3a58-4d63-91c0-30db2642e0af",
					Hari:       "Senin",
					JamMulai:   time.Date(2000, 1, 1, 9, 0, 0, 0, time.UTC),  // Jam mulai 09:00
					JamSelesai: time.Date(2000, 1, 1, 11, 0, 0, 0, time.UTC), // Jam selesai 11:00
				},
				{
					ID:         "f6cff6e2-53cc-4e6c-aa9e-81271555ccaa",
					KursusID:   "92ca1f3a-3a58-4d63-91c0-30db2642e0af",
					Hari:       "Selasa",
					JamMulai:   time.Date(2000, 1, 1, 9, 0, 0, 0, time.UTC),  // Jam mulai 09:00
					JamSelesai: time.Date(2000, 1, 1, 11, 0, 0, 0, time.UTC), // Jam selesai 11:00
				},
			},
			MateriPembelajaran: []data.MateriPembelajaran{
				{
					ID:        "1ec0cf31-86ea-4ab3-9f1a-8c1f64ea1098",
					KursusID:  "92ca1f3a-3a58-4d63-91c0-30db2642e0af",
					Position:  1,
					Deskripsi: "Materi 1",
				},
				{
					ID:        "9824c000-ef7e-41a4-acdd-899c94efb642",
					KursusID:  "92ca1f3a-3a58-4d63-91c0-30db2642e0af",
					Position:  2,
					Deskripsi: "Materi 2",
				},
			},
		},
	}
	for _, dataKursus := range kursus {
		result := s.db.FirstOrCreate(&dataKursus, data.Kursus{ID: dataKursus.ID})
		if result.Error != nil {
			fmt.Printf("Failed to seed voucher %v: %v\n", dataKursus.ID, result.Error)
		}
	}
}
