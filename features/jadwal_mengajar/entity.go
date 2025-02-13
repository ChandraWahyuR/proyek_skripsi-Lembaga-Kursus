package jadwal

import (
	"time"

	"github.com/labstack/echo/v4"
)

type JadwalMengajar struct {
	ID           string
	InstrukturID string
	Instruktur   Instruktur
	UserID       string
	User         User
	KursusID     string
	Kursus       Kursus
	Tanggal      time.Time
	JamMulai     time.Time
	JamAkhir     time.Time
	Status       bool
}

type Instruktur struct {
	ID             string
	InstrukturNama string
}
type User struct {
	ID       string
	UserName string
}
type Kursus struct {
	ID         string
	KursusNama string
}

type FeedbackMengajar struct {
	ID               string
	UserID           string
	JadwalMengajarID string
	Penilaian        int64
	Deskripsi        string
}

type MengajarHandlerInterface interface {
	GetJadwalMengajar() echo.HandlerFunc
	GetJadwalMengajarByID() echo.HandlerFunc
	GetJadwalMengajarForUser() echo.HandlerFunc
	CreateJadwalMengajar() echo.HandlerFunc
	EditJadwalMengajar() echo.HandlerFunc
}

type MengajarServiceInterface interface {
	GetJadwalMengajar() ([]*JadwalMengajar, error)
	GetJadwalMengajarByID(id string) (*JadwalMengajar, error)
	GetJadwalMengajarForUser(user_id string) ([]*JadwalMengajar, error)
	CreateJadwalMengajar(data *JadwalMengajar) error
	EditJadwalMengajar(data *JadwalMengajar) error
	DeleteJadwalMengajar(id string) error
}

type MengajarRepositoryInterface interface {
	GetJadwalMengajar() ([]*JadwalMengajar, error)
	GetJadwalMengajarByID(id string) (*JadwalMengajar, error)
	GetJadwalMengajarForUser(user_id string) ([]*JadwalMengajar, error)
	CreateJadwalMengajar(data *JadwalMengajar) error
	EditJadwalMengajar(data *JadwalMengajar) error
	DeleteJadwalMengajar(id string) error
	//
	CreateJadwalBatch(data *JadwalMengajar) error
}
