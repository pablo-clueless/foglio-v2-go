package config

import (
	"math/rand"
	"os"
	"strconv"
	"time"
)

type APIRoute struct {
	Endpoint string
	Method   string
}

type Config struct {
	AccessTokenExpiresIn time.Duration
	AllowMethods         []string
	AppEmail             string
	ApiUrl               string
	ClientUrl            string
	CloudinaryKey        string
	CloudinaryName       string
	CloudinarySecret     string
	CookieDomain         string
	CurrentUser          string
	CurrentUserId        string
	GithubClientId       string
	GithubClientSecret   string
	GoogleClientId       string
	GoogleClientSecret   string
	IsDevMode            bool
	JWTTokenSecret       []byte
	MaxFileSize          int
	NonAuthRoutes        []APIRoute
	Port                 string
	PostgresUrl          string
	ProjectId            string
	SmtpHost             string
	SmtpPort             int
	SmtpUser             string
	SmtpPassword         string
	Version              string
}

var AppConfig *Config

func InitializeConfig() {
	rand.Seed(time.Now().UnixNano())
	InitializeEnvFile()
	AppConfig = &Config{
		AccessTokenExpiresIn: time.Minute * 30,
		AppEmail:             os.Getenv("APP_EMAIL"),
		ApiUrl:               os.Getenv("API_URL"),
		ClientUrl:            os.Getenv("CLIENT_URL"),
		CloudinaryKey:        os.Getenv("CLOUDINARY_KEY"),
		CloudinaryName:       os.Getenv("CLOUDINARY_NAME"),
		CloudinarySecret:     os.Getenv("CLOUDINARY_SECRET"),
		CookieDomain:         os.Getenv("COOKIE_DOMAIN"),
		CurrentUserId:        "CURRENT USER ID",
		CurrentUser:          "CURRENT USER",
		GithubClientId:       os.Getenv("GITHUB_CLIENT_ID"),
		GithubClientSecret:   os.Getenv("GITHUB_CLIENT_SECRET"),
		GoogleClientId:       os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret:   os.Getenv("GOOGLE_CLIENT_SECRET"),
		IsDevMode:            os.Getenv("GO_ENV") == "development",
		JWTTokenSecret:       []byte(os.Getenv("JWT_TOKEN_SECRET")),
		MaxFileSize:          10 << 20, // default 10 MB
		Port:                 os.Getenv("PORT"),
		PostgresUrl:          os.Getenv("POSTGRES_URL"),
		ProjectId:            os.Getenv("PROJECT_ID"),
		SmtpHost:             os.Getenv("SMTP_HOST"),
		SmtpPort:             func() int { p, _ := strconv.Atoi(os.Getenv("SMTP_PORT")); return p }(),
		SmtpUser:             os.Getenv("SMTP_USER"),
		SmtpPassword:         os.Getenv("SMTP_PASSWORD"),
		Version:              os.Getenv("VERSION"),
		AllowMethods:         []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		NonAuthRoutes:        []APIRoute{},
	}
}
