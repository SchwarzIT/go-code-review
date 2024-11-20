package config

import (
	"log"

	"coupon_service/internal/api"

	"github.com/brumhard/alligotor"
)

type Config struct {
	API api.Config
}

func New() Config {
	cfg := Config{}
	if err := alligotor.Get(&cfg); err != nil {
		log.Fatal(err)
	}
	return cfg
}
