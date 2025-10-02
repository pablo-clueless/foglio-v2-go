package lib

import (
	"crypto/rand"
	"math/big"
	"strings"

	"github.com/google/uuid"
)

func GenerateOtp() string {
	const charset = "0123456789"
	otp := make([]byte, 6)
	for i := range otp {
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		otp[i] = charset[randomIndex.Int64()]
	}
	return string(otp)
}

func GenerateUrl(client, token string) string {
	return client + "?token=" + token
}

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[randomIndex.Int64()]
	}
	return string(b)
}

func GenerateUUID() string {
	return uuid.New().String()
}

func GenerateUsername(name string) string {
	base := strings.ToLower(strings.ReplaceAll(name, " ", "."))
	return base + "." + GenerateRandomString(6)
}
