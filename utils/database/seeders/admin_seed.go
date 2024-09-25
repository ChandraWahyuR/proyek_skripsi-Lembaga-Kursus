package seeders

import (
	"fmt"
	"log"
	"skripsi/features/admin/data"
	"skripsi/helper"

	"github.com/google/uuid"
)

func (s *Seeder) SeedAdmins() {
	hashedPassword, err := helper.HashPassword("adminKu123")
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	admin := []data.Admin{
		{
			ID:       uuid.New().String(),
			Username: "admin",
			Email:    "admin123@example.com",
			Password: hashedPassword,
		},
	}

	for _, admin := range admin {
		if err := s.db.FirstOrCreate(&admin, data.Admin{Email: admin.Email}).Error; err != nil {
			fmt.Printf("Failed to seed admin %v: %v\n", admin.Username, err)
		} else {
			fmt.Printf("Admin %v created/exists\n", admin.Username)
		}
	}
}
