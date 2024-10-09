package helper

import (
	"math/rand"
	"strconv"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateCode() string {
	// Seed acak
	rand.Seed(time.Now().UnixNano())

	// Fungsi untuk menghasilkan karakter acak dari charset
	generateRandomString := func(length int) string {
		b := make([]byte, length)
		for i := range b {
			b[i] = charset[rand.Intn(len(charset))]
		}
		return string(b)
	}

	// Menghasilkan string acak sepanjang 8 karakter
	randomString := generateRandomString(8)

	// Mendapatkan 2 digit terakhir dari tahun saat ini
	currentYear := time.Now().Year()
	yearSuffix := strconv.Itoa(currentYear)[2:] // Ambil 2 digit terakhir dari tahun

	// Menggabungkan string acak dengan 2 digit tahun
	code := randomString + yearSuffix

	return code
}
