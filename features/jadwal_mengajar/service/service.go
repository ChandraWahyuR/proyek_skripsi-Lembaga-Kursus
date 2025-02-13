package service

import (
	"skripsi/constant"
	jadwal "skripsi/features/jadwal_mengajar"
	"skripsi/helper"
)

type JadwalService struct {
	d jadwal.MengajarRepositoryInterface
	j helper.JWTInterface
}

func NewServiceCatatan(d jadwal.MengajarRepositoryInterface, j helper.JWTInterface) jadwal.MengajarServiceInterface {
	return &JadwalService{
		d: d,
		j: j,
	}
}

func (s *JadwalService) CreateJadwalMengajar(data *jadwal.JadwalMengajar) error {
	switch {
	case data.InstrukturID == "":
		return constant.ErrEmptyNameInstuktor
	case data.KursusID == "":
		return constant.ErrKursusName
	case data.Tanggal.IsZero():
		return constant.ErrJadwal
	case data.JamMulai.IsZero():
		return constant.ErrJamMulai
	case data.JamAkhir.IsZero():
		return constant.ErrJamAkhir
	}

	// Validasi format tanggal dan logika waktu
	if err := helper.ValidateDateFormat(data.Tanggal.Format("2006-01-02")); err != nil {
		return err
	}

	if err := helper.ValidateLogicalDate(data.Tanggal); err != nil {
		return err
	}

	if err := helper.ValidateTimeFormat(data.JamMulai.Format("15:04")); err != nil {
		return err
	}
	if err := helper.ValidateTimeFormat(data.JamAkhir.Format("15:04")); err != nil {
		return err
	}

	if err := helper.ValidateTimeLogic(data.JamMulai, data.JamAkhir); err != nil {
		return err
	}

	return s.d.CreateJadwalBatch(data)
}

func (s *JadwalService) EditJadwalMengajar(data *jadwal.JadwalMengajar) error {
	if data.ID == "" {
		return constant.ErrEmptyId
	}
	if data.InstrukturID == "" && data.Tanggal.IsZero() && data.JamAkhir.IsZero() && data.JamMulai.IsZero() {
		return constant.ErrUpdate
	}

	return s.d.EditJadwalMengajar(data)
}

func (s *JadwalService) DeleteJadwalMengajar(id string) error {
	return s.d.DeleteJadwalMengajar(id)
}

func (s *JadwalService) GetJadwalMengajar() ([]*jadwal.JadwalMengajar, error) {
	return s.d.GetJadwalMengajar()
}
func (s *JadwalService) GetJadwalMengajarByID(id string) (*jadwal.JadwalMengajar, error) {
	if id == "" {
		return &jadwal.JadwalMengajar{}, constant.ErrEmptyId
	}
	return s.d.GetJadwalMengajarByID(id)
}
func (s *JadwalService) GetJadwalMengajarForUser(user_id string) ([]*jadwal.JadwalMengajar, error) {
	return s.d.GetJadwalMengajarForUser(user_id)
}
