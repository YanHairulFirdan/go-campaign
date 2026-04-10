package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
}

type AppConfig struct {
	Port    string
	URL     string
	ENV     string
	Service ServiceConfig
}

type ServiceConfig struct {
	Payment PaymentConfig
}

type PaymentConfig struct {
	Vendor    string
	SecretKey string
}

type DatabaseConfig struct {
	URL  string
	Type string
}

func Load() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	return &Config{
		App: AppConfig{
			Port: getEnv("APP_PORT", "3030"),
			URL:  getEnv("APP_URL", "http://localhost:3030"),
			ENV:  getEnv("APP_ENV", "development"),
			Service: ServiceConfig{
				Payment: PaymentConfig{
					Vendor:    getEnv("PAYMENT_VENDOR", "xendit"),
					SecretKey: getEnv("PAYMENT_SECRET_KEY", ""),
				},
			},
		},
		Database: DatabaseConfig{
			URL:  getEnv("DATABASE_URL", "postgres://user:pass@localhost:5432/dbname"),
			Type: getEnv("DATABASE_TYPE", "postgres"),
		},
	}, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
