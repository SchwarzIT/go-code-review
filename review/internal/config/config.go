package config

import (
	"fmt"

	"github.com/brumhard/alligotor"
)

type Config struct {
	Host string
	Port int
}

func New() (Config, error) {
	cfg := Config{
		// set default values
		Host: "0.0.0.0",
		Port: 8080,
	}
	if err := alligotor.Get(&cfg); err != nil {
		return Config{}, fmt.Errorf("failed to load configurations: %w", err)
	}
	return cfg, nil
}
