package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type App struct {
	Name string
	Port string
	Env  string
}

type Database struct {
	Url            string
	MigrationsPath string
	RedisUrl       string
}

type SMSNet struct {
	APIKey      string
	ExpireInMin string
	Enabled     bool
}

type Auth struct {
	JWTSecret       string
	OTPSecret       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type Firebase struct {
	CredentialsBase64 string
}

type Brevo struct {
	APIKey      string
	SenderEmail string
	SenderName  string
}

// -------- Config --------
type Config struct {
	App      App
	DB       Database
	SMSNet   SMSNet
	Auth     Auth
	Firebase Firebase
	Brevo    Brevo
}

const minSecretLength = 32

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cnf := &Config{
		App: App{
			Name: getEnv("APP_NAME"),
			Port: getEnv("PORT"),
			Env:  getEnv("ENV"),
		},
		DB: Database{
			Url:            getEnv("DATABASE_URL"),
			MigrationsPath: getEnv("DB_MIGRATIONS_PATH"),
			RedisUrl:       getEnv("REDIS_URL"),
		},
		SMSNet: SMSNet{
			APIKey:      getEnv("SMS_NET_API_KEY"),
			ExpireInMin: getEnv("SMS_OTP_EXPIRE_IN_MIN"),
			Enabled:     getEnvOptional("SMS_ENABLED", "true") == "true",
		},
		Auth: Auth{
			JWTSecret:       getEnv("JWT_SECRET"),
			OTPSecret:       getEnv("OTP_SECRET"),
			AccessTokenTTL:  time.Duration(getEnvInt("ACCESS_TOKEN_TTL_MIN", 15)) * time.Minute,
			RefreshTokenTTL: time.Duration(getEnvInt("REFRESH_TOKEN_TTL_DAYS", 30)) * 24 * time.Hour,
		},
		Firebase: Firebase{
			CredentialsBase64: getEnv("FIREBASE_CREDENTIALS_BASE64"),
		},
		Brevo: Brevo{
			APIKey:      getEnvOptional("BREVO_API_KEY", ""),
			SenderEmail: getEnvOptional("BREVO_SENDER_EMAIL", ""),
			SenderName:  getEnvOptional("BREVO_SENDER_NAME", "GPI Connect"),
		},
	}

	if len(cnf.Auth.JWTSecret) < minSecretLength || len(cnf.Auth.OTPSecret) < minSecretLength {
		log.Fatalf("JWT_SECRET and OTP_SECRET must be at least %d characters, generate them with: go run ./cmd/tools/keygen", minSecretLength)
	}

	return cnf
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Error loading env variable: %s", key)
	}
	return value
}

func getEnvOptional(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	number, err := strconv.Atoi(value)
	if err != nil || number <= 0 {
		log.Fatalf("Env variable %s must be a positive number", key)
	}
	return number
}
