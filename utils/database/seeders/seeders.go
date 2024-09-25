package seeders

import "gorm.io/gorm"

type Seeder struct {
	db *gorm.DB
}

func NewSeeder(db *gorm.DB) *Seeder {
	return &Seeder{db: db}
}

func (s *Seeder) Seed() {
	s.SeedAdmins()
}
