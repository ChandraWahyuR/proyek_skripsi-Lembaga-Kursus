package data

import (
	instruktur "skripsi/features/instruktur/data"
	dataKategori "skripsi/features/kategori/data"
	"time"

	"gorm.io/gorm"
)

type Kursus struct {
	*gorm.Model
	ID                 string                `gorm:"type:varchar(50);primaryKey;not null;column:id"`
	Nama               string                `gorm:"type:varchar(255);not null;column:nama"`
	Deskripsi          string                `gorm:"type:text;not null;column:deskripsi"`
	Harga              int                   `gorm:"type:int;not null;column:harga"`
	InstrukturID       string                `gorm:"type:varchar(50);not null;column:instruktur_id"`
	Instruktur         instruktur.Instruktur `gorm:"foreignKey:InstrukturID;references:ID"`
	Jadwal             []JadwalKursus        `gorm:"foreignKey:KursusID"`
	Image              []ImageKursus         `gorm:"foreignKey:KursusID"`
	Kategori           []KategoriKursus      `gorm:"foreignKey:KursusID"`
	MateriPembelajaran []MateriPembelajaran  `gorm:"foreignKey:KursusID"`
}

type ImageKursus struct {
	*gorm.Model
	ID       string `gorm:"type:varchar(50);not null;column:id"`
	Name     string `gorm:"type:varchar(255);not null;column:name"`
	Url      string `gorm:"type:varchar(255);not null;column:url"`
	Position int    `gorm:"type:int;not null;column:position"`
	KursusID string `gorm:"type:varchar(50);not null;column:kursus_id"`
	Kursus   Kursus `gorm:"foreignKey:KursusID;references:ID"`
}

type KategoriKursus struct {
	*gorm.Model
	ID         string                `gorm:"type:varchar(50);not null;column:id"`
	KursusID   string                `gorm:"type:varchar(50);not null;column:kursus_id"`
	Kursus     Kursus                `gorm:"foreignKey:KursusID;references:ID"`
	KategoriID string                `gorm:"type:varchar(50);not null;column:kategori_id"`
	Kategori   dataKategori.Kategori `gorm:"foreignKey:KategoriID;references:ID"`
}

type MateriPembelajaran struct {
	*gorm.Model
	ID        string `gorm:"type:varchar(50);not null;column:id"`
	Position  int    `gorm:"type:int;not null;column:position"`
	KursusID  string `gorm:"type:varchar(50);not null;column:kursus_id"`
	Kursus    Kursus `gorm:"foreignKey:KursusID;references:ID"`
	Deskripsi string `gorm:"type:text;not null;column:deskripsi"`
}

type JadwalKursus struct {
	*gorm.Model
	ID         string    `gorm:"type:varchar(50);primaryKey;not null;column:id"`
	KursusID   string    `gorm:"type:varchar(50);not null;column:kursus_id"`
	Kursus     Kursus    `gorm:"foreignKey:KursusID;references:ID"`
	Hari       string    `gorm:"type:varchar(20);not null;column:hari"`
	JamMulai   time.Time `gorm:"not null;column:jam_mulai"`
	JamSelesai time.Time `gorm:"not null;column:jam_selesai"`
}

func (Kursus) TableName() string {
	return "kursus"
}

func (ImageKursus) TableName() string {
	return "image_kursus"
}

func (KategoriKursus) TableName() string {
	return "kategori_kursus"
}

func (MateriPembelajaran) TableName() string {
	return "materi_pembelajarans"
}

func (JadwalKursus) TableName() string {
	return "jadwal_kursus"
}
