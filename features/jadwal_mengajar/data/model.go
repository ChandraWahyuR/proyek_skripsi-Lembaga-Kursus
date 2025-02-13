package data

import (
	Instruktur "skripsi/features/instruktur/data"
	jadwal "skripsi/features/jadwal_mengajar"
	Kursus "skripsi/features/kursus/data"
	User "skripsi/features/users/data"
	"time"

	"gorm.io/gorm"
)

type JadwalMengajar struct {
	*gorm.Model
	ID           string                `gorm:"type:varchar(50);primaryKey;not null;column:id"`
	InstrukturID string                `gorm:"type:varchar(255);not null;column:instruktur_id"`
	Instruktur   Instruktur.Instruktur `gorm:"foreignKey:InstrukturID;references:ID"`
	UserID       string                `gorm:"type:varchar(255);not null;column:user_id"`
	User         User.User             `gorm:"foreignKey:UserID;references:ID"`
	KursusID     string                `gorm:"type:varchar(255);not null;column:kursus_id"`
	Kursus       Kursus.Kursus         `gorm:"foreignKey:KursusID;references:ID"`
	Tanggal      time.Time             `gorm:"type:date"`
	JamMulai     time.Time             `gorm:"type:time"`
	JamAkhir     time.Time             `gorm:"type:time"`
	Status       bool                  `gorm:"default:true"`
}

type FeedbackMengajar struct {
	*gorm.Model
	ID               string         `gorm:"type:varchar(50);primaryKey;not null;column:id"`
	UserID           string         `gorm:"type:varchar(255);not null;column:user_id"`
	User             User.User      `gorm:"foreignKey:UserID;references:ID"`
	JadwalMengajarID string         `gorm:"type:varchar(255);not null;column:jadwal_mengajar"`
	JadwalMengajar   JadwalMengajar `gorm:"foreignKey:JadwalMengajarID;references:ID"`
	Penilaian        int64          `gorm:"type:int64(10);not null;column:penilaian"`
	Deskripsi        string         `gorm:"type:varchar(255);not null;column:deskripsi"`
}

func (JadwalMengajar) TableName() string {
	return "jadwal_mengajars"
}

func (FeedbackMengajar) TableName() string {
	return "feedback_mengajar"
}

func (j *JadwalMengajar) ToEntity() *jadwal.JadwalMengajar {
	return &jadwal.JadwalMengajar{
		ID:           j.ID,
		InstrukturID: j.InstrukturID,
		Instruktur: jadwal.Instruktur{
			ID:             j.InstrukturID,
			InstrukturNama: j.Instruktur.Name,
		},
		UserID: j.UserID,
		User: jadwal.User{
			ID:       j.UserID,
			UserName: j.User.Username,
		},
		KursusID: j.KursusID,
		Kursus: jadwal.Kursus{
			ID:         j.KursusID,
			KursusNama: j.Kursus.Nama,
		},
		Tanggal:  j.Tanggal,
		JamMulai: j.JamMulai,
		JamAkhir: j.JamAkhir,
		Status:   j.Status,
	}
}

func FromEntity(e *jadwal.JadwalMengajar) *JadwalMengajar {
	return &JadwalMengajar{
		ID:           e.ID,
		InstrukturID: e.InstrukturID,
		UserID:       e.UserID,
		KursusID:     e.KursusID,
		Tanggal:      e.Tanggal,
		JamMulai:     e.JamMulai,
		JamAkhir:     e.JamAkhir,
		Status:       e.Status,
	}
}

func FromEntityFeedback(e *jadwal.FeedbackMengajar) *FeedbackMengajar {
	return &FeedbackMengajar{
		ID:               e.ID,
		UserID:           e.UserID,
		JadwalMengajarID: e.JadwalMengajarID,
		Penilaian:        e.Penilaian,
		Deskripsi:        e.Deskripsi,
	}
}
