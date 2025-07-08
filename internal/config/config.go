package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	Airba       AirbaConfig
}

type AirbaConfig struct {
	User         string
	Password     string
	TerminalID   string
	BaseURL      string
	SignatureKey string
}

func Load() (*Config, error) {
	// Загружаем переменные окружения из .env, если есть
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ .env not found, reading from environment variables")
	}

	cfg := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Airba: AirbaConfig{
			User:         os.Getenv("AIRBA_USER"),
			Password:     os.Getenv("AIRBA_PASSWORD"),
			TerminalID:   os.Getenv("AIRBA_TERMINAL_ID"),
			BaseURL:      os.Getenv("AIRBA_BASE_URL"),
			SignatureKey: os.Getenv("AIRBA_SIGNATURE_KEY"),
		},
	}

	return cfg, nil
}
