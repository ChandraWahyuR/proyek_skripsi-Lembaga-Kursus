package service

import (
	"encoding/csv"
	"fmt"
	"io"
	"skripsi/constant"
	"skripsi/features/admin"
	"skripsi/helper"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type AdminService struct {
	d admin.AdminDataInterface
	j helper.JWTInterface
}

func New(d admin.AdminDataInterface, j helper.JWTInterface) admin.AdminServiceInterface {
	return &AdminService{
		d: d,
		j: j,
	}
}

func (s *AdminService) RegisterAdmin(user admin.Admin) error {
	switch {
	case user.Email == "":
		return constant.ErrEmptyEmailRegister
	case user.Username == "":
		return constant.ErrEmptyNameRegister
	case user.Password == "":
		return constant.ErrEmptyPasswordRegister
	case user.Password != user.ConfirmPassword:
		return constant.ErrPasswordNotMatch
	}
	user.Email = strings.ToLower(user.Email)

	isEmailValid := helper.ValidateEmail(user.Email)
	if !isEmailValid {
		return constant.ErrInvalidEmail
	}

	isUsernameValid, err := helper.ValidateUsername(user.Username)
	if err != nil {
		return constant.ErrInvalidUsername
	}

	pass, err := helper.ValidatePassword(user.Password)
	if err != nil {
		return constant.ErrInvalidPassword
	}

	hashedPassword, err := helper.HashPassword(pass)
	if err != nil {
		return err
	}
	user.Username = isUsernameValid
	user.Password = hashedPassword

	err = s.d.RegisterAdmin(user)
	if err != nil {
		return err
	}
	return nil
}

func (s *AdminService) LoginAdmin(user admin.Admin) (admin.Login, error) {
	switch {
	case user.Username == "":
		return admin.Login{}, constant.ErrAdminUserNameEmpty
	case user.Password == "":
		return admin.Login{}, constant.ErrAdminPasswordEmpty
	}

	adminData, err := s.d.LoginAdmin(user)
	if err != nil {
		return admin.Login{}, err
	}

	var adminLogin helper.AdminJWT
	adminLogin.ID = adminData.ID
	adminLogin.Username = adminData.Username
	adminLogin.Email = adminData.Email

	token, err := s.j.GenerateAdminJWT(adminLogin)
	if err != nil {
		return admin.Login{}, err
	}

	var adminLoginData admin.Login
	adminLoginData.Token = token

	return adminLoginData, nil
}

func (s *AdminService) DownloadLaporanPembelian(startDate, endDate time.Time) ([]map[string]interface{}, error) {
	return s.d.DownloadLaporanPembelian(startDate, endDate)
}

func (s *AdminService) GenerateLaporanCSV(w io.Writer, histories []map[string]interface{}, startDate, endDate time.Time) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()
	writer.Write([]string{"Laporan Siswa Aktif \nLKP Mediakom Sidareja"})
	writer.Write([]string{"Periode", fmt.Sprintf("%s - %s", startDate.Format("02 January 2006"), endDate.Format("02 January 2006"))})
	writer.Write([]string{})
	writer.Write([]string{"No", "Nomor Induk", "Transaksi ID", "User ID", "Username", "Nama", "Jenis Kelamin", "Email", "Nomor Telepon", "Alamat", "Nama Kursus", "Tanggal Masuk", "Berlaku sampai", "Total Harga"})
	var totalUser int
	var totalHarga float64
	var userAktif int
	formatter := message.NewPrinter(language.Indonesian)
	for i, history := range histories {
		var harga float64
		if th, ok := history["total_harga"].(float64); ok {
			harga = th
		} else if th, ok := history["total_harga"].(string); ok {
			if floatVal, err := strconv.ParseFloat(th, 64); err == nil {
				harga = floatVal
			} else {
				harga = 0.0
			}
		} else {
			harga = 0.0
		}
		totalHargaStr := formatter.Sprintf("Rp.%d", int64(harga))
		writer.Write([]string{
			fmt.Sprintf("%d", i+1),
			history["nis"].(string),
			history["transaksi_id"].(string),
			history["user_id"].(string),
			history["username"].(string),
			history["jenis_kelamin"].(string),
			history["nama"].(string),
			history["email"].(string),
			history["hp"].(string),
			history["alamat"].(string),
			history["nama_kursus"].(string),
			history["tgl_masuk"].(time.Time).Format("2006-01-02"),
			history["valid_until"].(time.Time).Format("2006-01-02"),
			totalHargaStr,
		})

		totalUser++
		if history["status"] == "Active" && history["transaksi_status"] == "Success" {
			totalHarga += harga
			userAktif++
		}
	}

	hasil := formatter.Sprintf("Rp.%d", int64(totalHarga))
	writer.Write([]string{})
	writer.Write([]string{"Total Pengguna Yang Mendaftar", fmt.Sprintf("%d", totalUser)})
	writer.Write([]string{"Total Pemasukkan", hasil})
	return nil
}
