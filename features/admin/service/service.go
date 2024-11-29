package service

import (
	"encoding/csv"
	"fmt"
	"os"
	"skripsi/constant"
	"skripsi/features/admin"
	"skripsi/helper"
	"strconv"
	"strings"
	"time"
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

func (s *AdminService) DownloadLaporanPembelian(startDate, endDate time.Time) (string, error) {
	histories, err := s.d.DownloadLaporanPembelian(startDate, endDate)
	if err != nil {
		return "", err
	}

	filename := "laporan_pembelian_" + startDate.Format("2006-01-02") + "_to_" + endDate.Format("2006-01-02") + ".csv"
	file, err := os.CreateTemp("", filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Menambahkan header laporan
	writer.Write([]string{"Laporan Pembelian"})
	writer.Write([]string{"Periode", startDate.Format("02 January 2006") + " - " + endDate.Format("02 January 2006")})
	writer.Write([]string{})

	// Header tabel
	writer.Write([]string{"ID", "TransaksiID", "KursusID", "UserID", "UserNama", "Email", "Nama Kursus", "Status Pembelian", "ValidUntil", "TotalHarga", "Status Transaksi"})

	var totalUser int
	var totalHarga float64

	// Tulis data dari map ke CSV dan hitung total harga
	for _, history := range histories {
		var totalHargaStr string
		var harga float64

		if th, ok := history["total_harga"].(float64); ok {
			totalHargaStr = fmt.Sprintf("%.2f", th)
			harga = th
		} else if th, ok := history["total_harga"].(string); ok {
			if floatVal, err := strconv.ParseFloat(th, 64); err == nil {
				totalHargaStr = fmt.Sprintf("%.2f", floatVal)
				harga = floatVal
			} else {
				totalHargaStr = "0.00"
				harga = 0.0
			}
		} else {
			totalHargaStr = "0.00"
			harga = 0.0
		}

		// Tulis baris data ke CSV
		writer.Write([]string{
			history["id"].(string),
			history["transaksi_id"].(string),
			history["kursus_id"].(string),
			history["user_id"].(string),
			history["user_nama"].(string),
			history["email"].(string),
			history["nama_kursus"].(string),
			history["status"].(string),
			history["valid_until"].(time.Time).Format("2006-01-02"),
			totalHargaStr,
			history["transaksi_status"].(string),
		})

		totalUser++

		// Hanya total transaksi yang statusnya sudah aktif dan pembayarannya sukses
		if history["status"] == "Active" && history["transaksi_status"] == "Success" {
			totalHarga += harga
		}
	}

	// Footer
	writer.Write([]string{})
	writer.Write([]string{"Total User", fmt.Sprintf("%d", totalUser)})
	writer.Write([]string{"Total Dana yang Masuk", fmt.Sprintf("%.2f", totalHarga)})

	return filename, nil
}
