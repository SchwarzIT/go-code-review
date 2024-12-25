package config

import (
	"coupon_service/internal/mytypes"
	"time"
)

type environment string

const (
	shutdownTimeout = time.Duration(30) * time.Second
)

// WithDefaultEnv Option func to setup the default environment
func WithDefaultEnv(cfg *Config) error {
	cfg.API.Env = mytypes.DefaultEnvironment
	return nil
}

// WithDefaultShutdownTimeout Option func to setup the default shutdown timeout
func WithDefaultShutdownTimeout(cfg *Config) error {
	cfg.API.ShutdownTimeout = mytypes.MyDuration(shutdownTimeout)
	return nil
}

// NewDefault setup all default option func to create a new Config
func NewDefault(envFilePath string) (Config, error) {
	return New(envFilePath, WithDefaultEnv, WithDefaultShutdownTimeout)
}
