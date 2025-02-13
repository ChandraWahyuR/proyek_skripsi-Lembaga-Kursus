package data

import (
	"errors"
	"fmt"
	"skripsi/constant"
	jadwal "skripsi/features/jadwal_mengajar"
	"skripsi/features/notification/sse"
	"skripsi/features/transaksi"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JadwalRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) jadwal.MengajarRepositoryInterface {
	return &JadwalRepository{
		db: db,
	}
}

func (d *JadwalRepository) GetJadwalMengajar() ([]*jadwal.JadwalMengajar, error) {
	var jadwalList []*JadwalMengajar

	if err := d.db.Preload("User").
		Preload("Kursus").
		Preload("Instruktur").
		Where("deleted_at IS NULL").
		Where("status IS TRUE").
		Find(&jadwalList).Error; err != nil {

		return nil, constant.ErrDataNotfound
	}

	var result []*jadwal.JadwalMengajar
	for _, item := range jadwalList {
		result = append(result, item.ToEntity())
	}

	return result, nil
}
func (d *JadwalRepository) GetJadwalMengajarByID(id string) (*jadwal.JadwalMengajar, error) {
	var jadwalList *JadwalMengajar

	if err := d.db.Preload("User").
		Preload("Kursus").
		Preload("Instruktur").
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Where("status IS TRUE").
		First(&jadwalList).Error; err != nil {

		return nil, constant.ErrDataNotfound
	}

	result := jadwalList.ToEntity()
	return result, nil
}

func (d *JadwalRepository) GetJadwalMengajarForUser(user_id string) ([]*jadwal.JadwalMengajar, error) {
	var jadwalList []*JadwalMengajar

	if err := d.db.Preload("User").
		Preload("Kursus").
		Preload("Instruktur").
		Where("user_id = ?", user_id).
		Where("deleted_at IS NULL").
		Where("status IS TRUE").
		Find(&jadwalList).Error; err != nil {

		return nil, constant.ErrDataNotfound
	}

	var result []*jadwal.JadwalMengajar
	for _, item := range jadwalList {
		result = append(result, item.ToEntity())
	}

	return result, nil
}

func (d *JadwalRepository) CreateJadwalBatch(data *jadwal.JadwalMengajar) error {
	const MaxMahasiswa = 30

	// Query transaksi mahasiswa aktif
	var validUsers []transaksi.TransaksiHistory
	if err := d.db.Where("kursus_id = ?", data.KursusID).
		Where("valid_until > ?", time.Now()).
		Where("status = ?", "Active").
		Limit(MaxMahasiswa).
		Find(&validUsers).Error; err != nil {
		return err
	}

	if len(validUsers) == 0 {
		return errors.New("no valid users for this course")
	}

	tx := d.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Loop untuk insert jadwal
	for _, user := range validUsers {
		newJadwal := FromEntity(data)
		newJadwal.ID = uuid.New().String()
		newJadwal.UserID = user.UserID
		newJadwal.Status = true

		if err := tx.Create(&newJadwal).Error; err != nil {
			tx.Rollback()
			return err
		}
		// Kirim notifikasi SSE
		go func(userID, tanggal, jamMulai, jamAkhir string) {
			fmt.Printf("data: Jadwal baru pada %s pukul %s - %s untuk user %s\n\n", tanggal, jamMulai, jamAkhir, userID)
			sse.SendSSENotification(userID, tanggal, jamMulai, jamAkhir)
		}(user.UserID, newJadwal.Tanggal.Format("2006-01-02"), newJadwal.JamMulai.Format("15:04"), newJadwal.JamAkhir.Format("15:04"))

	}

	return tx.Commit().Error
}

func (d *JadwalRepository) CreateJadwalMengajar(data *jadwal.JadwalMengajar) error {
	tx := d.db.Begin()

	if err := tx.Create(&data).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (d *JadwalRepository) EditJadwalMengajar(data *jadwal.JadwalMengajar) error {
	tx := d.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	model := FromEntity(data)
	if err := tx.Model(&JadwalMengajar{}).Where("id = ?", data.ID).Updates(model).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (d *JadwalRepository) DeleteJadwalMengajar(id string) error {
	tx := d.db.Begin()

	if err := tx.Where("deleted_at IS NULL").Where("id = ?", id).Delete(&JadwalMengajar{}); err.Error != nil {
		tx.Rollback()
		return constant.ErrGetID
	} else if err.RowsAffected == 0 {
		tx.Rollback()
		return constant.ErrGetID
	}

	return tx.Commit().Error
}

func (d *JadwalRepository) ReviewPengajar(data *jadwal.FeedbackMengajar) error {
	// Konversi entity ke model
	model := FromEntityFeedback(data)

	if err := d.db.Create(model).Error; err != nil {
		return constant.ErrBadRequest
	}
	return nil
}
