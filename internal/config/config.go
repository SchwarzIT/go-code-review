package config

import (
	"coupon_service/internal/api"
	"fmt"
	"log"

	"github.com/brumhard/alligotor"
	"github.com/joho/godotenv"
)

type Config struct {
	API api.Config
}

func New() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: no .env file found (skipping)")
	}

	cfg := Config{}
	if err := alligotor.Get(&cfg); err != nil {
		return cfg, fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.API.Port == "" {
		return cfg, fmt.Errorf("critical environment variable API_PORT is missing")
	}

	return cfg, nil
}
