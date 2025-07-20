package utils

import "math/rand"

func GenerateOTP(length int) string {
	const charset = "1234567890"
	otp := make([]byte, length)
	for i := range otp {
		otp[i] = charset[rand.Intn(len(charset))]
	}
	return string(otp)
}
