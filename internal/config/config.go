package config

import (
	"coupon_service/internal/api"
	"fmt"
	"log"

	"github.com/brumhard/alligotor"
	"github.com/joho/godotenv"
)

// Config storage application configuration
type Config struct {
	API api.Config
}

type OptionsConfigFunc func(*Config) error

// New Create new config checking first the env file, opts and the environment variable
func New(envFilePath string, opts ...OptionsConfigFunc) (Config, error) {
	if envFilePath != "" {
		err := godotenv.Load(envFilePath)
		if err != nil {
			log.Printf("Warning: failed to load .env file at %s (skipping)", envFilePath)
		}
	} else {
		err := godotenv.Load()
		if err != nil {
			log.Println("Warning: no .env file found (skipping)")
		}
	}

	cfg := Config{}
	for _, opt := range opts {
		if err := opt(&cfg); err != nil {
			return cfg, err
		}
	}

	if err := alligotor.Get(&cfg); err != nil {
		return cfg, fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.API.Port == "" {
		return cfg, fmt.Errorf("critical environment variable API_PORT is missing")
	}

	return cfg, nil
}
