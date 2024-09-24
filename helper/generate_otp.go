package helper

import (
	"math/rand"
	"time"
)

func GenerateOTP(length int) string {
	rand.Seed(time.Now().UnixNano())
	otpChars := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	otp := make([]byte, length)
	for i := range otp {
		otp[i] = otpChars[rand.Intn(len(otpChars))]
	}
	return string(otp)
}
