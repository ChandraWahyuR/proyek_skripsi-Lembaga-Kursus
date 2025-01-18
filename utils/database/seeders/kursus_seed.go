package seeders

import (
	"fmt"
	"skripsi/features/kursus/data"
	"time"
)

func (s *Seeder) SeedKursus() {
	kursus := []data.Kursus{
		// {
		// 	ID:           "92ca1f3a-3a58-4d63-91c0-30db2642e0af",
		// 	Nama:         "Corel",
		// 	Deskripsi:    "Kursus belajar corel mulai dari hal basic hingga diakhir akan diberikan tugas akhir dalam membuat proyek",
		// 	Harga:        1500000,
		// 	InstrukturID: "231de840-1008-4b24-ba95-8586e5cac20b",
		// 	Image: []data.ImageKursus{
		// 		{
		// 			ID:       "67cf5efb-3fea-46f0-a1f8-18c23cb6586e",
		// 			Name:     "Corel Gambar 1",
		// 			Url:      "https://storage.googleapis.com/image_skripsi/gambar/kursus/default/Corel.png",
		// 			Position: 1,
		// 			KursusID: "92ca1f3a-3a58-4d63-91c0-30db2642e0af",
		// 		},
		// 	},
		// 	Kategori: []data.KategoriKursus{
		// 		{
		// 			ID:         "c19d699e-fb28-486a-baea-84604c0d9b31",
		// 			KursusID:   "92ca1f3a-3a58-4d63-91c0-30db2642e0af",
		// 			KategoriID: "18cc2761-6b57-4fd7-970b-d123a22cf287",
		// 		},
		// 	},
		// 	Jadwal: []data.JadwalKursus{
		// 		{
		// 			ID:         "7930aae2-4995-40b0-9e17-6c5e2a3de13d",
		// 			KursusID:   "92ca1f3a-3a58-4d63-91c0-30db2642e0af",
		// 			Hari:       "Senin",
		// 			JamMulai:   time.Date(2000, 1, 1, 9, 0, 0, 0, time.UTC),  // Jam mulai 09:00
		// 			JamSelesai: time.Date(2000, 1, 1, 11, 0, 0, 0, time.UTC), // Jam selesai 11:00
		// 		},
		// 		{
		// 			ID:         "f6cff6e2-53cc-4e6c-aa9e-81271555ccaa",
		// 			KursusID:   "92ca1f3a-3a58-4d63-91c0-30db2642e0af",
		// 			Hari:       "Selasa",
		// 			JamMulai:   time.Date(2000, 1, 1, 9, 0, 0, 0, time.UTC),  // Jam mulai 09:00
		// 			JamSelesai: time.Date(2000, 1, 1, 11, 0, 0, 0, time.UTC), // Jam selesai 11:00
		// 		},
		// 	},
		// 	MateriPembelajaran: []data.MateriPembelajaran{
		// 		{
		// 			ID:        "1ec0cf31-86ea-4ab3-9f1a-8c1f64ea1098",
		// 			KursusID:  "92ca1f3a-3a58-4d63-91c0-30db2642e0af",
		// 			Position:  1,
		// 			Deskripsi: "Materi 1",
		// 		},
		// 		{
		// 			ID:        "9824c000-ef7e-41a4-acdd-899c94efb642",
		// 			KursusID:  "92ca1f3a-3a58-4d63-91c0-30db2642e0af",
		// 			Position:  2,
		// 			Deskripsi: "Materi 2",
		// 		},
		// 	},
		// },
		// 1
		{
			ID:           "88abef12-1de5-4e76-bab0-82a9dd5f3d6c",
			Nama:         "Desain Grafis",
			Deskripsi:    "Era Digital seperti sekarang ini memiliki skill sebagai seorang Graphic Design sangat dibutuhkan. Hampir semua teknik dalam beriklan menggunakan bidang keilmuan Design Grafis. Untuk menciptakan sebuah hasil karya design yang bagus bukanlah sebuah hal yang mudah untuk kita lakukan apalagi jika kita tidak memiliki background pendidikan yang sesuai dengan bidang keilmuan ini Namun anda tidak perlu khawatir karena kami akan membimbing anda sampai. Materi yang diberikan dalam pelatihan ini mulai dari Nol dan sangat cocok untuk pemula yang baru belajar mengenai design grafis. Tujuan dari Pelatihan Design Grafis ini adalah peserta pelatihan diharapkan bisa membuat sebuah desain visual untuk promosi yang menarik pasar.",
			Harga:        2500000,
			InstrukturID: "231de840-1008-4b24-ba95-8586e5cac20b",
			Image: []data.ImageKursus{
				{
					ID:       "6ac3dd4b-6862-4c74-af0a-0f675cbdc46c",
					Name:     "Desain Gambar 1",
					Url:      "https://storage.googleapis.com/image_skripsi/gambar/kursus/default/Graphic%20Design",
					Position: 1,
					KursusID: "92ca1f3a-3a58-4d63-91c0-30db2642e0af",
				},
			},
			Kategori: []data.KategoriKursus{
				{
					ID:         "c19d699e-fb28-486a-baea-84604c0d9b31",
					KursusID:   "88abef12-1de5-4e76-bab0-82a9dd5f3d6c",
					KategoriID: "18cc2761-6b57-4fd7-970b-d123a22cf287",
				},
			},
			Jadwal: []data.JadwalKursus{
				{
					ID:         "ac12000d-85e0-4b1b-b6e3-6311253e3057",
					KursusID:   "88abef12-1de5-4e76-bab0-82a9dd5f3d6c",
					Hari:       "Senin",
					JamMulai:   time.Date(2000, 1, 1, 2, 0, 0, 0, time.UTC), // Jam mulai 09:00
					JamSelesai: time.Date(2000, 1, 1, 4, 0, 0, 0, time.UTC), // Jam selesai 11:00
				},
				{
					ID:         "5adc0e41-f37b-421d-a4b7-516f1725afaa",
					KursusID:   "88abef12-1de5-4e76-bab0-82a9dd5f3d6c",
					Hari:       "Selasa",
					JamMulai:   time.Date(2000, 1, 1, 2, 0, 0, 0, time.UTC), // Jam mulai 09:00
					JamSelesai: time.Date(2000, 1, 1, 4, 0, 0, 0, time.UTC), // Jam selesai 11:00
				},
				{
					ID:         "230b528f-cf07-4193-9e95-5166a16f305f",
					KursusID:   "88abef12-1de5-4e76-bab0-82a9dd5f3d6c",
					Hari:       "Rabu",
					JamMulai:   time.Date(2000, 1, 1, 2, 0, 0, 0, time.UTC), // Jam mulai 09:00
					JamSelesai: time.Date(2000, 1, 1, 4, 0, 0, 0, time.UTC), // Jam selesai 11:00
				},
				{
					ID:         "8b193c15-d655-40cd-90c1-64b971af1323",
					KursusID:   "88abef12-1de5-4e76-bab0-82a9dd5f3d6c",
					Hari:       "Kamis",
					JamMulai:   time.Date(2000, 1, 1, 2, 0, 0, 0, time.UTC), // Jam mulai 09:00
					JamSelesai: time.Date(2000, 1, 1, 4, 0, 0, 0, time.UTC), // Jam selesai 11:00
				},
				{
					ID:         "ceb4a079-f3ad-427b-af67-b706423a4c70",
					KursusID:   "88abef12-1de5-4e76-bab0-82a9dd5f3d6c",
					Hari:       "Jum'at",
					JamMulai:   time.Date(2000, 1, 1, 2, 0, 0, 0, time.UTC), // Jam mulai 09:00
					JamSelesai: time.Date(2000, 1, 1, 4, 0, 0, 0, time.UTC), // Jam selesai 11:00
				},
			},
			MateriPembelajaran: []data.MateriPembelajaran{
				{
					ID:        "1c4a5b6f-3d9b-4dfe-8a11-0b9463288ddf",
					KursusID:  "88abef12-1de5-4e76-bab0-82a9dd5f3d6c",
					Position:  1,
					Deskripsi: "Dasar-dasar layout",
				},
				{
					ID:        "a4a3b70b-682f-48cf-9e9e-da284d2950cb",
					KursusID:  "88abef12-1de5-4e76-bab0-82a9dd5f3d6c",
					Position:  2,
					Deskripsi: "Membuat Master Page pada beberapa media publikasi",
				},
				{
					ID:        "1687e73f-88b7-4529-828e-963d946981ba",
					KursusID:  "88abef12-1de5-4e76-bab0-82a9dd5f3d6c",
					Position:  3,
					Deskripsi: "Pengolahan shape dasar dan manipulasi shape",
				},
				{
					ID:        "6f1a2ee0-c336-4cfd-bde7-c5023bbe0a45",
					KursusID:  "88abef12-1de5-4e76-bab0-82a9dd5f3d6c",
					Position:  4,
					Deskripsi: "Melakukan export dan import data gambar",
				},
			},
		},
		// 2
		{
			ID:           "7af3c58b-fa0a-4221-9d1d-14035fe43718",
			Nama:         "Aplikasi Perkantoran",
			Deskripsi:    "Aplikasi perkantoran (Inggris:office suite) adalah sebuah perangkat lunak yang diperuntukkan khusus untuk pekerjaan di kantor. Komponen-komponennya umumnya didistribusikan bersamaan,memiliki antar muka pengguna yang konsisten dan dapat berinteraksi satu sama lain.Setelah pelatihan ini diharapkan peserta mampu menggunakan aplikasi perkantoran, serta mampu mengelolanya dalam lingkup pekerjaan kantor sebagai pengelola dokumen secara digitilisasi.",
			Harga:        1300000,
			InstrukturID: "d825743a-2698-4227-8964-bbbfe91e2d93",
			Image: []data.ImageKursus{
				{
					ID:       "dd31b57b-4cb0-40d1-863d-a6df702e6bb2",
					Name:     "Office Gambar 1",
					Url:      "https://storage.googleapis.com/image_skripsi/gambar/kursus/default/Operasi%20Office",
					Position: 1,
					KursusID: "7af3c58b-fa0a-4221-9d1d-14035fe43718",
				},
			},
			Kategori: []data.KategoriKursus{
				{
					ID:         "6fdfdcc6-42f6-416f-86ab-efb8f933c7ac",
					KursusID:   "7af3c58b-fa0a-4221-9d1d-14035fe43718",
					KategoriID: "f0adef3c-179a-4b35-a6e9-8badc4406a1b",
				},
			},
			Jadwal: []data.JadwalKursus{
				{
					ID:         "0f995bc4-1c30-4c86-9694-1b214c67b4e0",
					KursusID:   "7af3c58b-fa0a-4221-9d1d-14035fe43718",
					Hari:       "Senin",
					JamMulai:   time.Date(2000, 1, 1, 2, 0, 0, 0, time.UTC), // Jam mulai 09:00
					JamSelesai: time.Date(2000, 1, 1, 4, 0, 0, 0, time.UTC), // Jam selesai 11:00
				},
				{
					ID:         "fc1b36e8-670c-408c-9b23-f978e3353579",
					KursusID:   "7af3c58b-fa0a-4221-9d1d-14035fe43718",
					Hari:       "Selasa",
					JamMulai:   time.Date(2000, 1, 1, 2, 0, 0, 0, time.UTC), // Jam mulai 09:00
					JamSelesai: time.Date(2000, 1, 1, 4, 0, 0, 0, time.UTC), // Jam selesai 11:00
				},
				{
					ID:         "55f03250-f500-4e4c-bbda-48f11775f8ed",
					KursusID:   "7af3c58b-fa0a-4221-9d1d-14035fe43718",
					Hari:       "Rabu",
					JamMulai:   time.Date(2000, 1, 1, 2, 0, 0, 0, time.UTC), // Jam mulai 09:00
					JamSelesai: time.Date(2000, 1, 1, 4, 0, 0, 0, time.UTC), // Jam selesai 11:00
				},
				{
					ID:         "2cb9a577-4fbd-4961-b387-1f4183765f92",
					KursusID:   "7af3c58b-fa0a-4221-9d1d-14035fe43718",
					Hari:       "Kamis",
					JamMulai:   time.Date(2000, 1, 1, 2, 0, 0, 0, time.UTC), // Jam mulai 09:00
					JamSelesai: time.Date(2000, 1, 1, 4, 0, 0, 0, time.UTC), // Jam selesai 11:00
				},
				{
					ID:         "e9f21a1d-587a-4b30-ab09-116352875e2e",
					KursusID:   "7af3c58b-fa0a-4221-9d1d-14035fe43718",
					Hari:       "Jum'at",
					JamMulai:   time.Date(2000, 1, 1, 2, 0, 0, 0, time.UTC), // Jam mulai 09:00
					JamSelesai: time.Date(2000, 1, 1, 4, 0, 0, 0, time.UTC), // Jam selesai 11:00
				},
			},
			MateriPembelajaran: []data.MateriPembelajaran{
				{
					ID:        "ec7c9d3b-491c-4c07-8c5f-4498e15e1c20",
					KursusID:  "7af3c58b-fa0a-4221-9d1d-14035fe43718",
					Position:  1,
					Deskripsi: "Mengoperasikan sistem operasi",
				},
				{
					ID:        "a617de7a-588e-4a5b-9237-b217cc3886cd",
					KursusID:  "7af3c58b-fa0a-4221-9d1d-14035fe43718",
					Position:  2,
					Deskripsi: "Mengoperasikan piranti lunak Pengolah kata",
				},
				{
					ID:        "fa2c7114-8eb9-4345-8ed2-bcd9c7c4e49f",
					KursusID:  "7af3c58b-fa0a-4221-9d1d-14035fe43718",
					Position:  3,
					Deskripsi: "Mengoperasikan piranti lunak Pengolah Angka",
				},
				{
					ID:        "7cce77c9-6e31-4596-bd21-67df73b94489",
					KursusID:  "7af3c58b-fa0a-4221-9d1d-14035fe43718",
					Position:  4,
					Deskripsi: "Mengoperasikan piranti lunak Presentasi",
				},
				{
					ID:        "1175424b-09d7-46ed-bdad-0adb8c1b14fd",
					KursusID:  "7af3c58b-fa0a-4221-9d1d-14035fe43718",
					Position:  5,
					Deskripsi: "Mengoperasikan Internet",
				},
			},
		},
		// 3
		{
			ID:           "d2406e55-7fb5-4b92-8623-08cd334537bb",
			Nama:         "Komputer For Kids",
			Deskripsi:    "Setelah pelatihan ini diharapkan peserta mampu menggunakan aplikasi perkantoran, serta mampu mengelolanya dalam lingkup pekerjaan sebagai pengelola dokumen secara digitilisasi.",
			Harga:        600000,
			InstrukturID: "c8282a4e-1c2e-4a55-86ca-643089272922",
			Image: []data.ImageKursus{
				{
					ID:       "883d80b4-0631-435a-ab24-6c1666f010b5",
					Name:     "Komputer For Kids Gambar 1",
					Url:      "https://storage.googleapis.com/image_skripsi/gambar/kursus/default/Komputer%20For%20Kids",
					Position: 1,
					KursusID: "d2406e55-7fb5-4b92-8623-08cd334537bb",
				},
			},
			Kategori: []data.KategoriKursus{
				{
					ID:         "f176cfaf-fcd2-4ce8-bb2d-d6bd808988b8",
					KursusID:   "d2406e55-7fb5-4b92-8623-08cd334537bb",
					KategoriID: "f0adef3c-179a-4b35-a6e9-8badc4406a1b",
				},
			},
			Jadwal: []data.JadwalKursus{
				{
					ID:         "eeca6ac2-7e01-4c98-805f-a933096475da",
					KursusID:   "d2406e55-7fb5-4b92-8623-08cd334537bb",
					Hari:       "Senin",
					JamMulai:   time.Date(2000, 1, 1, 2, 0, 0, 0, time.UTC), // Jam mulai 09:00
					JamSelesai: time.Date(2000, 1, 1, 4, 0, 0, 0, time.UTC), // Jam selesai 11:00
				},
				{
					ID:         "2d582593-9c12-42b2-a33d-438fb569d173",
					KursusID:   "d2406e55-7fb5-4b92-8623-08cd334537bb",
					Hari:       "Selasa",
					JamMulai:   time.Date(2000, 1, 1, 2, 0, 0, 0, time.UTC), // Jam mulai 09:00
					JamSelesai: time.Date(2000, 1, 1, 4, 0, 0, 0, time.UTC), // Jam selesai 11:00
				},
				{
					ID:         "d79349fc-186e-45ae-8bde-bf53006c8900",
					KursusID:   "d2406e55-7fb5-4b92-8623-08cd334537bb",
					Hari:       "Rabu",
					JamMulai:   time.Date(2000, 1, 1, 2, 0, 0, 0, time.UTC), // Jam mulai 09:00
					JamSelesai: time.Date(2000, 1, 1, 4, 0, 0, 0, time.UTC), // Jam selesai 11:00
				},
				{
					ID:         "dc01b3e7-7d07-4b8b-9e17-4d6fc796b5df",
					KursusID:   "d2406e55-7fb5-4b92-8623-08cd334537bb",
					Hari:       "Kamis",
					JamMulai:   time.Date(2000, 1, 1, 2, 0, 0, 0, time.UTC), // Jam mulai 09:00
					JamSelesai: time.Date(2000, 1, 1, 4, 0, 0, 0, time.UTC), // Jam selesai 11:00
				},
				{
					ID:         "e0df5150-9878-4687-b75e-1b8551397613",
					KursusID:   "d2406e55-7fb5-4b92-8623-08cd334537bb",
					Hari:       "Jum'at",
					JamMulai:   time.Date(2000, 1, 1, 2, 0, 0, 0, time.UTC), // Jam mulai 09:00
					JamSelesai: time.Date(2000, 1, 1, 4, 0, 0, 0, time.UTC), // Jam selesai 11:00
				},
			},
			MateriPembelajaran: []data.MateriPembelajaran{
				{
					ID:        "d3dc36ae-5cd1-40c7-867c-a83bac46e75f",
					KursusID:  "d2406e55-7fb5-4b92-8623-08cd334537bb",
					Position:  1,
					Deskripsi: "Mengoperasikan sistem operasi",
				},
				{
					ID:        "ec671457-5071-4cee-b3b1-bbe692f6a633",
					KursusID:  "d2406e55-7fb5-4b92-8623-08cd334537bb",
					Position:  2,
					Deskripsi: "Mengoperasikan piranti lunak Pengolah kata",
				},
				{
					ID:        "481b2021-72c9-40a1-b203-a03d948811b0",
					KursusID:  "d2406e55-7fb5-4b92-8623-08cd334537bb",
					Position:  3,
					Deskripsi: "Mengoperasikan piranti lunak Pengolah Angka",
				},
				{
					ID:        "0b1d8322-8742-44bb-adef-0e931ca48c71",
					KursusID:  "d2406e55-7fb5-4b92-8623-08cd334537bb",
					Position:  4,
					Deskripsi: "Mengoperasikan piranti lunak Presentasi",
				},
				{
					ID:        "581c4a58-95d9-4fdb-87c6-a05ed64d4ad4",
					KursusID:  "d2406e55-7fb5-4b92-8623-08cd334537bb",
					Position:  5,
					Deskripsi: "Mengoperasikan Internet",
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
