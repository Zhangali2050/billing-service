package config

import (
	"os"
)

type Config struct {
	DatabaseURL     string
	AirbaURL        string
	AirbaUser       string
	AirbaPassword   string
	AirbaTerminalID string
}

func Load() (*Config, error) {
	return &Config{
		DatabaseURL:     getEnv("DATABASE_URL", "postgres://postgres:devpassword@billing_db:5432/dev_db?sslmode=disable"),
		AirbaURL:        getEnv("AIRBA_URL", "https://sandbox.airbapay.kz"),
		AirbaUser:       getEnv("AIRBA_USER", "sandbox_user"),
		AirbaPassword:   getEnv("AIRBA_PASSWORD", "sandbox_pass"),
		AirbaTerminalID: getEnv("AIRBA_TERMINAL_ID", "sandbox_terminal"),
	}, nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
