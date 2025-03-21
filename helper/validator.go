package helper

import (
	"errors"
	"fmt"
	"regexp"
	"skripsi/constant"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	MinCost     int = 4
	MaxCost     int = 31
	DefaultCost int = 14
)

func ValidateEmail(email string) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$`)
	return regex.MatchString(email)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), DefaultCost)
	if err != nil {
		return "", constant.ErrHashPassword
	}
	return string(bytes), nil
}

func ValidatePassword(password string) (string, error) {
	// ^ => diawali
	// (?=.*[a-z]) => satu huruf kecil
	// (?=.*[A-Z]) => satu huruf besar
	// (?=.*\d) => satu angka
	// (?=.*[@$!%*?&#]) => satu karakter spesial
	// [A-Za-z\d@$!%*?&#]{8,} => karakter yang ada harus sesuai ketentuan diatas
	// .{8,} => minimal 8
	// $ => diakhiri
	// passwordValid := regexp.MustCompile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&#])[A-Za-z\d@$!%*?&#]{8,}$`).MatchString(password)

	containsLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	containsUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	containsNumber := regexp.MustCompile(`\d`).MatchString(password)
	containsSpecial := regexp.MustCompile(`[@$!%*?&#]`).MatchString(password)

	// Panjang password 8 - 16
	if len(password) < 8 || len(password) > 16 {
		return "", constant.ErrLenPassword
	}
	if !containsLower || !containsUpper || !containsNumber || !containsSpecial {
		return "", constant.ErrInvalidPassword
	}
	return password, nil
}

func TelephoneValidator(phone string) (string, error) {
	if len(phone) < 10 {
		return "", errors.New("invalid phone number format")
	}

	// No international
	phoneRegexCode := regexp.MustCompile(`^[+]{1}[0-9]{10,12}$`)
	if phoneRegexCode.MatchString(phone) {
		return phone, nil
	}

	if strings.HasPrefix(phone, "+62") {
		phone = "0" + phone[3:]
	}

	var phoneRegex = regexp.MustCompile(`^[0-9]{10,12}$`)
	if !phoneRegex.MatchString(phone) {
		return "", errors.New("nomor 10 sampai 12 digit ")
	}

	return phone, nil
}

func ValidateUsername(username string) (string, error) {
	var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9._-]{3,16}$`)
	if !usernameRegex.MatchString(username) {
		return "", errors.New("invalid username")
	}
	return username, nil
}
func ValidateTime(times string) (time.Time, error) {
	currentDate := time.Now().Format("2006-01-02")
	jamComplete := fmt.Sprintf("%s %s", currentDate, times)

	// Parse ke wib
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return time.Time{}, err
	}

	jamParsed, err := time.ParseInLocation("2006-01-02 15:04", jamComplete, loc)
	if err != nil {
		return time.Time{}, err
	}

	return jamParsed, nil
}

func CodeVoucherValidator(code string) (string, error) {
	var codeVoucher = regexp.MustCompile(`^[a-zA-Z0-9._-]{10}$`)
	if !codeVoucher.MatchString(code) {
		return "", errors.New("invalid username")
	}
	return code, nil
}

// =====================================
func ValidateTimeFormat(timeStr string) error {
	const timePattern = `^(?:[01]\d|2[0-3]):[0-5]\d$`
	match, _ := regexp.MatchString(timePattern, timeStr)
	if !match {
		return errors.New("invalid time format, expected HH:mm")
	}
	return nil
}

func ValidateTimeLogic(jamMulai, jamAkhir time.Time) error {
	if !jamMulai.Before(jamAkhir) {
		return errors.New("jam mulai harus lebih awal dari jam akhir")
	}
	return nil
}

func ValidateDateFormat(date string) error {
	const datePattern = `^\d{4}-\d{2}-\d{2}$`
	match, _ := regexp.MatchString(datePattern, date)
	if !match {
		return errors.New("invalid date format, expected YYYY-MM-DD")
	}
	return nil
}

func ValidateLogicalDate(tanggal time.Time) error {
	now := time.Now()
	if tanggal.Before(now) {
		return errors.New("jadwal tidak boleh di masa lalu")
	}
	return nil
}
