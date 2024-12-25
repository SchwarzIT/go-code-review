package config

import (
	"coupon_service/internal/myduration"
	"time"
)

const (
	defaultEnvironment = "dev"
	defaultTimeAlive   = time.Duration(1) * myduration.HoursInDay * myduration.DaysInYear
	shutdownTimeout    = time.Duration(30) * time.Second
)

// WithDefaultEnv Option func to setup the default environment
func WithDefaultEnv(cfg *Config) error {
	cfg.API.Env = defaultEnvironment
	return nil
}

// WithDefaultShutdownTimeout Option func to setup the default shutdown timeout
func WithDefaultShutdownTimeout(cfg *Config) error {
	cfg.API.ShutdownTimeout = myduration.MyDuration(shutdownTimeout)
	return nil
}

// WithDefaultEnv Option func to setup the default time alive
func WithDefaultTimeAlive(cfg *Config) error {
	cfg.API.TimeAlive = myduration.MyDuration(defaultTimeAlive)
	return nil
}

// NewDefault setup all default option func to create a new Config
func NewDefault(envFilePath string) (Config, error) {
	return New(envFilePath, WithDefaultEnv, WithDefaultTimeAlive, WithDefaultShutdownTimeout)
}
