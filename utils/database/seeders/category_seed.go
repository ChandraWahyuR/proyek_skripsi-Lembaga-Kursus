package seeders

import (
	"fmt"
	"skripsi/features/kategori/data"
)

func (s *Seeder) SeedCategory() {
	categories := []data.Kategori{
		{
			ID:        "18cc2761-6b57-4fd7-970b-d123a22cf287",
			Nama:      "Category Desain",
			Deskripsi: "Category yang membahas tentang desain grafis",
			ImageUrl:  "https://storage.googleapis.com/image_skripsi/gambar/kategori/default/Question%20Mark.png",
		},
	}
	for _, category := range categories {
		result := s.db.FirstOrCreate(&category, data.Kategori{ID: category.ID})
		if result.Error != nil {
			fmt.Printf("Failed to seed voucher %v: %v\n", category.ID, result.Error)
		}
	}
}
