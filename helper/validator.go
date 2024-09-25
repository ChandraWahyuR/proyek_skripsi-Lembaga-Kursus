package helper

import (
	"errors"
	"regexp"
	"skripsi/constant"

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
		return "", errors.New("password must be between 8 and 16 characters long")
	}
	if !containsLower || !containsUpper || !containsNumber || !containsSpecial {
		return "", errors.New("password must contain at least 1 number, 1 uppercase letter, one punctuation symbol and 1 lowercase letter")
	}
	return password, nil
}

func TelephoneValidator(phone string) (string, error) {
	phoneRegexCode := regexp.MustCompile(`^[+]{1}[0-9]{10,12}$`)

	if !phoneRegexCode.MatchString(phone) {
		phone = phoneRegexCode.ReplaceAllString(phone, "0"+phone[2:])
	}
	var phoneRegex = regexp.MustCompile(`^[0-9]{10,12}$`)
	if !phoneRegex.MatchString(phone) {
		return "", errors.New("invalid phone number")
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
