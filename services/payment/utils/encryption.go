package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"

	"github.com/Creative-genius001/Stacklo/services/payment/config"
	"github.com/Creative-genius001/go-logger"
)

func GenerateHMAC(payload string) string {
	secretKey := config.Cfg.HMACKey
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(payload))
	signature := hex.EncodeToString(h.Sum(nil))

	logger.Info("HMAC signature", signature)
	return signature
}

func ValidateHMAC(payload, expectedHMAC string) bool {
	hmacSignature := GenerateHMAC(payload)
	return hmacSignature == expectedHMAC
}
