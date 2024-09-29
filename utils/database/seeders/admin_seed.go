package seeders

import (
	"fmt"
	"log"
	"skripsi/features/admin/data"
	"skripsi/helper"
)

func (s *Seeder) SeedAdmins() {
	hashedPassword, err := helper.HashPassword("adminKu123")
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	admin := []data.Admin{
		{
			ID:       "02288709-22bc-4d7b-ac91-15d8017c7845",
			Username: "admin",
			Email:    "admin123@example.com",
			Password: hashedPassword,
		},
	}

	for _, admin := range admin {
		result := s.db.FirstOrCreate(&admin, data.Admin{Email: admin.Email})

		// Mengecek apakah record sudah ada atau baru di-create
		if result.Error != nil {
			fmt.Printf("Failed to seed admin %v: %v\n", admin.Username, result.Error)
		} else if result.RowsAffected == 0 {
			fmt.Printf("Admin %v already exists, skipping...\n", admin.Username)
		} else {
			fmt.Printf("Admin %v created successfully\n", admin.Username)
		}
	}
}
