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

func (c *Config) Validate() error {
	if c.App.Port == "" {
		return fmt.Errorf("APP PORT is required")
	}

	if c.Database.URL == "" {
		return fmt.Errorf("DATABASE URL is required")
	}

	if c.App.JwtSecret == "" {
		return fmt.Errorf("JWT secret is required")
	}

	if c.App.Service.Payment.Vendor != "xendit" {
		return fmt.Errorf("only xendit payment gateway is currently supported")
	}

	if c.App.Service.Payment.SecretKey == "" {
		return fmt.Errorf("Secret key is required")
	}

	return nil
}

type AppConfig struct {
	Port      string
	URL       string
	ENV       string
	Service   ServiceConfig
	JwtSecret string
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
	_ = godotenv.Load(".env")

	return &Config{
		App: AppConfig{
			Port:      getEnv("APP_PORT", ""),
			URL:       getEnv("APP_URL", ""),
			ENV:       getEnv("APP_ENV", "development"),
			JwtSecret: getEnv("JWT_SECRET", ""),
			Service: ServiceConfig{
				Payment: PaymentConfig{
					Vendor:    getEnv("PAYMENT_VENDOR", "xendit"),
					SecretKey: getEnv("PAYMENT_SECRET_KEY", ""),
				},
			},
		},
		Database: DatabaseConfig{
			URL:  getEnv("DATABASE_URL", ""),
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
