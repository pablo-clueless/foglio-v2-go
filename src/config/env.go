package config

import (
	"os"

	"github.com/joho/godotenv"
)

func InitializeEnvFile() error {
	env := os.Getenv("GO_ENV")
	if env == "" || env == "development" {
		err := godotenv.Load()
		if err != nil {
			return err
		}
	}
	return nil
}
