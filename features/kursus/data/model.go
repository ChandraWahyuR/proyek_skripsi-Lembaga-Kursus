package data

import (
	dataInstruktur "skripsi/features/instruktur/data"
	dataKategori "skripsi/features/kategori/data"
	"time"

	"gorm.io/gorm"
)

type Kursus struct {
	*gorm.Model
	ID                 string               `gorm:"type:varchar(50);not null;column:id"`
	Nama               string               `gorm:"type:varchar(255);not null;column:nama"`
	Image              []ImageKursus        `gorm:"foreignKey:ID;references:ID"`
	Deskripsi          string               `gorm:"type:text;not null;column:deskripsi"`
	Kategori           []KategoriKursus     `gorm:"foreignKey:ID;references:ID"`
	Jadwal             time.Time            `gorm:"not null;column:jadwal"`
	Harga              int                  `gorm:"type:int;not null;column:harga"`
	Instruktur         InstrukturKursus     `gorm:"foreignKey:ID;references:ID"`
	MateriPembelajaran []MateriPembelajaran `gorm:"foreignKey:KursusID"`
}

type ImageKursus struct {
	*gorm.Model
	ID       string `gorm:"type:varchar(50);not null;column:id"`
	Name     string `gorm:"type:varchar(255);not null;column:name"`
	Url      string `gorm:"type:text;not null;column:url"`
	Position int    `gorm:"type:int;not null;column:position"`
	KursusID string `gorm:"type:varchar(50);not null;column:kursus_id"`
}

type KategoriKursus struct {
	*gorm.Model
	ID         string                `gorm:"type:varchar(50);not null;column:id"`
	KursusID   string                `gorm:"type:varchar(50);not null;column:kursus_id"`
	KategoriID string                `gorm:"type:varchar(50);not null;column:kategori_id"`
	Kategori   dataKategori.Kategori `gorm:"foreignKey:ID;references:ID"`
}

type InstrukturKursus struct {
	*gorm.Model
	ID           string                    `gorm:"type:varchar(50);not null;column:id"`
	KursusID     string                    `gorm:"type:varchar(50);not null;column:kursus_id"`
	InstrukturID string                    `gorm:"type:varchar(50);not null;column:instruktur_id"`
	Instruktur   dataInstruktur.Instruktur `gorm:"foreignKey:ID;references:ID"`
}

type MateriPembelajaran struct {
	*gorm.Model
	ID       string `gorm:"type:varchar(50);not null;column:id"`
	KursusID string `gorm:"type:varchar(50);not null;column:kursus_id"`
	MateriID string `gorm:"type:varchar(50);not null;column:materi_id"`
	Materi   Materi `gorm:"foreignKey:ID;references:ID"`
}

type Materi struct {
	*gorm.Model
	ID        string `gorm:"type:varchar(50);not null;column:id"`
	Name      string `gorm:"type:varchar(255);not null;column:name"`
	Deskripsi string `gorm:"type:text;not null;column:deskripsi"`
}

func (Kursus) TableName() string {
	return "kursuss"
}

func (ImageKursus) TableName() string {
	return "image_kursuss"
}

func (KategoriKursus) TableName() string {
	return "kategori_kursuss"
}

func (InstrukturKursus) TableName() string {
	return "instruktur_kursuss"
}

func (MateriPembelajaran) TableName() string {
	return "materi_pembelajaranss"
}

func (Materi) TableName() string {
	return "materis"
}
