package config

import (
	"math/rand"
	"net/http"
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
	GithubPrivateKey     string
	GithubRedirectUrl    string
	GoogleClientId       string
	GoogleClientSecret   string
	GoogleRedirectUrl    string
	IsDevMode            bool
	JWTTokenSecret       []byte
	MaxFileSize          int
	NonAuthRoutes        []APIRoute
	Port                 string
	PostgresUrl          string
	ProjectId            string
	RedisUrl             string
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
		ClientUrl:            "CLIENT_URL",
		CloudinaryKey:        os.Getenv("CLOUDINARY_KEY"),
		CloudinaryName:       os.Getenv("CLOUDINARY_NAME"),
		CloudinarySecret:     os.Getenv("CLOUDINARY_SECRET"),
		CookieDomain:         os.Getenv("COOKIE_DOMAIN"),
		CurrentUserId:        "CURRENT_USER_ID",
		CurrentUser:          "CURRENT_USER",
		GithubClientId:       os.Getenv("GITHUB_CLIENT_ID"),
		GithubClientSecret:   os.Getenv("GITHUB_CLIENT_SECRET"),
		GithubPrivateKey:     os.Getenv("GITHUB_CLIENT_PRIVATE_KEY"),
		GithubRedirectUrl:    os.Getenv("GITHUB_REDIRECT_URL"),
		GoogleClientId:       os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret:   os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectUrl:    os.Getenv("GOOGLE_REDIRECT_URL"),
		IsDevMode:            os.Getenv("GO_ENV") == "development",
		JWTTokenSecret:       []byte(os.Getenv("JWT_TOKEN_SECRET")),
		MaxFileSize:          10 << 20, // default 10 MB
		Port:                 os.Getenv("PORT"),
		PostgresUrl:          os.Getenv("POSTGRES_URL"),
		ProjectId:            os.Getenv("PROJECT_ID"),
		RedisUrl:             os.Getenv("REDIS_URL"),
		SmtpHost:             os.Getenv("SMTP_HOST"),
		SmtpPort:             func() int { p, _ := strconv.Atoi(os.Getenv("SMTP_PORT")); return p }(),
		SmtpUser:             os.Getenv("SMTP_USER"),
		SmtpPassword:         os.Getenv("SMTP_PASSWORD"),
		Version:              os.Getenv("VERSION"),
		NonAuthRoutes: []APIRoute{
			{Endpoint: "/public/*", Method: "*"},
			{Endpoint: "/swagger/*", Method: "*"},
			{Endpoint: "/docs", Method: http.MethodGet},
			{Endpoint: "/", Method: http.MethodGet},
			{Endpoint: "api/v2", Method: http.MethodGet},
			{Endpoint: "api/v2/ws", Method: http.MethodGet},
			{Endpoint: "/api/v2/ws/stats", Method: http.MethodGet},
			{Endpoint: "/api/v2/ws/send-notification", Method: http.MethodPost},
			{Endpoint: "/api/v2/ws/broadcast", Method: http.MethodPost},
			{Endpoint: "api/v2/health", Method: http.MethodGet},
			{Endpoint: "api/v2/auth/signup", Method: http.MethodPost},
			{Endpoint: "api/v2/auth/signin", Method: http.MethodPost},
			{Endpoint: "api/v2/auth/verification", Method: http.MethodPost},
			{Endpoint: "api/v2/auth/forgot-password", Method: http.MethodPost},
			{Endpoint: "api/v2/auth/reset-password", Method: http.MethodPost},
			{Endpoint: "api/v2/auth/github", Method: http.MethodGet},
			{Endpoint: "api/v2/auth/github/callback", Method: http.MethodGet},
			{Endpoint: "api/v2/auth/google", Method: http.MethodGet},
			{Endpoint: "api/v2/auth/google/callback", Method: http.MethodGet},
			{Endpoint: "api/v2/users", Method: http.MethodGet},
			{Endpoint: "api/v2/users/:id", Method: http.MethodGet},
			{Endpoint: "api/v2/jobs", Method: http.MethodGet},
			{Endpoint: "api/v2/jobs/:id", Method: http.MethodGet},
		},
	}
}
