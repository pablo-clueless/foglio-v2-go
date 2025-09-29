package lib

import (
	"crypto/rand"
	"errors"
	"foglio/v2/src/config"
	"math/big"
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

func GenerateUrl(token string) (string, error) {
	if client := config.AppConfig.ClientUrl; client != "" {
		return client + "?token=" + token, nil
	}
	return "", errors.New("client url not provided in env")
}
