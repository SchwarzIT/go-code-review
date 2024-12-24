package config

import (
	"coupon_service/internal/api"
	"fmt"
	"log"
	"os"

	"github.com/brumhard/alligotor"
	"github.com/joho/godotenv"
)

type Config struct {
	API api.Config
}

func New(envPath ...string) (Config, error) {
	envFilePath := ""
	if len(envPath) > 0 {
		envFilePath = envPath[0]
	}

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
	cfg.API.Env = getEnv("API_ENV", "dev")
	if err := alligotor.Get(&cfg); err != nil {
		return cfg, fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.API.Port == "" {
		return cfg, fmt.Errorf("critical environment variable API_PORT is missing")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return defaultValue
}
