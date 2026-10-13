package config

import (
	"log"
	"os"

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
}

// -------- Config --------
type Config struct {
	App    App
	DB     Database
	SMSNet SMSNet
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
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
		},
	}
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Error loading env variable: %s", key)
	}
	return value
}
